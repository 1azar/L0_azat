package main

import (
	"L0_azat/internal/config"
	"L0_azat/internal/http-service/handlers/order"
	"L0_azat/internal/lib/logger/sl"
	"L0_azat/internal/service"
	"L0_azat/internal/storage/postgres"
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// Config
	cfg := config.MustLoad()

	// Logger
	logger := spawnLogger(cfg.Env)
	logger.Info("initialising..")

	// Storage
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

	// Nats-streaming
	logger.Debug("starting nats-streaming listener service..")
	natsServ, err := service.New(cfg, storage, logger)
	logger.Debug("nats-streaming listener service started")

	// http-service
	logger.Debug("http service started", slog.String("address", cfg.HttpCfg.Address))

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
	// just to be able testing in local machine fixme: remove this later
	router.Use(cors.Handler(cors.Options{
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	router.Get("/{order_uid}", order.New(logger, storage))

	httpServ := &http.Server{
		Addr:         cfg.HttpCfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpCfg.Timeout,
		WriteTimeout: cfg.HttpCfg.Timeout,
		IdleTimeout:  cfg.HttpCfg.IdleTimeout,
	}

	go func() {
		if err := httpServ.ListenAndServe(); err != nil {
			logger.Error("http serving stopped:", sl.Err(err))
		}
	}()

	logger.Debug(fmt.Sprintf("listening http:%s to provide order info", cfg.HttpCfg.Address))
	logger.Info("service started")

	// client which spams msg.
	//go tests.StartMsgSpam(cfg, 20*time.Second)

	// healthy shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	<-done
	logger.Info("stopping service")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second) //todo: move timeout to cfg
	defer cancel()

	if err := httpServ.Shutdown(ctx); err != nil {
		logger.Error(err.Error())
	}
	storage.Close()
	natsServ.Terminate()

	logger.Info("service gracefully stopped")
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
