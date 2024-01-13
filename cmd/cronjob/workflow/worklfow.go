package workflow

import (
	"context"

	"github.com/elumbantoruan/feed/pkg/crawler"
	"github.com/elumbantoruan/feed/pkg/grpc/client"

	"log"
	"log/slog"
)

type Workflow struct {
	GRPCClient     client.GRPCFeedClient
	Logger         *slog.Logger
	DefaultCrawler crawler.Crawler
}

func New(grpcClient client.GRPCFeedClient, logger *slog.Logger, defaultCrawler ...crawler.Crawler) Workflow {
	workflow := Workflow{
		GRPCClient: grpcClient,
		Logger:     logger,
	}
	if len(defaultCrawler) != 0 {
		workflow.DefaultCrawler = defaultCrawler[0]
	}
	return workflow
}

func (w Workflow) Run(ctx context.Context) error {

	sites, err := w.GRPCClient.GetSitesFeed(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, site := range sites {

		var cr crawler.Crawler

		if w.DefaultCrawler != nil {
			cr = w.DefaultCrawler
		} else {
			cr = crawler.CrawlerFactory(site)
		}

		w.Logger.Info("Attempt to download", slog.String("url", site.RSS))

		f, err := cr.Download(site.RSS)
		if err != nil {
			w.Logger.Error("Run - Download", slog.Any("error", err))
		}
		f.ID = site.ID

		if f.Updated.Equal(*site.Updated) {
			w.Logger.Info("No update", slog.String("Site", site.Site))
			continue
		} else {
			w.Logger.Info("Update", slog.String("site", site.Site), slog.Time("current ts", *f.Updated), slog.Time("last ts", *site.Updated))

			err = w.GRPCClient.UpdateSiteFeed(ctx, *f)
			if err != nil {
				w.Logger.Error("Run - UpdateFeed", slog.Any("error", err))
			}
		}

		for _, article := range f.Articles {
			id, err := w.GRPCClient.AddArticle(ctx, article, site.ID)
			if err != nil {
				w.Logger.Error("Run - AddArticle", slog.Any("error", err))
			}
			if id == 0 {
				w.Logger.Info("AddArticle - article not added", slog.String("Existing Article", article.Link))
			} else {
				w.Logger.Info("AddArticle - new article added", slog.Int64("ArticleID", id), slog.String("New Article", article.Link))
			}
		}
	}

	return nil
}
