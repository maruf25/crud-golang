package models

import "time"

type Status string

const (
	Success Status = "success"
	Pending Status = "pending"
	Failed  Status = "failed"
)

type Transaction struct {
	Id              int               `gorm:"primaryKey" json:"id"`
	UserId          int               `json:"user_id"`
	ShippingAddress string            `json:"shipping_address" binding:"required"`
	TotalPrice      float64           `json:"total_price"`
	PaymentStatus   Status            `gorm:"type:enum('success', 'pending','failed')" json:"payment_status" binding:"required"`
	CreatedAt       *time.Time        `json:"created_at,omitempty"`
	UpdatedAt       *time.Time        `json:"updated_at,omitempty"`
	User            User              `json:"User,omitempty"`
	TransactionItem []TransactionItem `json:"TransactionItem,omitempty"`
}

type TransactionReq struct {
	ShippingAddress string `json:"shipping_address" binding:"required"`
}
