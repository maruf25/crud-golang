package models

import "time"

type Cart struct {
	Id         int        `gorm:"primaryKey" json:"id"`
	UserId     int        `json:"user_id"`
	ProductId  int        `json:"product_id" binding:"required"`
	Quantity   int        `json:"quantity" binding:"required"`
	TotalPrice float64    `json:"total_price"`
	CreatedAt  *time.Time `json:"created_at,omitempty"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
	Product    Product    `json:"Product,omitempty"`
}

type CreateCart struct {
	ProductId int `json:"product_id" binding:"required"`
	Quantity  int `json:"quantity" `
}
