package repositories

import (
	"context"
	"errors"

	"github.com/minq3010/Backend-React-Native-App/models"
	"gorm.io/gorm"
)


type EventRepository struct {
	db *gorm.DB
}

func (r *EventRepository) GetMany(ctx context.Context) ([]*models.Event, error) {
	events := []*models.Event{}

	res := r.db.Model(&models.Event{}).Order("updated_at desc").Find(&events)

	if res.Error != nil {
		return nil, res.Error
	}

	return events, nil
}

func (r *EventRepository) GetOne(ctx context.Context, eventId uint) (*models.Event, error) {
	event := &models.Event{}

	res := r.db.Model(event).Where("id = ?", eventId).First(event)

	if res.Error != nil {
		return nil, res.Error
	}

	return event, nil
}

func (r *EventRepository) CreateOne(ctx context.Context, event *models.Event) (*models.Event, error) {
	res := r.db.Model(event).Create(event)

	if res.Error != nil {
		return nil, res.Error
	}

	return event, nil
}

func (r *EventRepository) UpdateOne(ctx context.Context, eventId uint, updateData map[string]interface{}) (*models.Event, error) {
	event := &models.Event{}

	updateRes := r.db.Model(event).Where("id = ?", eventId).Updates(updateData)

	if updateRes.Error != nil {
		return nil, updateRes.Error
	}

	getRes := r.db.Model(event).Where("id = ?", eventId).First(event)

	if getRes.Error != nil {
		return nil, getRes.Error
	}

	return event, nil
}

func (r *EventRepository) DeleteOne(ctx context.Context, eventId uint) error {
	var count int64

	// Kiểm tra xem có vé nào đã entered = true cho event này không
	err := r.db.Model(&models.Ticket{}).
		Where("event_id = ? AND entered = true", eventId).
		Count(&count).Error

	if err != nil {
		return err
	}

	if count > 0 {
		return errors.New("this event cannot be deleted because tickets have already been confirmed - this is to ensure customer rights are protected")
	}

	// Nếu không có vé entered = true thì xoá
	res := r.db.Delete(&models.Event{}, eventId)
	return res.Error
}

func NewEventRepository(db *gorm.DB) models.EventRepository {
	return &EventRepository{
		db: db,
	}
}