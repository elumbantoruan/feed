package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/elumbantoruan/feed/cmd/grpc/config"
	pb "github.com/elumbantoruan/feed/pkg/feedproto"
	"github.com/elumbantoruan/feed/pkg/grpc/service"
	"github.com/elumbantoruan/feed/pkg/otelsetup"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"

	"github.com/elumbantoruan/feed/pkg/storage"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"log/slog"

	"github.com/heptiolabs/healthcheck"
	"google.golang.org/grpc"
)

const (
	healthCheckEndpoint = ":8086"
	metricsEndpoint     = ":9001"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.With(slog.String("service-name", "newsfeed-grpc"))
	logger.Info("main", slog.Time("start", time.Now()), slog.Int("cpu count", runtime.NumCPU()))

	config, err := config.NewConfig()
	if err != nil {
		logger.Error("main", slog.Any("error", err))
		os.Exit(1)
	}
	st, err := storage.NewMySQLStorage(config.DBConn)
	if err != nil {
		logger.Error("main", slog.Any("error", err))
		os.Exit(1)
	}

	// http endpoint for health check
	health := healthcheck.NewHandler()
	go http.ListenAndServe(healthCheckEndpoint, health)

	ctx := context.Background()
	tp := otelsetup.NewTraceProviderGrpc(ctx, config.OtelGRPCEndpoint)
	defer func(ctx context.Context) {
		tp.Shutdown(ctx)
	}(ctx)

	address := fmt.Sprintf(":%s", config.GRPCPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("main", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("main", slog.String("start listening at", address))

	// Setup metrics.
	srvMetrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	reg := prometheus.NewRegistry()
	reg.MustRegister(srvMetrics)

	svc := service.NewFeedServiceServer(st, logger)
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(srvMetrics.UnaryServerInterceptor()),
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	)
	pb.RegisterFeedServiceServer(grpcServer, svc)
	srvMetrics.InitializeMetrics(grpcServer)

	// http endpoint for metrics
	http.Handle("/metrics", promhttp.HandlerFor(
		reg,
		promhttp.HandlerOpts{
			EnableOpenMetrics: true,
		},
	))
	go http.ListenAndServe(metricsEndpoint, nil)

	grpcServer.Serve(lis)
}
