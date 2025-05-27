package handler

import (
	"ditto/pkg/logger"
	"ditto/pkg/wrapper"

	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	logger logger.Logger
}

func NewHealthHandler(logger logger.Logger) *HealthHandler {
	return &HealthHandler{logger: logger}
}

func (h *HealthHandler) HealthCheck() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.logger.Info("health check called")
		wrapper.JSONOk(ctx, map[string]interface{}{"status": "ok"})
	}
}
