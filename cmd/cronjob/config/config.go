package config

import (
	"errors"
	"os"
)

type Config struct {
	DBConn            string
	DiscordWebhook    string
	GRPCServerAddress string
	OtelGRPCEndpoint  string
	LokiGRPCAddress   string
}

func NewConfig() (*Config, error) {
	conn := os.Getenv("DB_CONN")
	discord := os.Getenv("DISCORD_WEBHOOK")
	grpcServerAddress := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcServerAddress == "" {
		return nil, errors.New("empty grpc server address in config")
	}
	otelGRPCEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_GRPC_ENDPOINT")
	if otelGRPCEndpoint == "" {
		return nil, errors.New("empty otel grpc endpoint in config")
	}
	lokiGRPCAddress := os.Getenv("LOKI_GRPC_ADDRESS")
	if lokiGRPCAddress == "" {
		return nil, errors.New("empty loki grpc address in config")
	}

	return &Config{
		DBConn:            conn,
		DiscordWebhook:    discord,
		GRPCServerAddress: grpcServerAddress,
		OtelGRPCEndpoint:  otelGRPCEndpoint,
		LokiGRPCAddress:   lokiGRPCAddress,
	}, nil
}
