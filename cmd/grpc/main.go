package main

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/elumbantoruan/feed/cmd/grpc/config"
	pb "github.com/elumbantoruan/feed/pkg/feedproto"
	"github.com/elumbantoruan/feed/pkg/grpc/service"

	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"

	"github.com/elumbantoruan/feed/pkg/storage"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"

	"log/slog"

	"github.com/heptiolabs/healthcheck"
	"google.golang.org/grpc"
)

const (
	healthCheckEndpoint = ":8086"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("main", slog.Time("start", time.Now()), slog.Int("cpu count", runtime.NumCPU()))

	config, _ := config.NewConfig()
	st, err := storage.NewMySQLStorage(config.DBConn)
	if err != nil {
		logger.Error("main", slog.Any("error", err))
		os.Exit(1)
	}

	health := healthcheck.NewHandler()

	go http.ListenAndServe(healthCheckEndpoint, health)

	dsp := honeycomb.NewDynamicAttributeSpanProcessor(func() []attribute.KeyValue {
		return []attribute.KeyValue{}
	})
	bsp := honeycomb.NewBaggageSpanProcessor()

	shutdown, err := otelconfig.ConfigureOpenTelemetry(
		otelconfig.WithSpanProcessor(dsp, bsp),
	)
	if err != nil {
		logger.Error("main - failed from ConfigurationOpenTelemetry", slog.Any("error", err))
		os.Exit(1)
	}

	defer shutdown()

	address := fmt.Sprintf(":%s", config.GRPCPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("main", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("main", slog.String("start listening at", address))

	svc := service.NewFeedServiceServer(st, logger)
	grpcServer := grpc.NewServer(grpc.StatsHandler(otelgrpc.NewServerHandler()))
	pb.RegisterFeedServiceServer(grpcServer, svc)
	grpcServer.Serve(lis)
}
