package config

import (
	"errors"
	"os"
)

type Config struct {
	DBConn           string
	GRPCPort         string
	OtelGRPCEndpoint string
}

func NewConfig() (*Config, error) {
	conn := os.Getenv("DB_CONN")
	if conn == "" {
		return nil, errors.New("config: empty connection string")
	}
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		return nil, errors.New("config: empty grpc port")
	}
	otelGRPCEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_GRPC_ENDPOINT")
	if otelGRPCEndpoint == "" {
		return nil, errors.New("config: empty otel grpc endpoint")
	}

	return &Config{
		DBConn:           conn,
		GRPCPort:         port,
		OtelGRPCEndpoint: otelGRPCEndpoint,
	}, nil
}
