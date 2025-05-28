package ditto

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"ditto/internal/influxdb"
)

// Service represents the Ditto service interface
type Service interface {
	Start(ctx context.Context) error
	Stop() error
	GetThing(thingID string) (json.RawMessage, error)
	CreateThing(thingID string, thing json.RawMessage) error
	UpdateThing(thingID string, thing json.RawMessage) error
	DeleteThing(thingID string) error
}

// service implements the Ditto service
type service struct {
	client   *Client
	influxDB *influxdb.Client
}

// NewService creates a new Ditto service
func NewService(client *Client, influxDB *influxdb.Client) Service {
	return &service{
		client:   client,
		influxDB: influxDB,
	}
}

// Start starts the Ditto service
func (s *service) Start(ctx context.Context) error {
	// Connect to Ditto
	if err := s.client.Connect(); err != nil {
		return fmt.Errorf("failed to connect to Ditto: %v", err)
	}

	// Subscribe to all thing events
	filter := "exists(thingId)"
	if err := s.client.Subscribe(filter); err != nil {
		return fmt.Errorf("failed to subscribe to events: %v", err)
	}

	// Start listening for events from Ditto
	go func() {
		if err := s.client.Listen(func(topic string, value json.RawMessage) {
			log.Printf("Received event from Ditto:")
			log.Printf("Topic: %s", topic)
			log.Printf("Content: %s", string(value))

			// Parse the event value
			var event struct {
				ThingID  string `json:"thingId"`
				Features map[string]struct {
					Properties map[string]interface{} `json:"properties"`
				} `json:"features"`
			}

			// Try to parse the value as a direct event
			if err := json.Unmarshal(value, &event); err != nil {
				// If direct parsing fails, try to parse as a wrapped event
				var wrappedEvent struct {
					Value struct {
						ThingID  string `json:"thingId"`
						Features map[string]struct {
							Properties map[string]interface{} `json:"properties"`
						} `json:"features"`
					} `json:"value"`
				}
				if err := json.Unmarshal(value, &wrappedEvent); err != nil {
					log.Printf("Failed to parse event payload: %v", err)
					return
				}
				event = wrappedEvent.Value
			}

			// Process features and store in InfluxDB
			for feature, data := range event.Features {
				if val, ok := data.Properties["value"]; ok {
					// Convert value to float64
					var floatVal float64
					switch v := val.(type) {
					case float64:
						floatVal = v
					case int:
						floatVal = float64(v)
					case int64:
						floatVal = float64(v)
					case string:
						// Try to parse string as float
						if f, err := strconv.ParseFloat(v, 64); err == nil {
							floatVal = f
						} else {
							log.Printf("Failed to parse value as float: %v", err)
							continue
						}
					default:
						log.Printf("Unsupported value type: %T", val)
						continue
					}

					// Get timestamp from properties or use current time
					timestamp := time.Now()
					if ts, ok := data.Properties["timestamp"].(string); ok {
						if t, err := time.Parse(time.RFC3339, ts); err == nil {
							timestamp = t
						}
					}

					// Store in InfluxDB
					err := s.influxDB.WriteEvent(
						event.ThingID,
						feature,
						floatVal,
						timestamp,
					)
					if err != nil {
						log.Printf("Failed to store event in InfluxDB: %v", err)
					}
				}
			}
		}); err != nil {
			log.Printf("Error listening to Ditto events: %v", err)
		}
	}()

	return nil
}

// Stop stops the Ditto service
func (s *service) Stop() error {
	if err := s.client.Close(); err != nil {
		return fmt.Errorf("failed to close Ditto client: %v", err)
	}
	return nil
}

// GetThing retrieves a thing by its ID
func (s *service) GetThing(thingID string) (json.RawMessage, error) {
	return s.client.GetThing(thingID)
}

// CreateThing creates a new thing
func (s *service) CreateThing(thingID string, thing json.RawMessage) error {
	return s.client.CreateThing(thingID, thing)
}

// UpdateThing updates an existing thing
func (s *service) UpdateThing(thingID string, thing json.RawMessage) error {
	return s.client.UpdateThing(thingID, thing)
}

// DeleteThing deletes a thing
func (s *service) DeleteThing(thingID string) error {
	return s.client.DeleteThing(thingID)
}
