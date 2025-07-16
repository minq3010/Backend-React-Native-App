package models

import (
	"context"
	"time"
)

type Payment struct {
	ID            uint      `gorm:"primaryKey"`
	OrderID       string    `gorm:"unique;not null"`         // Mã giao dịch duy nhất
	TicketID      string    `gorm:"default:null"`            // Có thể null nếu chưa tạo vé
	EventID       uint      `gorm:"not null"`                // Sự kiện liên quan
	UserID        uint      `gorm:"not null"`                // Người thực hiện thanh toán
	Amount        int       `gorm:"not null"`                // Số tiền thanh toán
	Status        string    `gorm:"default:'pending'"`       // pending / success / fail
	PaymentMethod string    `gorm:"default:'vnpay'"`         // vnpay / momo / ...
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
	GetByOrderID(ctx context.Context, orderID string) (*Payment, error)
	UpdateStatus(ctx context.Context, orderID string, status string) error
	UpdateTicketID(ctx context.Context, orderID string, ticketID string) error
}
