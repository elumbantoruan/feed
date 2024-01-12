package workflow

import (
	"context"

	"github.com/elumbantoruan/feed/pkg/crawler"
	"github.com/elumbantoruan/feed/pkg/grpc/client"

	"log"
	"log/slog"
)

type Workflow struct {
	Client *client.GrpcClient
	Logger *slog.Logger
}

func New(client *client.GrpcClient, logger *slog.Logger) Workflow {
	return Workflow{
		Client: client,
		Logger: logger,
	}
}

func (w Workflow) Run(ctx context.Context) error {

	sites, err := w.Client.GetSitesFeed(ctx)
	if err != nil {
		log.Fatal(err)
	}

	for _, site := range sites {
		cr := crawler.CrawlerFactory(site)

		w.Logger.Info("Attempt to download", slog.String("url", site.RSS))

		f, err := cr.Download(site.RSS)
		if err != nil {
			w.Logger.Error("crawler download: ", err)
		}
		f.ID = site.ID

		if f.Updated.Equal(*site.Updated) {
			w.Logger.Info("No update", slog.String("Site", site.Site))
			continue
		} else {
			w.Logger.Info("Update", slog.String("site", site.Site), slog.Time("current ts", *f.Updated), slog.Time("last ts", *site.Updated))

			err = w.Client.UpdateSiteFeed(ctx, *f)
			if err != nil {
				w.Logger.Error("UpdateFeed", err)
			}
		}

		for _, article := range f.Articles {
			id, err := w.Client.AddArticle(ctx, article, site.ID)
			if err != nil {
				w.Logger.Error("AddArticle", err)
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
