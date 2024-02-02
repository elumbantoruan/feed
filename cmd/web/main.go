package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/elumbantoruan/feed/cmd/web/config"
	"github.com/elumbantoruan/feed/pkg/grpc/client"
	"github.com/elumbantoruan/feed/pkg/otelsetup"
	"github.com/elumbantoruan/feed/pkg/web"
	"github.com/elumbantoruan/feed/pkg/web/storage"
	"github.com/heptiolabs/healthcheck"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

func main() {

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

	ctx := context.Background()
	tp := otelsetup.NewTraceProviderGrpc(ctx, cfg.OtelGRPCEndpoint)
	defer func(ctx context.Context) {
		tp.Shutdown(ctx)
	}(ctx)

	// Web handler
	var handler = web.NewContent(webstorage, logger)

	otelHandler := otelhttp.NewHandler(http.HandlerFunc(handler.RenderContent), "GET Content")
	http.Handle("/", otelHandler)

	health := healthcheck.NewHandler()

	healthCheckEndpoint := fmt.Sprintf(":%s", cfg.HealthCheckPort)

	go http.ListenAndServe(healthCheckEndpoint, health)

	logger.Info("Web UI start serving the service", slog.Time("Time", time.Now()))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", cfg.WebPort), nil); err != nil {
		logger.Error("Problem encountered serving web server", slog.Any("error", err))
	}
}
