package main

import (
	"github/elumbantoruan/feed/cmd/cronjob/config"
	"github/elumbantoruan/feed/cmd/cronjob/workflow"
	"github/elumbantoruan/feed/pkg/storage"
	"log"
	"log/slog"
	"os"
	"time"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("main", slog.Time("start", time.Now()))

	config, err := config.NewConfig()
	if err != nil {
		logger.Error("main", slog.Any("error", err))
	}
	st, err := storage.NewMySQLStorage(config.DBConn)
	if err != nil {
		log.Fatal(err)
	}

	workflow := workflow.New(st, config, logger)
	err = workflow.Run()
	if err != nil {
		log.Fatal(err)
	}
}
