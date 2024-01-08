package config

import "os"

type Config struct {
	DBConn         string
	DiscordWebhook string
}

func NewConfig() Config {
	conn := os.Getenv("DB_CONN")
	discord := os.Getenv("DISCORD_WEBHOOK")

	return Config{
		DBConn:         conn,
		DiscordWebhook: discord,
	}
}
