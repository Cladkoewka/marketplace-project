package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/Cladkoewka/marketplace-project/services/customers/internal/config"
	"github.com/Cladkoewka/marketplace-project/services/customers/internal/handler"
	"github.com/Cladkoewka/marketplace-project/services/customers/internal/kafka"
	"github.com/Cladkoewka/marketplace-project/services/customers/internal/repository"
	"github.com/Cladkoewka/marketplace-project/services/customers/internal/service"
)

func main() {
	cfg := config.Load()
	logger := initLogger(cfg.Log.Level)
	logger.Info("starting customer service")

	ctx := context.Background()
	db, err := repository.NewPostgresDB(ctx, cfg.DB)
	if err != nil {
		logger.Error("failed to connect to DB", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	repo := repository.NewCustomerRepository(db)
	svc := service.NewCustomerService(repo)

	// Kafka consumer
	consumer := kafka.NewConsumer(cfg.Kafka.Broker, cfg.Kafka.Topic, "customers-group")
	defer consumer.Close()

	go func() {
		time.Sleep(20*time.Second) // Костыль, нужно наверное через healthcheck сделать
		if err := consumer.Consume(ctx, svc.HandleOrderPlaced); err != nil {
			logger.Error("kafka consumer failed", "error", err)
		}
	}()


	h := handler.NewCustomerHandler(svc)


	router := handler.SetupRouter(h)
	logger.Info("HTTP server listening", "port", cfg.HTTP.Port)
	if err := http.ListenAndServe(":"+cfg.HTTP.Port, router); err != nil {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}

func initLogger(level string) *slog.Logger {
	var slogLevel slog.Level
	switch level {
	case "debug":
		slogLevel = slog.LevelDebug
	case "warn":
		slogLevel = slog.LevelWarn
	case "error":
		slogLevel = slog.LevelError
	default:
		slogLevel = slog.LevelInfo
	}

	handler := slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slogLevel})
	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
