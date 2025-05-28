package handler

import (
	"time"

	"github.com/gin-gonic/gin"
)

type TelemetryData struct {
	DeviceId     string                 `json:"deviceId"`
	Timestamp    time.Time              `json:"timestamp"`
	Measurements map[string]interface{} `json:"measurements"`
}

type TelemetryHandler struct {
	proxy *ProxyHandler
}

func NewTelemetryHandler(proxy *ProxyHandler) *TelemetryHandler {
	return &TelemetryHandler{
		proxy: proxy,
	}
}

// GetDeviceTelemetry handles GET /api/devices/:deviceId/telemetry
func (h *TelemetryHandler) GetDeviceTelemetry(c *gin.Context) {
	// Forward to Ditto with query parameters
	h.proxy.ProxyRequest(c)
}

// GetDeviceTelemetryHistory handles GET /api/devices/:deviceId/telemetry/history
func (h *TelemetryHandler) GetDeviceTelemetryHistory(c *gin.Context) {
	// Forward to Ditto with query parameters
	h.proxy.ProxyRequest(c)
}

// GetDeviceTelemetryStats handles GET /api/devices/:deviceId/telemetry/stats
func (h *TelemetryHandler) GetDeviceTelemetryStats(c *gin.Context) {
	// Forward to Ditto with query parameters
	h.proxy.ProxyRequest(c)
}
