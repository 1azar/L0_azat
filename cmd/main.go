package main

import (
	"L0_azat/internal/config"
	"L0_azat/internal/storage/postgres"
	"log/slog"
	"os"
)

func main() {

	cfg := config.MustLoad()

	logger := spawnLogger(cfg.Env)
	logger.Info("service started")

	logger.Debug("connecting to storage..")
	storage, err := postgres.New(cfg)
	if err != nil {
		logger.Error("connection failed: " + err.Error())
		os.Exit(1)
	}
	defer func() {
		storage.Close()
		logger.Debug("database connection closed")
	}()
	logger.Debug("connection succeeded")

}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func spawnLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
