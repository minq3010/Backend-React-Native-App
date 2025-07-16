package repositories

import (
	"context"
	"errors"

	"github.com/minq3010/Backend-React-Native-App/models"
	"gorm.io/gorm"
)

type TicketRepository struct {
	db *gorm.DB
}

func (r *TicketRepository) GetMany(ctx context.Context, userId uint) ([]*models.Ticket, error) {
	tickets := []*models.Ticket{}

	res := r.db.Model(&models.Ticket{}).
		Where("user_id = ?", userId).
		Preload("Event").
		Order("updated_at desc").
		Find(&tickets)

	if res.Error != nil {
		return nil, res.Error
	}

	return tickets, nil
}

func (r *TicketRepository) GetOne(ctx context.Context, userId uint, ticketId uint) (*models.Ticket, error) {
	ticket := &models.Ticket{}

	res := r.db.Model(ticket).
		Where("id = ?", ticketId).
		Where("user_id = ?", userId).
		Preload("Event").
		First(ticket)

	if res.Error != nil {
		return nil, res.Error
	}

	return ticket, nil
}

func (r *TicketRepository) CreateOne(ctx context.Context, userId uint, ticket *models.Ticket) (*models.Ticket, error) {
	ticket.UserID = userId

	res := r.db.Model(ticket).Create(ticket)
	if res.Error != nil {
		return nil, res.Error
	}

	return r.GetOne(ctx, userId, ticket.ID)
}

func (r *TicketRepository) UpdateOne(ctx context.Context, userId uint, ticketId uint, updateData map[string]interface{}) (*models.Ticket, error) {
	// Nếu có cập nhật giá, kiểm tra điều kiện
	if price, ok := updateData["price"]; ok {
		var ticket models.Ticket
		if err := r.db.First(&ticket, ticketId).Error; err != nil {
			return nil, err
		}
		if ticket.Entered {
			return nil, errors.New("cannot update price: ticket already entered")
		}

		// ép kiểu an toàn nếu cần
		switch v := price.(type) {
		case float64:
			updateData["price"] = int(v)
		case int:
			// ok
		default:
			return nil, errors.New("invalid price value")
		}
	}

	updateRes := r.db.Model(&models.Ticket{}).
		Where("id = ? AND user_id = ?", ticketId, userId).
		Updates(updateData)

	if updateRes.Error != nil {
		return nil, updateRes.Error
	}

	return r.GetOne(ctx, userId, ticketId)
}

func NewTicketRepository(db *gorm.DB) models.TicketRepository {
	return &TicketRepository{
		db: db,
	}
}