package config

import (
	"errors"
	"os"
)

type Config struct {
	GRPCServer       string
	WebPort          string
	HealthCheckPort  string
	OtelGRPCEndpoint string
}

func New() (*Config, error) {
	grpcServer := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcServer == "" {
		return nil, errors.New("config: empty grpc server address")
	}
	webPort := os.Getenv("WEB_PORT")
	if webPort == "" {
		return nil, errors.New("config: empty web port")
	}
	healthCheckPort := os.Getenv("HEALTH_CHECK_PORT")
	if healthCheckPort == "" {
		return nil, errors.New("config: empty health check port")
	}
	otelGRPCEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_GRPC_ENDPOINT")
	if otelGRPCEndpoint == "" {
		return nil, errors.New("config: empty otel grpc endpoint")
	}
	return &Config{
		GRPCServer:       grpcServer,
		WebPort:          webPort,
		HealthCheckPort:  healthCheckPort,
		OtelGRPCEndpoint: otelGRPCEndpoint,
	}, nil
}
