package main

import (
	"github/elumbantoruan/feed/cmd/crownjob/workflow"
	"github/elumbantoruan/feed/pkg/config"
	"log/slog"
	"os"
	"time"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("main", slog.Time("start", time.Now()))

	config := config.NewConfig()

	workflow := workflow.New(config, logger)
	workflow.Run()
}