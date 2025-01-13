package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sysu-ecnc-dev/sso-system/backend/internal/config"
	"github.com/sysu-ecnc-dev/sso-system/backend/internal/handler"
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

	// create the handler and register routes
	r := chi.NewRouter()
	handler.RegisterRoutes(r)

	// create a http server
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
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
