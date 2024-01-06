package storage

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github/elumbantoruan/feed/pkg/feed"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

type Storage interface {
	AddSiteFeed(ctx context.Context, feed feed.Feed) (int64, error)
	GetSiteFeeds(ctx context.Context) ([]feed.Feed, error)
	GetSiteFeed(ctx context.Context, id int) (*feed.Feed, error)
	UpdateSiteFeed(ctx context.Context, feed feed.Feed) error
	AddArticle(ctx context.Context, article feed.Article, siteID int) error
	AddArticles(ctx context.Context, articles []feed.Article) error
	GetArticle(ctx context.Context, id int) (*feed.Article, error)
	GetArticleHash(ctx context.Context, hash string) (*feed.Article, error)
	GetArticles(ctx context.Context) ([]feed.Article, error)
}

type MySQLStorage struct {
	conn string
}

func NewMySQLStorage(conn string) (*MySQLStorage, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()
	return &MySQLStorage{
		conn: conn,
	}, nil
}

func (ms *MySQLStorage) AddSiteFeed(ctx context.Context, feed feed.Feed) (int64, error) {
	query := fmt.Sprintf(`
		INSERT INTO feed.feed_site (
			name, url, type, updated
		) VALUES (
			?, ?, ?, ?
		)
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

	r, err := insert.ExecContext(ctx, feed.Site, feed.RSS, feed.Type, time.Now())
	if err != nil {
		return -1, err
	}

	return r.LastInsertId()
}
func (ms *MySQLStorage) GetSiteFeeds(ctx context.Context) ([]feed.Feed, error) {
	db, err := sql.Open("mysql", ms.conn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, name, url, type, updated FROM feed.feed_site"
	selectQ, err := db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer selectQ.Close()

	rows, err := selectQ.QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	var (
		feedSite  feed.Feed
		feedSites []feed.Feed
	)

	for rows.Next() {
		err := rows.Scan(&feedSite.ID, &feedSite.Site, &feedSite.RSS, &feedSite.Type, &feedSite.Updated)
		if err != nil {
			return nil, err
		}
		feedSites = append(feedSites, feedSite)
	}

	return feedSites, nil
}
func (ms *MySQLStorage) GetSiteFeed(ctx context.Context, id int) (*feed.Feed, error) {
	return nil, nil
}

func (ms *MySQLStorage) UpdateSiteFeed(ctx context.Context, feed feed.Feed) error {
	query := fmt.Sprintf("UPDATE feed.feed_site SET feed.feed_site.updated = ? WHERE feed.feed_site.id = ?")

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

func (ms *MySQLStorage) AddArticle(ctx context.Context, article feed.Article, siteID int64) (int64, error) {

	var authors string
	for i, author := range article.Authors {
		authors += author
		if i < len(article.Authors)-1 {
			authors += " "
		}
	}
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
		return 0, nil
	}

	query := fmt.Sprintf(`
		INSERT INTO feed.feed_content (
			feed_site_id, content_id, title, link, pub_date, description, content, authors, hash
		) VALUES (
			?, ?, ?, ?,	?, ?, ?, ?, ?
		)
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

	r, err := insert.ExecContext(ctx, siteID, article.ID, article.Title, article.Link, article.Published, article.Description, article.Content, authors, hash)
	if err != nil {
		return -1, err
	}

	return r.LastInsertId()
}
func (ms *MySQLStorage) AddArticles(ctx context.Context, articles []feed.Article) error {
	return nil
}
func (ms *MySQLStorage) GetArticle(ctx context.Context, id int) (*feed.Article, error) {
	return nil, nil
}
func (ms *MySQLStorage) GetArticleHash(ctx context.Context, hash string) (*feed.Article, error) {
	query := "SELECT id, title FROM feed.feed_content WHERE hash = ?"
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
func (ms *MySQLStorage) GetArticles(ctx context.Context) ([]feed.Article, error) {
	return nil, nil
}
