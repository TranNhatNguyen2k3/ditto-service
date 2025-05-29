package main

import (
	"context"
	"ditto/config"
	"ditto/internal/app"
	"ditto/internal/ditto"
	"ditto/internal/http"
	"ditto/internal/influxdb"
	"ditto/pkg/database"
	"ditto/pkg/logger"
	"log"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/fx"
)

func main() {
	application := fx.New(
		config.Module,
		database.Module,
		app.Module,
		http.Module,
		logger.Module,
		fx.Provide(
			// Initialize InfluxDB client
			func(cfg *config.Config) *influxdb.Client {
				return influxdb.NewClient(cfg.InfluxDB.URL, cfg.InfluxDB.Token, cfg.InfluxDB.Org, cfg.InfluxDB.Bucket)
			},
			// Initialize Ditto client
			func(cfg *config.Config) *ditto.Client {
				return ditto.NewClient(cfg.Ditto.WSURL, cfg.Ditto.Username, cfg.Ditto.Password)
			},
			// Initialize Ditto service
			ditto.NewService,
		),
		fx.Invoke(func(lc fx.Lifecycle, app *app.App, dittoService ditto.Service, logger logger.Logger) {
			// Start Ditto service and HTTP server
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Info("Starting Ditto service...")
					if err := dittoService.Start(ctx); err != nil {
						return err
					}
					return app.Start(ctx)
				},
				OnStop: func(ctx context.Context) error {
					logger.Info("Stopping Ditto service...")
					return dittoService.Stop()
				},
			})
		}),
	)

	// Start the application
	startCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		if err := application.Start(startCtx); err != nil {
			log.Printf("Failed to start application: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	// Graceful shutdown
	stopCtx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := application.Stop(stopCtx); err != nil {
		log.Printf("Failed to stop application: %v", err)
		os.Exit(1)
	}
}
