package influxdb

import (
	"context"
	"fmt"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
)

// Client represents an InfluxDB client
type Client struct {
	client   influxdb2.Client
	writeAPI api.WriteAPIBlocking
}

// NewClient creates a new InfluxDB client
func NewClient(url, token, org, bucket string) *Client {
	client := influxdb2.NewClient(url, token)
	writeAPI := client.WriteAPIBlocking(org, bucket)

	return &Client{
		client:   client,
		writeAPI: writeAPI,
	}
}

// WriteEvent writes a WebSocket event to InfluxDB
func (c *Client) WriteEvent(deviceID, eventType string, eventCount int) error {
	point := influxdb2.NewPoint(
		"websocket_events", // measurement
		map[string]string{
			"device_id":  deviceID,
			"event_type": eventType,
		},
		map[string]interface{}{
			"event_count": eventCount,
		},
		time.Now(),
	)

	err := c.writeAPI.WritePoint(context.Background(), point)
	if err != nil {
		return fmt.Errorf("failed to write point: %v", err)
	}

	return nil
}

// Close closes the InfluxDB client
func (c *Client) Close() {
	c.client.Close()
}
