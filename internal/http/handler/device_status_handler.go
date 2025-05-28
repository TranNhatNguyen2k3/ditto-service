package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type DeviceStatus struct {
	Status    string    `json:"status"`
	LastSeen  time.Time `json:"lastSeen"`
	Battery   int       `json:"battery,omitempty"`
	Connected bool      `json:"connected"`
	Error     string    `json:"error,omitempty"`
}

type DeviceStatusHandler struct {
	proxy *ProxyHandler
}

func NewDeviceStatusHandler(proxy *ProxyHandler) *DeviceStatusHandler {
	return &DeviceStatusHandler{
		proxy: proxy,
	}
}

// GetDeviceStatus handles GET /api/devices/:deviceId/status
func (h *DeviceStatusHandler) GetDeviceStatus(c *gin.Context) {
	// Forward to Ditto to get device status
	h.proxy.ProxyRequest(c)
}

// UpdateDeviceStatus handles PUT /api/devices/:deviceId/status
func (h *DeviceStatusHandler) UpdateDeviceStatus(c *gin.Context) {
	var status DeviceStatus
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Update last seen time
	status.LastSeen = time.Now()

	// Forward to Ditto
	h.proxy.ProxyRequest(c)
}

// GetDeviceConnectionStatus handles GET /api/devices/:deviceId/connection
func (h *DeviceStatusHandler) GetDeviceConnectionStatus(c *gin.Context) {
	// Forward to Ditto to get connection status
	h.proxy.ProxyRequest(c)
}
