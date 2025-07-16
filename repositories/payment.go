package repositories

import (
	"context"

	"github.com/minq3010/Backend-React-Native-App/models"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{
		db: db,
	}
}

func (r *PaymentRepository)	Create(ctx context.Context, p *models.Payment) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *PaymentRepository) GetMany(ctx context.Context) ([]models.Payment, error) {
	var payments []models.Payment
	err := r.db.WithContext(ctx).Find(&payments).Error
	return payments, err
}

func (r *PaymentRepository) UpdateStatus(ctx context.Context, orderID, status string) error {
	return r.db.WithContext(ctx).
		Model(&models.Payment{}).
		Where("order_id = ?", orderID).
		Update("status", status).
		Error
} 


