package service

import (
	"context"
	"time"

	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/elumbantoruan/feed/pkg/feedproto"
	pb "github.com/elumbantoruan/feed/pkg/feedproto"
	"github.com/elumbantoruan/feed/pkg/storage"

	"log/slog"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type feedServiceServer struct {
	pb.UnimplementedFeedServiceServer
	storage storage.Storage[int64]
	logger  *slog.Logger
}

func NewFeedServiceServer(st storage.Storage[int64], logger *slog.Logger) feedproto.FeedServiceServer {
	return &feedServiceServer{
		storage: st,
		logger:  logger,
	}
}

func (f feedServiceServer) AddSiteFeed(ctx context.Context, pbfeed *pb.Feed) (*empty.Empty, error) {

	feed := feed.Feed{
		Site: pbfeed.Site,
		RSS:  pbfeed.Rss,
		Type: pbfeed.Type,
	}
	_, err := f.storage.AddSiteFeed(ctx, feed)
	if err != nil {
		f.logger.Error("AddSiteFeed", slog.Any("error", err))
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (f feedServiceServer) GetSitesFeed(ctx context.Context, e *empty.Empty) (*pb.Feeds, error) {

	sites, err := f.storage.GetSitesFeed(ctx)
	if err != nil {
		f.logger.Error("GetSitesFeed", slog.Any("error", err))
		return nil, err
	}
	var pbFeeds []*pb.Feed
	for _, site := range sites {
		if site.Updated == nil {
			site.Updated = &time.Time{}
		}
		pbFeed := &pb.Feed{
			Id:      site.ID,
			Site:    site.Site,
			Rss:     site.RSS,
			Type:    site.Type,
			Updated: timestamppb.New(*site.Updated),
		}
		pbFeeds = append(pbFeeds, pbFeed)
	}
	return &pb.Feeds{Feeds: pbFeeds}, nil
}

func (f feedServiceServer) UpdateSiteFeed(ctx context.Context, pbfeed *pb.Feed) (*empty.Empty, error) {

	ts := pbfeed.Updated.AsTime()
	feed := feed.Feed{
		ID:      pbfeed.Id,
		Updated: &ts,
	}
	err := f.storage.UpdateSiteFeed(ctx, feed)
	if err != nil {
		f.logger.Error("UpdateSiteFeed", slog.Any("error", err))
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (f feedServiceServer) UpsertArticle(ctx context.Context, pbarticle *pb.ArticleSite) (*pb.ArticleIdentifier, error) {
	var authors []string
	for _, author := range pbarticle.Article.Authors {
		authors = append(authors, author)
	}

	article := feed.Article{
		ID:          pbarticle.Article.Id,
		Title:       pbarticle.Article.Title,
		Link:        pbarticle.Article.Link,
		Published:   pbarticle.Article.Published.AsTime(),
		Description: pbarticle.Article.Description,
		Content:     pbarticle.Article.Content,
		Authors:     authors,
	}
	id, err := f.storage.UpsertArticle(ctx, article, pbarticle.Siteid)
	if err != nil {
		f.logger.Error("UpsertArticle", slog.Any("error", err))
		return nil, err
	}
	return &pb.ArticleIdentifier{
		Identifier: id,
		Article:    pbarticle.Article,
	}, nil
}

func (f feedServiceServer) GetArticles(ctx context.Context, e *empty.Empty) (*pb.ArticlesSite, error) {
	articles, err := f.storage.GetArticles(ctx)
	if err != nil {
		f.logger.Error("GetArticles", slog.Any("error", err))
		return nil, err
	}
	var articlespb []*pb.ArticleSite
	for _, article := range articles {
		var authors []string
		for _, author := range article.Article.Authors {
			authors = append(authors, author)
		}
		articlepb := pb.ArticleSite{
			Siteid: article.SiteID,
			Article: &pb.Article{
				Id:          article.Article.ID,
				Title:       article.Article.Title,
				Link:        article.Article.Link,
				Published:   timestamppb.New(article.Article.Published),
				Description: article.Article.Description,
				Content:     article.Article.Content,
				Authors:     authors,
			},
		}
		articlespb = append(articlespb, &articlepb)
	}
	return &pb.ArticlesSite{
		ArticlesSite: articlespb,
	}, nil
}

func (f *feedServiceServer) GetArticlesWithSite(ctx context.Context, in *pb.SiteId) (*pb.Articles, error) {
	articles, err := f.storage.GetArticlesWithSite(ctx, in.Id, in.LimitRecords)
	if err != nil {
		f.logger.Error("GetArticlesWithSite", slog.Any("error", err))
		return nil, err
	}
	var articlespb []*pb.Article
	for _, article := range articles {
		var authors []string
		for _, author := range article.Authors {
			authors = append(authors, author)
		}
		articlepb := pb.Article{
			Id:          article.ID,
			Title:       article.Title,
			Link:        article.Link,
			Published:   timestamppb.New(article.Published),
			Description: article.Description,
			Content:     article.Content,
			Authors:     authors,
		}
		articlespb = append(articlespb, &articlepb)
	}
	return &pb.Articles{
		Articles: articlespb,
	}, nil
}
