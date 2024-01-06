package workflow

import (
	"context"
	"github/elumbantoruan/feed/pkg/config"
	"github/elumbantoruan/feed/pkg/crawler"
	"github/elumbantoruan/feed/pkg/storage"
	"log"
	"log/slog"
)

type Workflow struct {
	Config config.Config
	Logger *slog.Logger
}

func New(config config.Config, logger *slog.Logger) Workflow {
	return Workflow{
		Config: config,
		Logger: logger,
	}
}

func (w Workflow) Run() error {
	st, err := storage.NewMySQLStorage(w.Config.DBConn)
	if err != nil {
		w.Logger.Error("bad connection string", "error", err)
		return err
	}

	sites, err := st.GetSiteFeeds(context.Background())
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

			err = st.UpdateSiteFeed(context.Background(), *f)
			if err != nil {
				w.Logger.Error("UpdateFeed", err)
			}
		}

		for _, article := range f.Articles {
			n, err := st.AddArticle(context.Background(), article, site.ID)
			if err != nil {
				w.Logger.Error("AddArticle", err)
			}
			if n == 0 {
				w.Logger.Info("AddArticle - article not added", slog.String("Existing Article", article.Link))
			} else {
				w.Logger.Info("AddArticle - new article added", slog.Int64("ArticleID", n), slog.String("New Article", article.Link))
			}
		}
	}
	return nil
}
