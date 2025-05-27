package service

import (
	"context"
	"fmt"
	"time"

	"ditto/internal/ditto"
	"ditto/internal/model"
	"ditto/internal/repository"
)

// ThingService handles business logic for things
type ThingService struct {
	repo        repository.ThingRepository
	dittoClient *ditto.Client
}

// NewThingService creates a new thing service
func NewThingService(repo repository.ThingRepository, dittoClient *ditto.Client) *ThingService {
	return &ThingService{
		repo:        repo,
		dittoClient: dittoClient,
	}
}

// Create creates a new thing
func (s *ThingService) Create(ctx context.Context, input *model.ThingCreate) (*model.Thing, error) {
	// Create policy if it doesn't exist
	if err := s.dittoClient.CreatePolicy(input.PolicyID); err != nil {
		return nil, fmt.Errorf("failed to create policy: %w", err)
	}

	thing := &model.Thing{
		ID:         input.ID,
		PolicyID:   input.PolicyID,
		Definition: input.Definition,
		Attributes: input.Attributes,
		Features:   input.Features,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := s.repo.Create(ctx, thing); err != nil {
		return nil, err
	}

	return thing, nil
}

// GetByID retrieves a thing by its ID
func (s *ThingService) GetByID(ctx context.Context, id string) (*model.Thing, error) {
	return s.repo.GetByID(ctx, id)
}

// Update updates an existing thing
func (s *ThingService) Update(ctx context.Context, id string, input *model.ThingUpdate) (*model.Thing, error) {
	thing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Update(ctx, id, input); err != nil {
		return nil, err
	}

	// Update the thing object with new values
	if input.PolicyID != "" {
		thing.PolicyID = input.PolicyID
	}
	if input.Definition != "" {
		thing.Definition = input.Definition
	}
	if input.Attributes != nil {
		thing.Attributes = input.Attributes
	}
	if input.Features != nil {
		thing.Features = input.Features
	}
	thing.UpdatedAt = time.Now()

	return thing, nil
}

// Delete removes a thing
func (s *ThingService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// List retrieves a list of things with pagination
func (s *ThingService) List(ctx context.Context, offset, limit int) ([]*model.Thing, error) {
	return s.repo.List(ctx, offset, limit)
}
