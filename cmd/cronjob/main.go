package main

import (
	"github/elumbantoruan/feed/cmd/cronjob/config"
	"github/elumbantoruan/feed/cmd/cronjob/workflow"
	"github/elumbantoruan/feed/pkg/grpc/client"
	"log/slog"
	"os"
	"time"
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
	err = workflow.Run()
	if err != nil {
		logger.Error("main - run worklow", slog.Any("error", err))
		os.Exit(1)
	}

	endTime := time.Now()
	elapsed := endTime.Sub(startTime)
	logger.Info("main", slog.Any("elapsed time ms", elapsed.Milliseconds()))

}
