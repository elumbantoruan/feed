package config

import (
	"errors"
	"os"
)

type Config struct {
	DBConn   string
	GRPCPort string
}

func NewConfig() (*Config, error) {
	conn := os.Getenv("DB_CONN")
	if conn == "" {
		return nil, errors.New("empty connection string")
	}
	port := os.Getenv("GRPC_PORT")
	if port == "" {
		return nil, errors.New("empty port")
	}

	return &Config{
		DBConn:   conn,
		GRPCPort: port,
	}, nil
}
