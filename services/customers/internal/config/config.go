package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB   DBConfig
	HTTP HTTPConfig
	Log  LogConfig
}

type DBConfig struct {
	DSN string // Data source name
}

type HTTPConfig struct {
	Port string
}

type LogConfig struct {
	Level string
}

func Load() Config {
	_ = loadEnv()

	return Config{
		DB: DBConfig{
			DSN: getEnv("DSN", "postgres://user:pass@localhost:5432/customers"),
		},
		HTTP: HTTPConfig{
			Port: getEnv("HTTP_PORT", "8080"),
		},
		Log: LogConfig{
			Level: getEnv("LOG_LEVEL", "info"),
		},
	}
}

func loadEnv() error {
	return godotenv.Load()
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
