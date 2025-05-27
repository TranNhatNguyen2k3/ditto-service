package repository

import (
	"context"

	"ditto/internal/model"
)

// ThingRepository defines the interface for thing data operations
type ThingRepository interface {
	Create(ctx context.Context, thing *model.Thing) error
	GetByID(ctx context.Context, id string) (*model.Thing, error)
	Update(ctx context.Context, id string, thing *model.ThingUpdate) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, offset, limit int) ([]*model.Thing, error)
}
