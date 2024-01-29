package service

import (
	"context"
	"time"

	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/elumbantoruan/feed/pkg/feedproto"
	pb "github.com/elumbantoruan/feed/pkg/feedproto"
	"github.com/elumbantoruan/feed/pkg/storage"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"log/slog"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type feedServiceServer struct {
	pb.UnimplementedFeedServiceServer
	storage storage.Storage[int64]
	logger  *slog.Logger
	tracer  trace.Tracer
}

func NewFeedServiceServer(st storage.Storage[int64], logger *slog.Logger) feedproto.FeedServiceServer {
	tracer := otel.Tracer("newsfeed-grpc")

	return &feedServiceServer{
		storage: st,
		logger:  logger,
		tracer:  tracer,
	}
}

func (f feedServiceServer) AddSiteFeed(ctx context.Context, pbfeed *pb.Feed) (*empty.Empty, error) {
	ctx, span := f.tracer.Start(ctx, "FeedService.AddSiteFeed")
	defer span.End()

	site := feed.Site[int64]{
		Site: pbfeed.Site,
		RSS:  pbfeed.Rss,
		Type: pbfeed.Type,
	}
	_, err := f.storage.AddSite(ctx, site)
	if err != nil {
		f.logger.Error("AddSiteFeed", slog.Any("error", err))
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (f feedServiceServer) GetSites(ctx context.Context, e *empty.Empty) (*pb.Sites, error) {

	ctx, span := f.tracer.Start(ctx, "FeedService.GetSites")
	defer span.End()

	sites, err := f.storage.GetSites(ctx)
	if err != nil {
		f.logger.Error("GetSites", slog.Any("error", err))
		return nil, err
	}
	var pbSites []*pb.Site
	for _, site := range sites {
		if site.Updated == nil {
			site.Updated = &time.Time{}
		}

		id := &pb.Site_Id{Id: site.ID}
		pbSite := &pb.Site{
			Idtype:       id,
			Site:         site.Site,
			Rss:          site.RSS,
			Type:         site.Type,
			Updated:      timestamppb.New(*site.Updated),
			ArticlesHash: site.ArticlesHash,
		}
		pbSites = append(pbSites, pbSite)
	}
	return &pb.Sites{Sites: pbSites}, nil
}

func (f feedServiceServer) UpdateSite(ctx context.Context, pbsite *pb.Site) (*empty.Empty, error) {
	ctx, span := f.tracer.Start(ctx, "FeedService.UpdateSite")
	defer span.End()

	ts := pbsite.Updated.AsTime()
	site := feed.Site[int64]{
		ID:           pbsite.GetId(),
		Updated:      &ts,
		ArticlesHash: pbsite.ArticlesHash,
	}
	err := f.storage.UpdateSite(ctx, site)
	if err != nil {
		f.logger.Error("UpdateSiteFeed", slog.Any("error", err))
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (f feedServiceServer) UpsertArticle(ctx context.Context, pbarticle *pb.ArticleSite) (*pb.ArticleIdentifier, error) {
	ctx, span := f.tracer.Start(ctx, "FeedService.UpsertArticle")
	defer span.End()

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
	ctx, span := f.tracer.Start(ctx, "FeedService.GetArticles")
	defer span.End()

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
	ctx, span := f.tracer.Start(ctx, "FeedService.GetArticlesWithSite")
	defer span.End()

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
