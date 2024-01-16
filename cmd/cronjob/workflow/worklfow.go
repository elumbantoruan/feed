package workflow

import (
	"context"
	"runtime"

	"github.com/elumbantoruan/feed/pkg/crawler"
	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/elumbantoruan/feed/pkg/grpc/client"

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

type Result struct {
	WorkerID int
	Metric   Metric
	Error    error
}

type Results []Result

type Metric struct {
	Site            string
	Updated         bool
	NewArticles     int
	UpdatedArticles int
}

func (w Workflow) Run(ctx context.Context) (Results, error) {
	sites, err := w.GRPCClient.GetSitesFeed(ctx)
	if err != nil {
		w.Logger.Error("Run", slog.Any("error", err))
		return nil, err
	}

	workerPools := runtime.NumCPU()
	w.Logger.Info("Run", slog.Int("Number of worker pools", workerPools))

	jobs := len(sites)
	siteC := make(chan feed.Feed, jobs)
	resultC := make(chan Result, jobs)

	for i := 1; i <= workerPools; i++ {
		go func(i int) {
			w.worker(ctx, i, siteC, resultC)
		}(i)
	}

	for i := 0; i < jobs; i++ {
		siteC <- sites[i]
	}

	close(siteC)

	var results Results
	for i := 0; i < jobs; i++ {
		result := <-resultC
		if result.Error != nil {
			w.Logger.Error("Result", slog.Any("error", result.Error))
		}
		results = append(results, result)
	}

	close(resultC)

	return results, nil
}

func (w Workflow) worker(ctx context.Context, wID int, siteC <-chan feed.Feed, resultC chan<- Result) {
	for site := range siteC {
		var cr crawler.Crawler

		if w.DefaultCrawler != nil {
			cr = w.DefaultCrawler
		} else {
			cr = crawler.CrawlerFactory(site)
		}

		w.Logger.Info("Attempt to download", slog.String("url", site.RSS), slog.Int("worker", wID))

		f, err := cr.Download(site.RSS)
		if err != nil {
			w.Logger.Error("Run - Download", slog.Any("error", err))
			resultC <- Result{WorkerID: wID, Metric: Metric{Site: site.Site}, Error: err}
			return
		}
		f.ID = site.ID

		if f.Updated.Equal(*site.Updated) {
			// no update
			resultC <- Result{WorkerID: wID, Metric: Metric{Site: site.Site}, Error: nil}
			return
		} else {
			w.Logger.Info("Update", slog.String("site", site.Site), slog.Time("current ts", *f.Updated), slog.Time("last ts", *site.Updated))

			err = w.GRPCClient.UpdateSiteFeed(ctx, *f)
			if err != nil {
				w.Logger.Error("Run - UpdateFeed", slog.Any("error", err))
				resultC <- Result{WorkerID: wID, Metric: Metric{Site: site.Site}, Error: err}
				return
			}
		}

		var (
			newArticle     = 0
			updatedArticle = 0
			updated        bool
		)

		for _, article := range f.Articles {
			newID, err := w.GRPCClient.UpsertArticle(ctx, article, site.ID)
			if err != nil {
				w.Logger.Error("Run - UpsertArticle", slog.Any("error", err))
				resultC <- Result{WorkerID: wID, Metric: Metric{Site: site.Site}, Error: err}
				return
			}
			if newID == 0 {
				updatedArticle++
			} else if newID > 0 {
				newArticle++
			}
		}
		if newArticle != 0 || updatedArticle != 0 {
			updated = true
		}
		resultC <- Result{WorkerID: wID, Metric: Metric{Site: site.Site, Updated: updated, NewArticles: newArticle, UpdatedArticles: updatedArticle}, Error: err}
	}
}
