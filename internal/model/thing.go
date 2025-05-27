package model

import "time"

// Thing represents a thing in the system
type Thing struct {
	ID         string                 `json:"thingId"`
	PolicyID   string                 `json:"policyId"`
	Definition string                 `json:"definition,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
	Features   map[string]Feature     `json:"features"`
	CreatedAt  time.Time              `json:"createdAt,omitempty"`
	UpdatedAt  time.Time              `json:"updatedAt,omitempty"`
}

// Feature represents a feature of a thing
type Feature struct {
	Definition        []string               `json:"definition,omitempty"`
	Properties        map[string]interface{} `json:"properties"`
	DesiredProperties map[string]interface{} `json:"desiredProperties,omitempty"`
}

// ThingCreate represents the data needed to create a thing
type ThingCreate struct {
	ID         string                 `json:"thingId" validate:"required"`
	PolicyID   string                 `json:"policyId" validate:"required"`
	Definition string                 `json:"definition,omitempty"`
	Attributes map[string]interface{} `json:"attributes"`
	Features   map[string]Feature     `json:"features,omitempty"`
}

// ThingUpdate represents the data needed to update a thing
type ThingUpdate struct {
	PolicyID   string                 `json:"policyId,omitempty"`
	Definition string                 `json:"definition,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
	Features   map[string]Feature     `json:"features,omitempty"`
}
