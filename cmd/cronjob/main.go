package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/elumbantoruan/feed/cmd/cronjob/config"
	"github.com/elumbantoruan/feed/cmd/cronjob/workflow"
	"github.com/elumbantoruan/feed/pkg/grpc/client"
	"github.com/elumbantoruan/feed/pkg/lokiwriter"
	"github.com/elumbantoruan/feed/pkg/otelsetup"
)

func main() {

	config, err := config.NewConfig()
	if err != nil {
		log.Fatal(err)
	}

	labels := lokiwriter.Labels{
		"app": "newsfeed-cronjob",
	}

	lkw, shutdown, err := lokiwriter.NewLokiWriter(config.LokiGRPCAddress, labels)
	if err != nil {
		log.Fatal(err)
	}
	defer shutdown()

	// create an instance of slog with loki writer as an exporter
	logger := slog.New(slog.NewJSONHandler(lkw, nil))

	startTime := time.Now()
	logger.Info("main", slog.Time("start", startTime), slog.Int("cpu_count", runtime.NumCPU()))

	svc, err := client.NewGRPCClient(config.GRPCServerAddress)
	if err != nil {
		logger.Error("cronjob - grpc client connection", slog.Any("error", err))
		os.Exit(1)
	}

	ctx := context.Background()
	tp := otelsetup.NewTraceProviderGrpc(ctx, config.OtelGRPCEndpoint)
	defer func(ctx context.Context) {
		tp.Shutdown(ctx)
	}(ctx)

	job := workflow.New(svc, logger)

	res, err := job.Run(ctx)
	if err != nil {
		logger.Error("cronjob - run worklow", slog.Any("error", err))
		os.Exit(1)
	}

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	logger.Info("main", slog.Any("elapsed_time_ms", elapsed.Milliseconds()))

	for _, re := range res {
		if re.Error != nil {
			logger.Error("run workflow result", slog.String("site", re.Metric.Site), slog.Any("error", err))
		} else {
			logger.Info(
				"run workflow result",
				slog.String("site", re.Metric.Site),
				slog.Int("new_articles", re.Metric.NewArticles),
				slog.Int("updated_articles", re.Metric.UpdatedArticles),
			)
		}
	}
}
