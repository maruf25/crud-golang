package models

import "time"

type TransactionItem struct {
	Id            int        `gorm:"primaryKey" json:"id"`
	ProductId     int        `json:"product_id"`
	TransactionId int        `json:"transaction_id"`
	Quantity      int        `json:"quantity"`
	TotalPrice    float64    `json:"total_price"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
	Product       Product    `json:"Product,omitempty"`
}
