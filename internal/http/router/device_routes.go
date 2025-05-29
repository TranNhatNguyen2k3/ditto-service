package router

import (
	"log"

	"ditto/config"
	"ditto/internal/http/handler"

	"github.com/gin-gonic/gin"
)

// SetupDeviceRoutes configures all device-related routes
func SetupDeviceRoutes(router *gin.RouterGroup, config *config.Config) {
	// Initialize handler
	deviceHandler := handler.NewDeviceHandler(config.Proxy.TargetURL, config.Proxy.Username, config.Proxy.Password)

	// Device routes group
	deviceGroup := router.Group("/devices")
	{
		// List things with filtering
		deviceGroup.GET("", deviceHandler.ListThings)

		// Create/Update thing
		deviceGroup.PUT("/:thingId", deviceHandler.CreateThing)

		// Get thing state
		deviceGroup.GET("/:thingId/state", deviceHandler.GetThingState)

		// Send commands
		deviceGroup.PUT("/:thingId/features/:feature/command", deviceHandler.SendCommand)
		deviceGroup.POST("/:thingId/features/:feature/command", deviceHandler.SendCommand)
	}

	// Log registered device routes
	log.Printf("Registered device routes:")
	log.Printf("GET /api/devices")
	log.Printf("PUT /api/devices/:thingId")
	log.Printf("GET /api/devices/:thingId/state")
	log.Printf("PUT /api/devices/:thingId/features/:feature/command")
	log.Printf("POST /api/devices/:thingId/features/:feature/command")
}
