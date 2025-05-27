package ditto

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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

	// Subscribe to temperature events for specific thing
	thingID := "org.eclipse.ditto:device-1"
	filter := fmt.Sprintf("and(eq(thingId,'%s'),gt(features/temperature/properties/value,50))", thingID)
	if err := s.client.Subscribe(filter); err != nil {
		return fmt.Errorf("failed to subscribe to events: %v", err)
	}

	// Start listening for events from Ditto
	go func() {
		if err := s.client.Listen(func(topic string, value json.RawMessage) {
			log.Printf("Received event from Ditto:")
			log.Printf("Topic: %s", topic)
			log.Printf("Content: %s", string(value))

			// Store event in InfluxDB (tự sinh event từ features)
			var thing struct {
				ThingId  string `json:"thingId"`
				Features map[string]struct {
					Properties map[string]interface{} `json:"properties"`
				} `json:"features"`
			}
			if err := json.Unmarshal(value, &thing); err == nil {
				for feature, data := range thing.Features {
					if val, ok := data.Properties["value"]; ok {
						// Chỉ ghi nếu value là số
						if floatVal, ok := val.(float64); ok {
							err := s.influxDB.WriteEvent(
								thing.ThingId,
								feature,
								int(floatVal),
							)
							if err != nil {
								log.Printf("Failed to store event in InfluxDB: %v", err)
							}
						}
					}
				}
			} else {
				log.Printf("Failed to parse thing payload: %v", err)
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
