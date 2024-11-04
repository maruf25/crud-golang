package models

import "time"

type Product struct {
	Id          int        `gorm:"primaryKey" json:"id"`
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description" binding:"required"`
	Image       string     `json:"image" binding:"required"`
	Stock       int        `json:"stock" binding:"required"`
	Price       float64    `json:"price" binding:"required"`
	CreatedAt   *time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}
