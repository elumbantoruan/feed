package config

import (
	"os"
)

type Config struct {
	DBConn            string
	DiscordWebhook    string
	GRPCServerAddress string
}

func NewConfig() (*Config, error) {
	conn := os.Getenv("DB_CONN")
	discord := os.Getenv("DISCORD_WEBHOOK")
	serverAddress := os.Getenv("GRPC_SERVER_ADDRESS")

	return &Config{
		DBConn:            conn,
		DiscordWebhook:    discord,
		GRPCServerAddress: serverAddress,
	}, nil
}
