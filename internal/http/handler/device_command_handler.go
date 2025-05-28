package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type DeviceCommand struct {
	Command    string                 `json:"command"`
	Parameters map[string]interface{} `json:"parameters,omitempty"`
	Headers    map[string]string      `json:"headers,omitempty"`
}

type DeviceCommandHandler struct {
	proxy *ProxyHandler
}

func NewDeviceCommandHandler(proxy *ProxyHandler) *DeviceCommandHandler {
	return &DeviceCommandHandler{
		proxy: proxy,
	}
}

// SendCommand handles POST /api/devices/:deviceId/commands
func (h *DeviceCommandHandler) SendCommand(c *gin.Context) {
	var cmd DeviceCommand
	if err := c.ShouldBindJSON(&cmd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Forward to Ditto
	h.proxy.ProxyRequest(c)
}

// GetCommandStatus handles GET /api/devices/:deviceId/commands/:commandId
func (h *DeviceCommandHandler) GetCommandStatus(c *gin.Context) {
	// Forward to Ditto
	h.proxy.ProxyRequest(c)
}

// ListCommands handles GET /api/devices/:deviceId/commands
func (h *DeviceCommandHandler) ListCommands(c *gin.Context) {
	// Forward to Ditto
	h.proxy.ProxyRequest(c)
}
