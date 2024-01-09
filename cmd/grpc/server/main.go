package main

import (
	"fmt"
	"github/elumbantoruan/feed/cmd/grpc/server/config"
	"github/elumbantoruan/feed/cmd/grpc/server/service"
	pb "github/elumbantoruan/feed/pkg/feedproto"
	"github/elumbantoruan/feed/pkg/storage"
	"log"
	"net"
	"os"
	"time"

	"log/slog"

	"google.golang.org/grpc"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	logger.Info("main", slog.Time("start", time.Now()))

	config, _ := config.NewConfig()
	st, err := storage.NewMySQLStorage(config.DBConn)
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", config.GRPCPort))

	svc := service.NewFeedServiceServer(st, logger)
	grpcServer := grpc.NewServer()
	pb.RegisterFeedServiceServer(grpcServer, svc)
	grpcServer.Serve(lis)
}
