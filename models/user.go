package models

import "time"

type Role string

const (
	Admin  Role = "admin"
	Member Role = "member"
)

type User struct {
	Id          int           `gorm:"primaryKey" json:"id"`
	Name        string        `json:"name" binding:"required"`
	Email       string        `json:"email,omitempty" binding:"required,email" gorm:"unique;not null"`
	Password    string        `json:"password,omitempty" binding:"required"`
	Role        Role          `gorm:"type:enum('admin', 'member')" json:"role,omitempty" binding:"required"`
	CreatedAt   *time.Time    `json:"created_at,omitempty"`
	UpdatedAt   *time.Time    `json:"updated_at,omitempty"`
	Transaction []Transaction `json:"Transaction,omitempty"`
	Cart        []Cart        `json:"Cart,omitempty"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
