package main

import (
	"context"
	"log/slog"
	"os"
	"runtime"
	"time"

	"github.com/elumbantoruan/feed/cmd/cronjob/config"
	"github.com/elumbantoruan/feed/cmd/cronjob/workflow"
	"github.com/elumbantoruan/feed/pkg/grpc/client"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	startTime := time.Now()
	logger.Info("main", slog.Time("start", startTime), slog.Int("cpu count", runtime.NumCPU()))

	config, err := config.NewConfig()
	if err != nil {
		logger.Error("main - config", slog.Any("error", err))
		os.Exit(1)
	}

	svc, err := client.NewGRPCClient(config.GRPCServerAddress)
	if err != nil {
		logger.Error("main - grpc client connection", slog.Any("error", err))
		os.Exit(1)
	}

	workflow := workflow.New(svc, logger)
	ctx := context.Background()

	res, err := workflow.Run(ctx)
	if err != nil {
		logger.Error("main - run worklow", slog.Any("error", err))
		os.Exit(1)
	}

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	logger.Info("main", slog.Any("elapsed time ms", elapsed.Milliseconds()))

	for _, re := range res {
		if re.Error != nil {
			logger.Error("run workflow result", slog.String("site", re.Metric.Site), slog.Any("error", err))
		} else {
			logger.Info(
				"run workflow result",
				slog.String("site", re.Metric.Site),
				slog.Int("new articles", re.Metric.NewArticles),
				slog.Int("updated articles", re.Metric.UpdatedArticles),
			)
		}
	}
}
