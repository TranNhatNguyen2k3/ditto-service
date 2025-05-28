package app

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"ditto/config"
	"ditto/pkg/logger"

	"github.com/gin-gonic/gin"
)

type App struct {
	engine *gin.Engine
	config *config.Config
	logger logger.Logger
}

func NewApp(engine *gin.Engine, config *config.Config, logger logger.Logger) *App {
	return &App{
		engine: engine,
		config: config,
		logger: logger,
	}
}

func (a *App) Start(ctx context.Context) error {
	// Create HTTP server
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", a.config.Server.Port),
		Handler: a.engine,
	}

	// Start server in a goroutine
	go func() {
		a.logger.Info(fmt.Sprintf("Starting HTTP server on port %s...", a.config.Server.Port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			a.logger.Error(fmt.Sprintf("Failed to start HTTP server: %v", err))
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	// Shutdown server gracefully
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("failed to shutdown HTTP server: %v", err)
	}

	return nil
}
