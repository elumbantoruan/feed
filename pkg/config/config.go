package config

import "os"

type Config struct {
	DBConn string
}

func NewConfig() Config {
	conn := os.Getenv("DB_CONN")
	return Config{
		DBConn: conn,
	}
}
