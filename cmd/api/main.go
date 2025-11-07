package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
	"github.com/yourusername/go-sqlc-starter/internal/api"
	"github.com/yourusername/go-sqlc-starter/internal/config"
	"github.com/yourusername/go-sqlc-starter/internal/db"
)

func main() {
	// Load configuration from environment
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Setup structured logging
	logger := setupLogger(cfg.Env)
	logger.Info().
		Str("env", cfg.Env).
		Str("port", cfg.Port).
		Msg("Starting application")

	// Connect to database
	database, err := db.Connect(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}
	defer database.Close()

	logger.Info().Msg("Database connection established")

	// Run migrations automatically
	if err := db.RunMigrations(cfg.DatabaseURL); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run migrations")
	}

	logger.Info().Msg("Database migrations completed")

	// Initialize router with all dependencies
	router := api.NewRouter(cfg, database, logger)

	// Configure HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Info().
			Str("address", server.Addr).
			Msg("Server starting")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info().Msg("Shutting down server...")

	// Give outstanding requests 10 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server exited gracefully")
}

// setupLogger configures structured logging based on environment
func setupLogger(env string) zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if env == "production" {
		// JSON logging for production
		return zerolog.New(os.Stdout).With().Timestamp().Logger()
	}

	// Pretty console logging for development
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
	}
	return zerolog.New(output).With().Timestamp().Caller().Logger()
}
