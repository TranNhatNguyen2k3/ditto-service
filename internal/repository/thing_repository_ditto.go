package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"ditto/internal/ditto"
	"ditto/internal/model"
)

// ThingRepositoryDitto implements ThingRepository interface using Ditto service
type ThingRepositoryDitto struct {
	dittoService ditto.Service
	client       *ditto.Client
}

// NewThingRepositoryDitto creates a new instance of ThingRepositoryDitto
func NewThingRepositoryDitto(dittoService ditto.Service, client *ditto.Client) *ThingRepositoryDitto {
	return &ThingRepositoryDitto{
		dittoService: dittoService,
		client:       client,
	}
}

// Create implements ThingRepository
func (r *ThingRepositoryDitto) Create(ctx context.Context, thing *model.Thing) error {
	// Convert thing to Ditto format
	dittoThing := map[string]interface{}{
		"thingId":  thing.ID,
		"policyId": thing.PolicyID,
	}

	if thing.Definition != "" {
		dittoThing["definition"] = thing.Definition
	}

	if thing.Attributes != nil {
		dittoThing["attributes"] = thing.Attributes
	}

	if thing.Features != nil {
		dittoThing["features"] = thing.Features
	}

	// Marshal to JSON
	thingJSON, err := json.Marshal(dittoThing)
	if err != nil {
		return fmt.Errorf("failed to marshal thing: %w", err)
	}

	// Create thing in Ditto
	err = r.dittoService.CreateThing(thing.ID, thingJSON)
	if err != nil {
		return fmt.Errorf("failed to create thing in Ditto: %w", err)
	}

	return nil
}

// GetByID implements ThingRepository
func (r *ThingRepositoryDitto) GetByID(ctx context.Context, id string) (*model.Thing, error) {
	// Get thing from Ditto
	dittoThingJSON, err := r.dittoService.GetThing(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get thing from Ditto: %w", err)
	}

	// Unmarshal JSON
	var dittoThing map[string]interface{}
	if err := json.Unmarshal(dittoThingJSON, &dittoThing); err != nil {
		return nil, fmt.Errorf("failed to unmarshal thing: %w", err)
	}

	// Convert Ditto thing to model.Thing
	thing := &model.Thing{
		ID: id,
	}

	// Extract fields
	if policyID, ok := dittoThing["policyId"].(string); ok {
		thing.PolicyID = policyID
	}
	if definition, ok := dittoThing["definition"].(string); ok {
		thing.Definition = definition
	}
	if attrs, ok := dittoThing["attributes"].(map[string]interface{}); ok {
		thing.Attributes = attrs
	}
	if features, ok := dittoThing["features"].(map[string]interface{}); ok {
		thing.Features = make(map[string]model.Feature)
		for name, feature := range features {
			if featureMap, ok := feature.(map[string]interface{}); ok {
				f := model.Feature{}

				if def, ok := featureMap["definition"].([]interface{}); ok {
					f.Definition = make([]string, len(def))
					for i, d := range def {
						f.Definition[i] = d.(string)
					}
				}

				if props, ok := featureMap["properties"].(map[string]interface{}); ok {
					f.Properties = props
				}

				if desiredProps, ok := featureMap["desiredProperties"].(map[string]interface{}); ok {
					f.DesiredProperties = desiredProps
				}

				thing.Features[name] = f
			}
		}
	}

	return thing, nil
}

// Update implements ThingRepository
func (r *ThingRepositoryDitto) Update(ctx context.Context, id string, thing *model.ThingUpdate) error {
	// Create update payload
	updatePayload := make(map[string]interface{})

	// Update fields if provided
	if thing.PolicyID != "" {
		updatePayload["policyId"] = thing.PolicyID
	}
	if thing.Definition != "" {
		updatePayload["definition"] = thing.Definition
	}
	if thing.Attributes != nil {
		updatePayload["attributes"] = thing.Attributes
	}
	if thing.Features != nil {
		updatePayload["features"] = thing.Features
	}

	// Marshal to JSON
	updateJSON, err := json.Marshal(updatePayload)
	if err != nil {
		return fmt.Errorf("failed to marshal update: %w", err)
	}

	// Update thing in Ditto
	err = r.dittoService.UpdateThing(id, updateJSON)
	if err != nil {
		return fmt.Errorf("failed to update thing in Ditto: %w", err)
	}

	return nil
}

// Delete implements ThingRepository
func (r *ThingRepositoryDitto) Delete(ctx context.Context, id string) error {
	// Delete thing from Ditto
	err := r.dittoService.DeleteThing(id)
	if err != nil {
		return fmt.Errorf("failed to delete thing from Ditto: %w", err)
	}

	return nil
}

// List implements ThingRepository
func (r *ThingRepositoryDitto) List(ctx context.Context, offset, limit int) ([]*model.Thing, error) {
	// Create search query
	query := map[string]interface{}{
		"filter": "exists(thingId)",
		"options": map[string]interface{}{
			"size":   limit,
			"offset": offset,
		},
	}

	// Marshal query to JSON
	queryJSON, err := json.Marshal(query)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal search query: %w", err)
	}

	// Create HTTP request
	url := "http://localhost:8080/api/2/search/things"
	req, err := http.NewRequest("POST", url, strings.NewReader(string(queryJSON)))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Basic ZGl0dG86ZGl0dG8=") // base64 encoded "ditto:ditto"

	// Send request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("search request failed with status: %d", resp.StatusCode)
	}

	// Parse response
	var result struct {
		Items []map[string]interface{} `json:"items"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert to model.Thing slice
	things := make([]*model.Thing, 0, len(result.Items))
	for _, item := range result.Items {
		thing := &model.Thing{
			ID: item["thingId"].(string),
		}

		if policyID, ok := item["policyId"].(string); ok {
			thing.PolicyID = policyID
		}
		if definition, ok := item["definition"].(string); ok {
			thing.Definition = definition
		}
		if attrs, ok := item["attributes"].(map[string]interface{}); ok {
			thing.Attributes = attrs
		}
		if features, ok := item["features"].(map[string]interface{}); ok {
			thing.Features = make(map[string]model.Feature)
			for name, feature := range features {
				if featureMap, ok := feature.(map[string]interface{}); ok {
					f := model.Feature{}

					if def, ok := featureMap["definition"].([]interface{}); ok {
						f.Definition = make([]string, len(def))
						for i, d := range def {
							f.Definition[i] = d.(string)
						}
					}

					if props, ok := featureMap["properties"].(map[string]interface{}); ok {
						f.Properties = props
					}

					if desiredProps, ok := featureMap["desiredProperties"].(map[string]interface{}); ok {
						f.DesiredProperties = desiredProps
					}

					thing.Features[name] = f
				}
			}
		}

		things = append(things, thing)
	}

	return things, nil
}
