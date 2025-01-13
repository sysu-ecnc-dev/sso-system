package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sysu-ecnc-dev/sso-system/backend/internal/config"
	"github.com/sysu-ecnc-dev/sso-system/backend/internal/handler"
	"github.com/sysu-ecnc-dev/sso-system/backend/internal/repository"
	"github.com/sysu-ecnc-dev/sso-system/backend/internal/utils"
)

func main() {
	// create a slog logger
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// read config
	cfg, err := config.ReadConfig()
	if err != nil {
		slog.Error("Failed to read config.", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// create db connection pool
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.Database.PingTimeout)*time.Second)
	defer cancel()

	dbpool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://postgres:%s@localhost:5432/ecnc_sso_db?sslmode=disable", cfg.Database.Password))
	if err != nil {
		slog.Error("Failed to create database connection pool.", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer dbpool.Close()

	// create sqlc queries instance
	queries := repository.New(dbpool)

	// self check
	if err := utils.EnsureRolesTableInitialized(cfg, queries, dbpool); err != nil {
		slog.Error("Failed to ensure roles table initialized.", slog.String("error", err.Error()))
		os.Exit(1)
	}
	if err := utils.EnsureInitialAdminExists(cfg, queries, dbpool); err != nil {
		slog.Error("Failed to ensure initial admin exists.", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// create the handler and register routes
	handler := handler.New(queries)
	handler.RegisterRoutes()

	// create a http server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      handler.Router,
		IdleTimeout:  time.Duration(cfg.Server.IdleTimeout),
		ReadTimeout:  time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(cfg.Server.WriteTimeout) * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	// start the server
	slog.Info("Starting server...", slog.Int("port", cfg.Server.Port))
	if err := srv.ListenAndServe(); err != nil {
		slog.Error("Failed to start server.", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
