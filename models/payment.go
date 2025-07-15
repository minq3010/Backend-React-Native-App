package models

import "time"

type Payment struct {
	ID            uint   `gorm:"primaryKey"`
	OrderID       string `gorm:"unique;not null"`
	TicketID      string `gorm:"not null"`
	Amount        int    `gorm:"not null"`
	Status        string `gorm:"default:'pending'"` // pending | success | failed
	PaymentMethod string `gorm:"default:'vnpay'"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
