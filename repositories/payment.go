package repositories

import (
	"context"

	"github.com/minq3010/Backend-React-Native-App/models"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) models.PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *models.Payment) error {
	return r.db.WithContext(ctx).Create(payment).Error
}

func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID string) (*models.Payment, error) {
	var payment models.Payment
	err := r.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error
	return &payment, err
}


func (r *PaymentRepository) UpdateStatus(ctx context.Context, orderID string, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.Payment{}).
		Where("order_id = ?", orderID).
		Update("status", status).Error
}

func (r *PaymentRepository) UpdateTicketID(ctx context.Context, orderID string, ticketID string) error {
	return r.db.Model(&models.Payment{}).
		Where("order_id = ?", orderID).
		Update("ticket_id", ticketID).Error
}
