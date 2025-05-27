package router

import (
	"github.com/gin-gonic/gin"

	"ditto/internal/http/handler"
)

// RegisterThingRoutes registers all thing-related routes
func RegisterThingRoutes(r *gin.RouterGroup, h *handler.ThingHandler) {
	things := r.Group("/things")
	{
		things.POST("", h.Create)
		things.GET("", h.List)
		things.GET("/:id", h.GetByID)
		things.PUT("/:id", h.Update)
		things.DELETE("/:id", h.Delete)
	}
}
