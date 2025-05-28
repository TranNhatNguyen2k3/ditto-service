package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type ProxyHandler struct {
	dittoURL string
	username string
	password string
}

func NewProxyHandler(dittoURL, username, password string) *ProxyHandler {
	// Ensure dittoURL ends with /api/2
	if !strings.HasSuffix(dittoURL, "/api/2") {
		dittoURL = strings.TrimRight(dittoURL, "/") + "/api/2"
	}

	return &ProxyHandler{
		dittoURL: dittoURL,
		username: username,
		password: password,
	}
}

func (h *ProxyHandler) ProxyRequest(c *gin.Context) {
	// Get the path from the request
	path := c.Param("path")
	if path == "" {
		path = c.Request.URL.Path
	}

	// Remove /api prefix if present
	path = strings.TrimPrefix(path, "/api")

	// Map /devices to /things for Ditto API
	path = strings.Replace(path, "/devices", "/things", 1)

	// Map /policies to /api/2/policies for Ditto API
	if strings.Contains(path, "/policies") {
		// Extract the policy ID from the path
		parts := strings.Split(path, "/policies/")
		if len(parts) == 2 {
			// Reconstruct the path in Ditto format
			path = "/api/2/policies/" + parts[1]
		}
	}

	// Map /commands to /inbox/messages for Ditto API
	if strings.Contains(path, "/commands") {
		// Extract the command subject from the path
		parts := strings.Split(path, "/commands/")
		if len(parts) == 2 {
			// Reconstruct the path in Ditto format
			path = strings.Replace(parts[0], "/commands", "/inbox/messages", 1) + "/" + parts[1]
		}
	}

	// Construct the target URL
	targetURL := fmt.Sprintf("%s%s", strings.TrimRight(h.dittoURL, "/"), path)
	log.Printf("Proxying request to: %s", targetURL)

	// Create a new request
	req, err := http.NewRequest(c.Request.Method, targetURL, c.Request.Body)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create request: %v", err)})
		return
	}

	// Copy headers from the original request
	for name, values := range c.Request.Header {
		for _, value := range values {
			req.Header.Add(name, value)
		}
	}

	// Add Basic Auth
	req.SetBasicAuth(h.username, h.password)

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to send request: %v", err)})
		return
	}
	defer resp.Body.Close()

	log.Printf("Received response with status: %d", resp.StatusCode)

	// Copy response headers
	for name, values := range resp.Header {
		for _, value := range values {
			c.Header(name, value)
		}
	}

	// Set response status code
	c.Status(resp.StatusCode)

	// Copy response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read response: %v", err)})
		return
	}

	// Try to parse as JSON for pretty printing
	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err == nil {
		c.JSON(resp.StatusCode, jsonData)
	} else {
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}
