package model

import "time"

type DeliveryReceipt struct {
	ID             int       `json:"id"`
	DeliveryID     int       `json:"delivery_id"`
	ReceiptNumber  string    `json:"receipt_number"`
	IssuedDate     time.Time `json:"issued_date"`
	DriverName     string    `json:"driver_name"`
	CarPlateNumber string    `json:"car_plate_number"`
	ReceiverName   string    `json:"receiver_name"`
	ReceiverNote   string    `json:"receiver_note,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
