package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/BernsteinMondy/subscription-service/internal/controller"
	"github.com/BernsteinMondy/subscription-service/internal/middleware"
	"github.com/BernsteinMondy/subscription-service/internal/repository"
	"github.com/BernsteinMondy/subscription-service/internal/service"
	"github.com/BernsteinMondy/subscription-service/pkg/database"
	"log"
	"log/slog"
	"net/http"
	"os/signal"
	"sync"
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
	slog.Info("Config loaded:", slog.Any("config", cfg))

	slog.Info("Creating new database connection...")
	db, err := newDatabaseConnection(cfg.DB)
	if err != nil {
		return fmt.Errorf("newDatabaseConnection() returned error: %v", err)
	}
	slog.Info("Database connection created")
	defer func() {
		slog.Info("Closing database connection...")
		closeErr := db.Close()
		if closeErr != nil {
			err = errors.Join(err, fmt.Errorf("closing database connection: %w", closeErr))
		} else {
			slog.Info("Database connection closed")
		}
	}()

	// New repository
	repo := repository.New(db)

	// New service
	srvc := service.NewService(repo)

	// New controller
	ctrl := controller.New(srvc)

	// New HTTP multiplexer
	mux := http.NewServeMux()

	// Map handlers
	ctrl.MapHandlers(mux)

	// Middleware
	handlerWithMw := middleware.LoggingMiddleware(mux)

	// New HTTP server
	httpServer := &http.Server{
		Addr:    cfg.HTTPServer.ListenAddr,
		Handler: handlerWithMw,
	}

	wg := &sync.WaitGroup{}

	wg.Add(1)
	go func(ctx context.Context, server *http.Server) {
		defer wg.Done()
		defer slog.Info("HTTP server shutdown")

		slog.Info("Launching HTTP server...")
		launchErr := launchHTTPServer(ctx, httpServer)
		if launchErr != nil {
			err = errors.Join(err, launchErr)
		}
	}(ctx, httpServer)

	wg.Wait()
	<-ctx.Done()
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
	errCh := make(chan error, 1)

	go func() {
		launchErr := httpServer.ListenAndServe()
		if launchErr != nil {
			errCh <- fmt.Errorf("launch http server: %v", launchErr)
		}

		return
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := httpServer.Shutdown(shutdownCtx)
	if err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	return nil
}
