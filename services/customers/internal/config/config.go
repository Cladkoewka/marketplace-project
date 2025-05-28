package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB    DBConfig
	HTTP  HTTPConfig
	Log   LogConfig
	Kafka KafkaConfig
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

type KafkaConfig struct {
	Broker string
	Topic  string
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
		Kafka: KafkaConfig{
			Broker: getEnv("KAFKA_BROKER", ""),
			Topic:  getEnv("KAFKA_TOPIC", ""),
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
