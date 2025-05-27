package main

import (
	"context"
	"ditto/config"
	"ditto/internal/app"
	"ditto/internal/ditto"
	"ditto/internal/http/handler"
	"ditto/internal/http/router"
	"ditto/internal/repository"
	"ditto/internal/service"
	"ditto/pkg/database"
	"ditto/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		config.Module,
		database.Module,
		app.Module,
		logger.Module,
		fx.Invoke(func(lc fx.Lifecycle, dittoService ditto.Service, dittoClient *ditto.Client, r *gin.Engine) {
			// Start Ditto service
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					return dittoService.Start(ctx)
				},
				OnStop: func(ctx context.Context) error {
					return dittoService.Stop()
				},
			})

			// Initialize repository with Ditto service
			thingRepo := repository.NewThingRepositoryDitto(dittoService, dittoClient)

			// Initialize service
			thingService := service.NewThingService(thingRepo, dittoClient)

			// Initialize handler
			thingHandler := handler.NewThingHandler(thingService)

			// Register routes
			router.RegisterThingRoutes(r.Group("/api/v1"), thingHandler)
		}),
	)

	app.Run()
}
