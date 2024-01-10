package main

import (
	"fmt"
	"github/elumbantoruan/feed/cmd/grpc/server/config"
	"github/elumbantoruan/feed/cmd/grpc/server/service"
	pb "github/elumbantoruan/feed/pkg/feedproto"
	"github/elumbantoruan/feed/pkg/storage"
	"net"
	"net/http"
	"os"
	"time"

	"log/slog"

	"github.com/heptiolabs/healthcheck"
	"google.golang.org/grpc"
)

const (
	healthCheckEndpoint = ":8086"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("main", slog.Time("start", time.Now()))

	config, _ := config.NewConfig()
	st, err := storage.NewMySQLStorage(config.DBConn)
	if err != nil {
		logger.Error("main", slog.Any("error", err))
		os.Exit(1)
	}

	health := healthcheck.NewHandler()

	go http.ListenAndServe(healthCheckEndpoint, health)

	address := fmt.Sprintf(":%s", config.GRPCPort)
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logger.Error("main", slog.Any("error", err))
		os.Exit(1)
	}
	logger.Info("main", slog.String("start listening at", address))

	svc := service.NewFeedServiceServer(st, logger)
	grpcServer := grpc.NewServer()
	pb.RegisterFeedServiceServer(grpcServer, svc)
	grpcServer.Serve(lis)
}
