package model

import "time"

type Delivery struct {
	ID                int             `json:"id"`
	TransactionID     int             `json:"transaction_id"`
	DriverID          int             `json:"driver_id"`
	Driver            Driver          `json:"driver"`
	CarID             int             `json:"car_id"`
	Car               Car             `json:"car"`
	TrackingNumber    string          `json:"tracking_number"`
	Content           string          `json:"content"`
	Volume            int16           `json:"volume"`
	AddressFrom       string          `json:"address_from"`
	AddressTo         string          `json:"address_to"`
	DeliveryDate      time.Time       `json:"delivery_date"`
	DeliveryDeadline  time.Time       `json:"delivery_deadline"`
	Total             float64         `json:"total"`
	Status            string          `json:"status"`
	DeliveryReceiptID DeliveryReceipt `json:"delivery_receipt"`
}
