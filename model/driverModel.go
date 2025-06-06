package model

import "time"

type DriverStatus string

const (
	DriverAvailable    DriverStatus = "tersedia"
	DriverNotAvailable DriverStatus = "tidak tersedia"
)

type Driver struct {
	ID         int          `json:"id" gorm:"primaryKey"`
	Name       string       `json:"name"`
	Photo      string       `json:"photo"`
	Phone      string       `json:"phone"`
	Address    string       `json:"address"`
	Status     DriverStatus `json:"status"`
	Username   string       `json:"username"`
	Password   string       `json:"password"`
	CreatedAt  time.Time    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time    `json:"updated_at" gorm:"autoUpdateTime"`
	Deliveries []Delivery   `json:"deliveries,omitempty" gorm:"foreignKey:DriverID"`
}
