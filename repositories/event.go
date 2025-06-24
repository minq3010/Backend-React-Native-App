package repositories

import (
	"context"

	"github.com/minq3010/Backend-React-Native-App/models"
)


type EventRepository struct {
	db any
}

func (r *EventRepository) GetMany(ctx context.Context) ([]*models.Event, error) {
	return nil, nil
}

func (r *EventRepository) GetOne(ctx context.Context, eventId string) (*models.Event, error) {
	return nil, nil
}

func (r *EventRepository) CreateOne(ctx context.Context, event models.Event) (*models.Event, error) {
	return nil, nil
}

func NewRepository(db any) models.EventRepository {
	return &EventRepository {
		db: db,
	}
}