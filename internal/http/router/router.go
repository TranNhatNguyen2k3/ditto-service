package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/fx"

	"ditto/internal/http/handler"
)

type Router struct {
	engine *gin.Engine
	proxy  *handler.ProxyHandler
}

func NewRouter(engine *gin.Engine, proxy *handler.ProxyHandler) *Router {
	return &Router{
		engine: engine,
		proxy:  proxy,
	}
}

func (r *Router) Setup() {
	// Health check endpoint
	r.engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	// Setup device routes
	r.SetupDeviceRoutes()

	// Proxy all other requests to Ditto
	api := r.engine.Group("/api")
	{
		// Proxy /api/things/* to Ditto
		api.Any("/things/*path", r.proxy.ProxyRequest)
	}
}

func (r *Router) SetupDeviceRoutes() {
	// Initialize handlers
	statusHandler := handler.NewDeviceStatusHandler(r.proxy)
	telemetryHandler := handler.NewTelemetryHandler(r.proxy)
	commandHandler := handler.NewDeviceCommandHandler(r.proxy)

	// Device management routes
	devices := r.engine.Group("/api/devices")
	{
		// Basic CRUD operations
		devices.GET("", r.proxy.ProxyRequest)
		devices.GET("/:deviceId", r.proxy.ProxyRequest)
		devices.POST("", r.proxy.ProxyRequest)
		devices.PUT("/:deviceId", r.proxy.ProxyRequest)
		devices.DELETE("/:deviceId", r.proxy.ProxyRequest)

		// Device status routes
		devices.GET("/:deviceId/status", statusHandler.GetDeviceStatus)
		devices.PUT("/:deviceId/status", statusHandler.UpdateDeviceStatus)
		devices.GET("/:deviceId/connection", statusHandler.GetDeviceConnectionStatus)

		// Device telemetry routes
		devices.GET("/:deviceId/telemetry", telemetryHandler.GetDeviceTelemetry)
		devices.GET("/:deviceId/telemetry/history", telemetryHandler.GetDeviceTelemetryHistory)
		devices.GET("/:deviceId/telemetry/stats", telemetryHandler.GetDeviceTelemetryStats)

		// Device command routes
		devices.POST("/:deviceId/commands", commandHandler.SendCommand)
		devices.GET("/:deviceId/commands", commandHandler.ListCommands)
		devices.GET("/:deviceId/commands/:commandId", commandHandler.GetCommandStatus)
	}
}

var Module = fx.Options(
	fx.Provide(
		NewRouter,
	),
)
