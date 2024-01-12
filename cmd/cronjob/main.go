package main

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/elumbantoruan/feed/cmd/cronjob/config"
	"github.com/elumbantoruan/feed/cmd/cronjob/workflow"
	"github.com/elumbantoruan/feed/pkg/grpc/client"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	startTime := time.Now()
	logger.Info("main", slog.Time("start", startTime))

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

	err = workflow.Run(ctx)
	if err != nil {
		logger.Error("main - run worklow", slog.Any("error", err))
		os.Exit(1)
	}

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	logger.Info("main", slog.Any("elapsed time ms", elapsed.Milliseconds()))

}
