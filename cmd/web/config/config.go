package config

import (
	"errors"
	"os"
)

type Config struct {
	GRPCServer      string
	WebPort         string
	HealthCheckPort string
}

func New() (*Config, error) {
	grpcServer := os.Getenv("GRPC_SERVER_ADDRESS")
	if grpcServer == "" {
		return nil, errors.New("empty grpc server address config")
	}
	webPort := os.Getenv("WEB_PORT")
	if webPort == "" {
		return nil, errors.New("empty port config")
	}
	healthCheckPort := os.Getenv("HEALTH_CHECK_PORT")
	if healthCheckPort == "" {
		return nil, errors.New("empty health check port config")
	}
	return &Config{
		GRPCServer:      grpcServer,
		WebPort:         webPort,
		HealthCheckPort: healthCheckPort,
	}, nil
}
