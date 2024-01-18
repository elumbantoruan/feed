package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/elumbantoruan/feed/cmd/web/config"
	"github.com/elumbantoruan/feed/pkg/grpc/client"
	"github.com/elumbantoruan/feed/pkg/web"
	"github.com/elumbantoruan/feed/pkg/web/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/heptiolabs/healthcheck"
)

func main() {
	app := fiber.New()

	// components dependency
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}
	grpcClient, err := client.NewGRPCClient(cfg.GRPCServer)
	if err != nil {
		log.Fatal(err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	webstorage := storage.NewWebStorage(grpcClient)

	// Web handler
	var handler = web.NewContent(webstorage, logger)
	app.Get("/", handler.RenderContentRoute)

	health := healthcheck.NewHandler()

	healthCheckEndpoint := fmt.Sprintf(":%s", cfg.HealthCheckPort)

	go http.ListenAndServe(healthCheckEndpoint, health)

	logger.Info("Web UI start serving the service", slog.Time("Time", time.Now()))

	if err := app.Listen(fmt.Sprintf(":%s", cfg.WebPort)); err != nil {
		logger.Error("Problem encountered serving web server", slog.Any("error", err))
	}
}
