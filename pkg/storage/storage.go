package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	_ "github.com/go-sql-driver/mysql"
)

type Storage[T any] interface {
	AddSite(ctx context.Context, feed feed.Site[T]) (T, error)
	GetSitesFeed(ctx context.Context) ([]feed.Feed, error)
	UpdateSiteFeed(ctx context.Context, feed feed.Feed) error
	GetSites(ctx context.Context) ([]feed.Site[T], error)
	UpdateSite(ctx context.Context, site feed.Site[T]) error
	UpsertArticle(ctx context.Context, article feed.Article, siteID T) (T, error)
	AddArticles(ctx context.Context, articles []feed.Article) error
	GetArticle(ctx context.Context, id T) (*feed.Article, error)
	GetArticleHash(ctx context.Context, hash string) (*feed.Article, error)
	GetArticles(ctx context.Context) ([]feed.ArticleSite[T], error)
	GetArticlesWithSite(ctx context.Context, siteID T, limit int32) ([]feed.Article, error)
}

type MySQLStorage struct {
	conn   string
	tracer trace.Tracer
}

func NewMySQLStorage(conn string) (Storage[int64], error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	defer db.Close()
	tracer := otel.Tracer("MySQL")
	return &MySQLStorage{
		conn:   conn,
		tracer: tracer,
	}, nil
}

func (ms *MySQLStorage) AddSite(ctx context.Context, site feed.Site[int64]) (int64, error) {
	query := fmt.Sprintf(`
		INSERT INTO feed_site (
			name, url, type, updated
		) VALUES (
			?, ?, ?, ?
		)
	`)

	db, err := otelsql.Open("mysql", ms.conn)
	if err != nil {
		return -1, err
	}
	defer db.Close()

	insert, err := db.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer insert.Close()

	r, err := insert.ExecContext(ctx, site.Site, site.RSS, site.Type, time.Now())
	if err != nil {
		return -1, err
	}

	return r.LastInsertId()
}

func (ms *MySQLStorage) GetSites(ctx context.Context) ([]feed.Site[int64], error) {
	db, err := otelsql.Open("mysql", ms.conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	dbAttr := attribute.KeyValue{Key: attribute.Key("db.name"), Value: attribute.StringValue("mysql")}
	ctx, span := ms.tracer.Start(ctx, "MySQL.GetSites", trace.WithSpanKind(trace.SpanKindClient), trace.WithAttributes(dbAttr))
	defer span.End()

	query := "SELECT id, name, url, type, updated, articles_hash FROM feed_site"
	selectQ, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer selectQ.Close()

	rows, err := selectQ.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	var sites []feed.Site[int64]

	for rows.Next() {
		var site feed.Site[int64]

		err := rows.Scan(&site.ID, &site.Site, &site.RSS, &site.Type, &site.Updated, &site.ArticlesHash)
		if err != nil {
			return nil, err
		}
		sites = append(sites, site)
	}

	return sites, nil
}

func (ms *MySQLStorage) GetSitesFeed(ctx context.Context) ([]feed.Feed, error) {
	db, err := sql.Open("mysql", ms.conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, name, url, type, updated FROM feed_site"
	selectQ, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer selectQ.Close()

	rows, err := selectQ.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	var feedSites []feed.Feed

	for rows.Next() {
		var feedSite feed.Feed

		err := rows.Scan(&feedSite.ID, &feedSite.Site, &feedSite.RSS, &feedSite.Type, &feedSite.Updated)
		if err != nil {
			return nil, err
		}
		feedSites = append(feedSites, feedSite)
	}

	return feedSites, nil
}

func (ms *MySQLStorage) UpdateSite(ctx context.Context, feed feed.Site[int64]) error {
	query := fmt.Sprintf("UPDATE feed_site SET feed_site.updated = ?, feed_site.articles_hash = ? WHERE feed_site.id = ?")

	db, err := sql.Open("mysql", ms.conn)
	if err != nil {
		return err
	}
	defer db.Close()

	update, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer update.Close()

	r, err := update.ExecContext(ctx, feed.Updated, feed.ArticlesHash, feed.ID)
	if err != nil {
		return err
	}

	if v, err := r.RowsAffected(); err != nil || v == 0 {
		if v == 0 {
			return fmt.Errorf("MySQLStorage.UpdateSite. No record being updated for siteID: %d", feed.ID)
		} else if err != nil {
			return fmt.Errorf("MySQLStorage.UpdateSite: %w", err)
		}
	}
	return nil
}

func (ms *MySQLStorage) UpdateSiteFeed(ctx context.Context, feed feed.Feed) error {
	query := fmt.Sprintf("UPDATE feed_site SET feed_site.updated = ? WHERE feed_site.id = ?")

	db, err := sql.Open("mysql", ms.conn)
	if err != nil {
		return err
	}
	defer db.Close()

	update, err := db.Prepare(query)
	if err != nil {
		return err
	}
	defer update.Close()

	r, err := update.ExecContext(ctx, feed.Updated, feed.ID)
	if err != nil {
		return err
	}

	if v, err := r.RowsAffected(); err != nil || v == 0 {
		if v == 0 {
			return fmt.Errorf("MySQLStorage.UpdateFeed. No record being updated for siteID: %d", feed.ID)
		} else if err != nil {
			return fmt.Errorf("MySQLStorage.UpdateFeed: %w", err)
		}
	}
	return nil
}

func (ms *MySQLStorage) UpsertArticle(ctx context.Context, article feed.Article, siteID int64) (int64, error) {
	var authors string
	for i, author := range article.Authors {
		authors += author
		if i < len(article.Authors)-1 {
			authors += ", "
		}
	}
	// some site, the updated value changed but other fields remain the same
	// let set it the same with published, and hashing will compute other fields
	article.Updated = article.Published
	data, _ := json.Marshal(article)
	h := sha256.New()
	h.Write(data)
	bs := h.Sum(nil)
	hash := fmt.Sprintf("%x", bs)

	content, err := ms.GetArticleHash(ctx, hash)
	if err != nil {
		return -1, fmt.Errorf("GetArticleHash %w", err)
	}
	if content != nil {
		// existing hash found
		// -1 indicated no update
		return -1, nil
	}

	// ON DUPLICATE KEY UPDATE means there are duplicate in unique index, thus UPDATE is performed
	query := fmt.Sprintf(`
		INSERT INTO feed_content 
			(feed_site_id, content_id, title, link, pub_date, description, content, authors, hash)
		VALUES 
			(?, ?, ?, ?, ?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE
			title = ?, 
			link = ?, 
			pub_date = ?, 
			description = ?, 
			content = ?, 
			authors = ?, 
			hash = ?;
	`)

	db, err := sql.Open("mysql", ms.conn)
	if err != nil {
		return -1, err
	}
	defer db.Close()

	insert, err := db.Prepare(query)
	if err != nil {
		return -1, err
	}
	defer insert.Close()

	r, err := insert.ExecContext(ctx, siteID, article.ID, article.Title, article.Link, article.Published, article.Description, article.Content, authors, hash,
		article.Title, article.Link, article.Published, article.Description, article.Content, authors, hash)
	if err != nil {
		return -1, err
	}

	rowsUpdated, err := r.RowsAffected()
	// rowsUpdated will be two if there are existing record and get updated.
	// upsert inserts 1 (new record) delete 1 existing record
	if rowsUpdated == 2 {
		return 0, nil
	}

	return r.LastInsertId()
}

func (ms *MySQLStorage) AddArticles(ctx context.Context, articles []feed.Article) error {
	return nil
}

func (ms *MySQLStorage) GetArticle(ctx context.Context, id int64) (*feed.Article, error) {
	return nil, nil
}

func (ms *MySQLStorage) GetArticleHash(ctx context.Context, hash string) (*feed.Article, error) {
	query := "SELECT id, title FROM feed_content WHERE hash = ?"
	db, err := sql.Open("mysql", ms.conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	selectQ, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer selectQ.Close()

	row := selectQ.QueryRowContext(ctx, hash)
	var article feed.Article

	err = row.Scan(&article.ID, &article.Title)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		} else {
			return nil, err
		}
	}

	return &article, nil
}

func (ms *MySQLStorage) GetArticles(ctx context.Context) ([]feed.ArticleSite[int64], error) {
	query := "SELECT feed_site_id, content_id, title, link, pub_date, description, content, authors FROM feed_content ORDER BY id desc LIMIT 100"
	db, err := sql.Open("mysql", ms.conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	selectQ, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer selectQ.Close()

	rows, err := selectQ.QueryContext(ctx)
	if err != nil {
		return nil, err
	}

	var articles []feed.ArticleSite[int64]

	for rows.Next() {
		var (
			authors string
			article feed.ArticleSite[int64]
		)

		err := rows.Scan(&article.SiteID, &article.Article.ID, &article.Article.Title, &article.Article.Link, &article.Article.Published, &article.Article.Description, &article.Article.Content, &authors)
		if err != nil {
			return nil, err
		}
		items := strings.Split(authors, ", ")
		for _, item := range items {
			article.Article.Authors = append(article.Article.Authors, item)
		}
		articles = append(articles, article)
	}
	return articles, nil
}

func (ms *MySQLStorage) GetArticlesWithSite(ctx context.Context, siteID int64, limit int32) ([]feed.Article, error) {
	query := "SELECT content_id, title, link, pub_date, description, content, authors FROM feed_content WHERE feed_site_id = ? ORDER BY id desc LIMIT ?"
	db, err := sql.Open("mysql", ms.conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	selectQ, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer selectQ.Close()

	rows, err := selectQ.QueryContext(ctx, siteID, limit)
	if err != nil {
		return nil, err
	}

	var articles []feed.Article

	for rows.Next() {
		var (
			authors string
			article feed.Article
		)

		err := rows.Scan(&article.ID, &article.Title, &article.Link, &article.Published, &article.Description, &article.Content, &authors)
		if err != nil {
			return nil, err
		}
		items := strings.Split(authors, ", ")
		for _, item := range items {
			article.Authors = append(article.Authors, item)
		}
		articles = append(articles, article)
	}
	return articles, nil
}
