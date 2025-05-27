package ditto

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/websocket"
)

// Client represents a Ditto WebSocket client
type Client struct {
	conn     *websocket.Conn
	host     string
	username string
	password string
}

// NewClient creates a new Ditto WebSocket client
func NewClient(host, username, password string) *Client {
	return &Client{
		host:     host,
		username: username,
		password: password,
	}
}

// Connect establishes a WebSocket connection to Ditto
func (c *Client) Connect() error {
	u := url.URL{Scheme: "ws", Host: c.host, Path: "/ws/2"}

	// Basic Auth header
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
	header := http.Header{}
	header.Add("Authorization", "Basic "+auth)

	log.Printf("Connecting to Ditto WebSocket at %s...", u.String())

	// Connect WebSocket
	conn, resp, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		if resp != nil {
			return fmt.Errorf("cannot connect to WebSocket: %v (HTTP status %s)", err, resp.Status)
		}
		return fmt.Errorf("cannot connect to WebSocket: %v", err)
	}

	c.conn = conn
	log.Printf("Successfully connected to Ditto WebSocket: %s", u.String())
	return nil
}

// Subscribe subscribes to events with the given filter
func (c *Client) Subscribe(filter string) error {
	if c.conn == nil {
		return fmt.Errorf("not connected to Ditto")
	}

	msg := "START-SEND-EVENTS?filter=" + url.QueryEscape(filter)
	log.Printf("Subscribing to events with filter: %s", filter)

	if err := c.conn.WriteMessage(websocket.TextMessage, []byte(msg)); err != nil {
		return fmt.Errorf("failed to send subscription message: %v", err)
	}

	log.Printf("Successfully subscribed to events")
	return nil
}

// SendMessage sends a message to Ditto
func (c *Client) SendMessage(topic string, value json.RawMessage) error {
	if c.conn == nil {
		return fmt.Errorf("not connected to Ditto")
	}

	msg := struct {
		Topic string          `json:"topic"`
		Value json.RawMessage `json:"value"`
	}{
		Topic: topic,
		Value: value,
	}

	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %v", err)
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, msgBytes); err != nil {
		return fmt.Errorf("failed to send message: %v", err)
	}

	log.Printf("Successfully sent message to topic: %s", topic)
	return nil
}

// Listen starts listening for events
func (c *Client) Listen(handler func(topic string, value json.RawMessage)) error {
	if c.conn == nil {
		return fmt.Errorf("not connected to Ditto")
	}

	log.Printf("Starting to listen for events...")

	for {
		messageType, msgBytes, err := c.conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("read error: %v", err)
		}

		// Log raw message for debugging
		log.Printf("Received message type: %d", messageType)
		log.Printf("Raw message: %s", string(msgBytes))

		// Skip non-text messages
		if messageType != websocket.TextMessage {
			log.Printf("Skipping non-text message")
			continue
		}

		// Check if message is a control message
		if strings.HasPrefix(string(msgBytes), "START-SEND-EVENTS") {
			log.Printf("Received subscription confirmation: %s", string(msgBytes))
			continue
		}

		var msg struct {
			Topic string          `json:"topic"`
			Value json.RawMessage `json:"value"`
		}
		if err := json.Unmarshal(msgBytes, &msg); err != nil {
			log.Printf("Failed to parse message as JSON: %v", err)
			log.Printf("Message content: %s", string(msgBytes))
			continue
		}

		if !strings.HasSuffix(msg.Topic, "/things/twin/events/merged") {
			log.Printf("Processing non-merged event with topic: %s", msg.Topic)
		}
		handler(msg.Topic, msg.Value)
	}
}

// Close closes the WebSocket connection
func (c *Client) Close() error {
	if c.conn != nil {
		log.Printf("Closing Ditto WebSocket connection...")
		err := c.conn.Close()
		c.conn = nil
		return err
	}
	return nil
}

// GetThing retrieves a thing by its ID
func (c *Client) GetThing(thingID string) (json.RawMessage, error) {
	url := fmt.Sprintf("http://%s/api/2/things/%s", c.host, thingID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Add Basic Auth
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
	req.Header.Add("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get thing: status code %d", resp.StatusCode)
	}

	var result json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %v", err)
	}

	return result, nil
}

// CreateThing creates a new thing
func (c *Client) CreateThing(thingID string, thing json.RawMessage) error {
	url := fmt.Sprintf("http://%s/api/2/things/%s", c.host, thingID)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(thing))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add Basic Auth
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to create thing: status code %d", resp.StatusCode)
	}

	return nil
}

// UpdateThing updates an existing thing
func (c *Client) UpdateThing(thingID string, thing json.RawMessage) error {
	url := fmt.Sprintf("http://%s/api/2/things/%s", c.host, thingID)

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(thing))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add Basic Auth
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to update thing: status code %d", resp.StatusCode)
	}

	return nil
}

// DeleteThing deletes a thing
func (c *Client) DeleteThing(thingID string) error {
	url := fmt.Sprintf("http://%s/api/2/things/%s", c.host, thingID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add Basic Auth
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
	req.Header.Add("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to delete thing: status code %d", resp.StatusCode)
	}

	return nil
}

// CreatePolicy creates a new policy in Ditto
func (c *Client) CreatePolicy(policyID string) error {
	url := fmt.Sprintf("http://%s/api/2/policies/%s", c.host, policyID)

	// Create default policy
	policy := map[string]interface{}{
		"entries": map[string]interface{}{
			"owner": map[string]interface{}{
				"subjects": map[string]interface{}{
					"nginx:ditto": map[string]interface{}{
						"type": "nginx basic auth user",
					},
				},
				"resources": map[string]interface{}{
					"thing:/": map[string]interface{}{
						"grant":  []string{"READ", "WRITE"},
						"revoke": []string{},
					},
					"policy:/": map[string]interface{}{
						"grant":  []string{"READ", "WRITE"},
						"revoke": []string{},
					},
					"message:/": map[string]interface{}{
						"grant":  []string{"READ", "WRITE"},
						"revoke": []string{},
					},
				},
			},
		},
	}

	// Marshal policy to JSON
	policyJSON, err := json.Marshal(policy)
	if err != nil {
		return fmt.Errorf("failed to marshal policy: %v", err)
	}

	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(policyJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add Basic Auth
	auth := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", c.username, c.password)))
	req.Header.Add("Authorization", "Basic "+auth)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("failed to create policy: status code %d", resp.StatusCode)
	}

	return nil
}
