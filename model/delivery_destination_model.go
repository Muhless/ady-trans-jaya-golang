package model

import "time"

type DeliveryDestinations struct {
	ID              int        `json:"id" gorm:"primaryKey"`
	DeliveryID      int        `json:"delivery_id"`
	Address         string     `json:"address"`
	Lat             float64    `json:"lat"`
	Lng             float64    `json:"lng"`
	ArrivalTime     *time.Time `json:"arrival_time"`
	ArrivalPhotoURL string     `json:"arrival_photo_url"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
