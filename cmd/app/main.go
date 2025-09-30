package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/subscription-service/internal/controller"
	"github.com/BernsteinMondy/subscription-service/internal/middleware"
	"github.com/BernsteinMondy/subscription-service/internal/migrations"
	"github.com/BernsteinMondy/subscription-service/internal/repository"
	"github.com/BernsteinMondy/subscription-service/internal/service"
	"github.com/BernsteinMondy/subscription-service/pkg/database"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	err := run()
	if err != nil {
		log.Fatalf("run() returned error: %v", err)
	}
}

func run() (err error) {
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	defer cancel()

	slog.Info("Loading config...")
	cfg, err := loadConfigFromEnv()
	if err != nil {
		return fmt.Errorf("load config: %v", err)
	}
	slog.Info("Config loaded", slog.Any("config", cfg))

	slog.Info("Creating new database connection...")
	db, err := newDatabaseConnection(cfg.DB)
	if err != nil {
		return fmt.Errorf("newDatabaseConnection: %w", err)
	}
	slog.Info("Database connection created")

	defer func() {
		slog.Info("Closing database connection...")
		if closeErr := db.Close(); closeErr != nil {
			slog.Error("Failed to close database", slog.Any("error", closeErr))
		}
		slog.Info("Database connection closed")
	}()

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	// Migrations
	if cfg.Migrations.Enabled {
		slog.Info("Migrations enabled - running migrations")
		err = migrations.Run(ctx, db, cfg.Migrations.Dir)
		if err != nil {
			return fmt.Errorf("run migrations: %w", err)
		}
		slog.Info("Successfully run migrations")
	} else {
		slog.Info("Migrations disabled - skipping migrations")
	}

	select {
	case <-ctx.Done():
		return nil
	default:
	}

	// Repository - Service - Controller
	repo := repository.New(db)
	srvc := service.NewService(repo)
	ctrl := controller.New(srvc)

	// HTTP mux and middleware
	mux := http.NewServeMux()
	ctrl.MapHandlers(mux)
	handlerWithMw := middleware.LoggingMiddleware(mux)

	// HTTP server
	httpServer := &http.Server{
		Addr:    cfg.HTTPServer.ListenAddr,
		Handler: handlerWithMw,
	}

	if err = launchHTTPServer(ctx, httpServer); err != nil {
		slog.Error("HTTP server error", slog.Any("error", err))
		return err
	}

	slog.Info("Server stopped gracefully")
	return nil
}

func newDatabaseConnection(c DB) (*sql.DB, error) {
	dbCfg := &database.Config{
		Host:     c.Host,
		Port:     c.Port,
		User:     c.User,
		Password: c.Password,
		DBName:   c.DatabaseName,
		SSLMode:  c.SSLMode,
	}

	return database.NewConnection(dbCfg)
}

func launchHTTPServer(ctx context.Context, httpServer *http.Server) error {
	serverErr := make(chan error, 1)

	go func() {
		slog.Info("Starting HTTP server", slog.String("address", httpServer.Addr))
		if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- fmt.Errorf("listen on %s: %w", httpServer.Addr, err)
		} else {
			serverErr <- nil
		}
	}()

	<-ctx.Done()
	slog.Info("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("Error during server shutdown", slog.Any("error", err))
		return fmt.Errorf("shutdown http server: %w", err)
	}

	if err := <-serverErr; err != nil {
		return err
	}

	slog.Info("Server shutdown complete")
	return nil
}
