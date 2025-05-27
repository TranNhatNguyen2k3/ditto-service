package router

import (
	"context"
	"ditto/config"
	"ditto/internal/middleware"
	"ditto/pkg/graceful"
	"ditto/pkg/logger"
	"ditto/pkg/swagger"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"go.uber.org/fx"
)

type RouterParams struct {
	fx.In
	HealthRouter HealthRouter
	Logger       logger.Logger
	ErrorHandler *middleware.ErrorHandler
	Config       *config.Config
}

func registerSwaggerHandler(g *gin.Engine) {
	swaggerAPI := g.Group("/swagger")
	swag := swagger.NewSwagger()
	swaggerAPI.Use(swag.SwaggerHandler(false))
	swag.Register(swaggerAPI)
}

func startServer(g *gin.Engine, lifecycle fx.Lifecycle, logger logger.Logger, config *config.Config) {
	gracefulService := graceful.NewService(graceful.WithStopTimeout(time.Second), graceful.WithWaitTime(time.Second))
	gracefulService.Register(g)
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				port := fmt.Sprintf("%d", cast.ToInt(config.Server.Port))
				fmt.Println("run on port:", port)
				go gracefulService.StartServer(g, port)
				return nil
			},
			OnStop: func(context.Context) error {
				gracefulService.Close(logger)
				return nil
			},
		},
	)
}

func NewRouter(params RouterParams) *gin.Engine {
	router := gin.Default()

	router.Use(middleware.CorsMiddleware())
	router.Use(middleware.LoggingMiddleware(params.Logger))
	router.Use(params.ErrorHandler.Handle())

	api := router.Group("/api/v1")
	params.HealthRouter.Register(api)

	// Setup proxy router if target URL is configured
	if params.Config.Proxy.TargetURL != "" {
		SetupProxyRouter(router, params.Config.Proxy.TargetURL)
	}

	return router
}

var Module = fx.Options(
	fx.Provide(middleware.NewErrorHandler),
	fx.Provide(NewHealthRouter),
	fx.Provide(NewRouter),
	fx.Invoke(registerSwaggerHandler, startServer),
)
