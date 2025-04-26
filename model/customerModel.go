package model

import "time"

type Customer struct {
	ID           int           `json:"id" gorm:"primaryKey"`
	Name         string        `json:"name"`
	Company      string        `json:"company"`
	Email        string        `json:"email"`
	Phone        string        `json:"phone"`
	Address      string        `json:"address"`
	CreatedAt    time.Time     `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time     `json:"updated_at" gorm:"autoUpdateTime"`
	Transactions []Transaction `json:"transactions,omitempty"`
}
