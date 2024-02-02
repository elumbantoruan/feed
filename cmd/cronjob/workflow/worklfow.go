package workflow

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"runtime"

	"github.com/elumbantoruan/feed/pkg/crawler"
	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/elumbantoruan/feed/pkg/grpc/client"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"log/slog"
)

type Workflow struct {
	GRPCClient client.GRPCFeedClient
	Logger     *slog.Logger
	Tracer     trace.Tracer
}

func New(grpcClient client.GRPCFeedClient, logger *slog.Logger) Workflow {
	tracer := otel.Tracer("newsfeed-cronjob")

	workflow := Workflow{
		GRPCClient: grpcClient,
		Logger:     logger,
		Tracer:     tracer,
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

	ctx, rootSpan := w.Tracer.Start(ctx, "Workflow.Run", trace.WithSpanKind(trace.SpanKindClient))
	defer rootSpan.End()

	sites, err := w.GRPCClient.GetSites(ctx)
	if err != nil {
		w.Logger.Error("Run", slog.Any("error", err))
		return nil, err
	}

	workerPools := runtime.NumCPU()
	w.Logger.Info("Run", slog.Int("Number of worker pools", workerPools))

	jobs := len(sites)
	siteC := make(chan feed.Site[int64], jobs)
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

	var anyError bool

	var results Results
	for i := 0; i < jobs; i++ {
		result := <-resultC
		if result.Error != nil {
			w.Logger.Error("Result", slog.Any("error", result.Error))
			anyError = true
		}
		results = append(results, result)
	}

	close(resultC)

	if anyError {
		rootSpan.SetStatus(codes.Error, "Workflow-Run completed with error. Check the log")
	} else {
		rootSpan.SetStatus(codes.Ok, "Workflow-Run completed successfully")
	}

	return results, nil
}

func (w Workflow) worker(ctx context.Context, wID int, siteC <-chan feed.Site[int64], resultC chan<- Result) {
	for site := range siteC {

		cr := crawler.Factory(site)

		w.Logger.Info("Attempt to download", slog.String("url", site.RSS), slog.Int("worker", wID))

		ctx, span := w.Tracer.Start(ctx, "Workflow-worker",
			trace.WithAttributes(attribute.Int("workerID", wID)),
		)
		defer span.End()

		f, err := cr.Download(ctx, site.RSS)
		if err != nil {
			w.Logger.Error("Run - Download", slog.Any("error", err))
			resultC <- Result{WorkerID: wID, Metric: Metric{Site: site.Site}, Error: err}
			return
		}
		f.Site.ID = site.ID

		articlesHash := computeHash(f.Articles)

		if articlesHash == site.ArticlesHash {
			// no update
			resultC <- Result{WorkerID: wID, Metric: Metric{Site: site.Site}, Error: nil}
			return
		} else {
			w.Logger.Info("Update", slog.String("site", site.Site), slog.Time("current ts", *f.Site.Updated), slog.Time("last ts", *site.Updated))

			f.Site.ArticlesHash = articlesHash
			err = w.GRPCClient.UpdateSite(ctx, f.Site)
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

func computeHash(articles []feed.Article) string {
	data, _ := json.Marshal(articles)
	h := sha256.New()
	h.Write(data)
	bs := h.Sum(nil)
	hash := fmt.Sprintf("%x", bs)
	return hash
}
