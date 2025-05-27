package router

import (
	"ditto/internal/http/handler"

	"github.com/gin-gonic/gin"
)

type healthRouter struct {
	healthHandler *handler.HealthHandler
}

type HealthRouter interface {
	Register(r *gin.RouterGroup)
}

func NewHealthRouter(healthHandler *handler.HealthHandler) HealthRouter {
	return &healthRouter{
		healthHandler: healthHandler,
	}
}

func (h *healthRouter) Register(rg *gin.RouterGroup) {
	rg.GET("/health", h.healthHandler.HealthCheck())
}
