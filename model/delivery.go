package model

import "time"

type DeliveryStatus string

const (
	DeliveryStatusPending   DeliveryStatus = "menunggu persetujuan"
	DeliveryStatusOngoing   DeliveryStatus = "sedang berlangsung"
	DeliveryStatusDelivered DeliveryStatus = "dalam perjalanan"
	DeliveryStatusCancelled DeliveryStatus = "dibatalkan"
)

type VolumeUnit string

const (
	UnitKilogram VolumeUnit = "kg"
	UnitTon      VolumeUnit = "ton"
	UnitCubicM   VolumeUnit = "m3"
	UnitLiter    VolumeUnit = "liter"
)

type Delivery struct {
	ID                int             `json:"id"`
	TransactionID     int             `json:"transaction_id"`
	DriverID          int             `json:"driver_id"`
	Driver            Driver          `json:"driver"`
	CarID             int             `json:"car_id"`
	Car               Car             `json:"car"`
	Content           string          `json:"content"`
	Volume            float64         `json:"volume"`
	VolumeUnit        VolumeUnit      `json:"volume_unit"`
	AddressFrom       string          `json:"address_from"`
	AddressTo         string          `json:"address_to"`
	DeliveryDate      time.Time       `json:"delivery_date"`
	DeliveryDeadline  time.Time       `json:"delivery_deadline"`
	Total             float64         `json:"total"`
	Status            string          `json:"status"`
	// DeliveryReceiptID DeliveryReceipt `json:"delivery_receipt"`
	CreatedAt         time.Time       `json:"created_at"`
	UpdatedAt         time.Time       `json:"updated_at"`
}
