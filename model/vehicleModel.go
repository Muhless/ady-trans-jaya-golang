package model

import "time"

type Status string

const (
	Available    Status = "tersedia"
	NotAvailable Status = "tidak tersedia"
)

type Vehicle struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name"`
	Type         string    `json:"type"`
	LicensePlate string    `json:"license_plate"`
	Capacity     string    `json:"capacity"`
	RatePerKM    float64   `json:"rate_per_km"`
	Status       Status    `json:"status"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
