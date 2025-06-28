package model

import "time"

type DeliveryProgress struct {
	ID                int        `json:"id" gorm:"primaryKey"`
	DeliveryID        int        `json:"delivery_id"`
	DeliveryStartTime *time.Time `json:"delivery_start_time"`
	PickupTime        *time.Time `json:"pickup_time"`
	PickupPhotoURL    string     `json:"pickup_photo_url"`
	ArrivalTime       *time.Time `json:"arrival_time"`
	ArrivalPhotoURL   string     `json:"arrival_photo_url"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
