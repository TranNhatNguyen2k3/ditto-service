package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"ditto/config"

	"github.com/gin-gonic/gin"
)

type DeviceHandler struct {
	config *config.Config
}

type Thing struct {
	ThingID    string                 `json:"thingId"`
	PolicyID   string                 `json:"policyId"`
	Attributes map[string]interface{} `json:"attributes"`
	Features   map[string]interface{} `json:"features"`
}

type ThingsResponse struct {
	Items []Thing `json:"items"`
	Total int     `json:"total"`
}

type ThingState struct {
	ThingID  string                 `json:"thingId"`
	Features map[string]interface{} `json:"features"`
}

type CommandRequest struct {
	Command string                 `json:"command"`
	Params  map[string]interface{} `json:"params"`
}

// NewDeviceHandler creates a new DeviceHandler
func NewDeviceHandler(config *config.Config) *DeviceHandler {
	return &DeviceHandler{
		config: config,
	}
}

// ListThings handles GET /api/devices
func (h *DeviceHandler) ListThings(c *gin.Context) {
	// Get filter parameters
	namespace := c.Query("namespace")
	company := c.Query("company")
	location := c.Query("location")

	// Build filter query
	filter := fmt.Sprintf("namespace=\"%s\"", namespace)
	if company != "" {
		filter += fmt.Sprintf(" AND attributes/company=\"%s\"", company)
	}
	if location != "" {
		filter += fmt.Sprintf(" AND attributes/location=\"%s\"", location)
	}

	// Build Ditto API URL
	url := fmt.Sprintf("%s/things?filter=%s", h.config.Ditto.URL, url.QueryEscape(filter))

	// Create request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create request"})
		return
	}

	// Add Ditto authentication
	req.SetBasicAuth(h.config.Ditto.Username, h.config.Ditto.Password)

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request to Ditto: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to send request: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read response: %v", err)})
		return
	}

	// Parse response as array of things
	var things []Thing
	if err := json.Unmarshal(body, &things); err != nil {
		log.Printf("Failed to parse things response: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse things response: %v", err)})
		return
	}

	// Add additional filtering
	filteredThings := make([]Thing, 0)
	for _, thing := range things {
		// Check if thing matches all filter criteria
		matches := true

		// Filter by company if specified
		if company != "" {
			if thingCompany, ok := thing.Attributes["company"].(string); !ok || thingCompany != company {
				matches = false
			}
		}

		// Filter by location if specified
		if location != "" {
			if thingLocation, ok := thing.Attributes["location"].(string); !ok || thingLocation != location {
				matches = false
			}
		}

		// Add thing to filtered list if it matches all criteria
		if matches {
			filteredThings = append(filteredThings, thing)
		}
	}

	// Create response with filtered items
	response := ThingsResponse{
		Items: filteredThings,
		Total: len(filteredThings),
	}

	log.Printf("Found %d things matching filters (location: %s, company: %s)",
		response.Total, location, company)

	c.JSON(http.StatusOK, response)
}

// CreateThing handles creating a new thing
func (h *DeviceHandler) CreateThing(c *gin.Context) {
	thingID := c.Param("thingId")
	if thingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameter: thingId"})
		return
	}

	// Parse the thing data from request body
	var thing Thing
	if err := c.ShouldBindJSON(&thing); err != nil {
		log.Printf("Failed to parse thing data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid thing data: %v", err)})
		return
	}

	// Ensure ThingID matches the URL parameter
	thing.ThingID = thingID

	// Convert to JSON
	payload, err := json.Marshal(thing)
	if err != nil {
		log.Printf("Failed to marshal thing data: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to marshal thing data: %v", err)})
		return
	}

	// Construct the target URL for Ditto API
	targetURL := fmt.Sprintf("%s/things/%s",
		strings.TrimRight(h.config.Ditto.URL, "/"),
		thingID)

	log.Printf("Creating thing at: %s", targetURL)
	log.Printf("Thing data: %s", string(payload))

	// Create a new request
	req, err := http.NewRequest(http.MethodPut, targetURL, strings.NewReader(string(payload)))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create request: %v", err)})
		return
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Add Ditto authentication
	req.SetBasicAuth(h.config.Ditto.Username, h.config.Ditto.Password)

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request to Ditto: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to send request: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read response: %v", err)})
		return
	}

	log.Printf("Response from Ditto: %s", string(body))

	// Set response status code
	c.Status(resp.StatusCode)

	// Try to parse as JSON for pretty printing
	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err == nil {
		c.JSON(resp.StatusCode, jsonData)
	} else {
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}

// GetThingState handles getting the current state of a thing
func (h *DeviceHandler) GetThingState(c *gin.Context) {
	thingID := c.Param("thingId")
	if thingID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameter: thingId"})
		return
	}

	// Construct the target URL for Ditto API
	targetURL := fmt.Sprintf("%s/things/%s",
		strings.TrimRight(h.config.Ditto.URL, "/"),
		thingID)

	log.Printf("Getting thing state from: %s", targetURL)

	// Create a new request
	req, err := http.NewRequest(http.MethodGet, targetURL, nil)
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create request: %v", err)})
		return
	}

	// Add Ditto authentication
	req.SetBasicAuth(h.config.Ditto.Username, h.config.Ditto.Password)

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request to Ditto: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to send request: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read response: %v", err)})
		return
	}

	// Parse response
	var thingState ThingState
	if err := json.Unmarshal(body, &thingState); err != nil {
		log.Printf("Failed to parse thing state: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to parse thing state: %v", err)})
		return
	}

	c.JSON(http.StatusOK, thingState)
}

// SendCommand handles sending a command to a thing's feature
func (h *DeviceHandler) SendCommand(c *gin.Context) {
	thingID := c.Param("thingId")
	feature := c.Param("feature")

	if thingID == "" || feature == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required parameters: thingId or feature"})
		return
	}

	// Parse command request
	var cmdReq CommandRequest
	if err := c.ShouldBindJSON(&cmdReq); err != nil {
		log.Printf("Failed to parse command request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Invalid command request: %v", err)})
		return
	}

	// Construct the target URL for Ditto API
	targetURL := fmt.Sprintf("%s/things/%s/features/%s/inbox/messages/%s",
		strings.TrimRight(h.config.Ditto.URL, "/"),
		thingID,
		feature,
		cmdReq.Command)

	log.Printf("Sending command to: %s", targetURL)
	log.Printf("Command data: %+v", cmdReq)

	// Convert params to JSON
	payload, err := json.Marshal(cmdReq.Params)
	if err != nil {
		log.Printf("Failed to marshal command params: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to marshal command params: %v", err)})
		return
	}

	// Create a new request
	req, err := http.NewRequest(http.MethodPut, targetURL, strings.NewReader(string(payload)))
	if err != nil {
		log.Printf("Failed to create request: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create request: %v", err)})
		return
	}

	// Set content type
	req.Header.Set("Content-Type", "application/json")

	// Add Ditto authentication
	req.SetBasicAuth(h.config.Ditto.Username, h.config.Ditto.Password)

	// Create HTTP client and send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Failed to send request to Ditto: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to send request: %v", err)})
		return
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to read response: %v", err)})
		return
	}

	log.Printf("Response from Ditto: %s", string(body))

	// Set response status code
	c.Status(resp.StatusCode)

	// Try to parse as JSON for pretty printing
	var jsonData interface{}
	if err := json.Unmarshal(body, &jsonData); err == nil {
		c.JSON(resp.StatusCode, jsonData)
	} else {
		c.Data(resp.StatusCode, resp.Header.Get("Content-Type"), body)
	}
}
