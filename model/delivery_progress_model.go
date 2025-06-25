package model

import "time"

type DeliveryProgress struct {
	ID         int      `json:"id" gorm:"primaryKey"`
	DeliveryID int      `json:"delivery_id"`
	Delivery   Delivery `json:"delivery" gorm:"foreignKey:DeliveryID"`

	PickupTime        *time.Time `json:"pickup_time"`
	ArrivalTime       *time.Time `json:"destination_arrival_time"`
	ReceiverName      string     `json:"receiver_name"`
	ReceiverPhone     string     `json:"receiver_phone"`
	ReceivedAt        *time.Time `json:"received_at"`
	ReceiverSignature string     `json:"receiver_signature_url"`
	DeliveryPhotoURL  string     `json:"delivery_photo_url"`

	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
