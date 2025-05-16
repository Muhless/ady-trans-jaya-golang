package model

import "time"

type DriverStatus string

const (
	DriverAvailable    DriverStatus = "tersedia"
	DriverNotAvailable DriverStatus = "tidak tersedia"
)

type Driver struct {
	ID         int          `json:"id" gorm:"primaryKey"`
	Name       string       `json:"name" form:"name"`
	Photo      string       `json:"photo"`
	Phone      string       `json:"phone" form:"phone"`
	Address    string       `json:"address" form:"address"`
	Status     DriverStatus `json:"status" form:"status"`
	CreatedAt  time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
	Deliveries []Delivery   `json:"deliveries,omitempty" gorm:"foreignKey:DriverID"`
}
