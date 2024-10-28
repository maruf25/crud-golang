package models

import "time"

type Transaction struct {
	Id         int       `gorm:"primaryKey" json:"id"`
	UserId     string    `json:"user_id"`
	ProductId  string    `json:"product_id"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	User       User
	Product    Product
}
