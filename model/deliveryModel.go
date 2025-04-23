package model

import "time"

type DeliveryStatus string

const (
	DeliveryStatusPending     DeliveryStatus = "menunggu persetujuan"
	DeliveryStatusApproved    DeliveryStatus = "disetujui"
	DeliveryStatusNotApproved DeliveryStatus = "ditolak"
	DeliveryStatusOngoing     DeliveryStatus = "sedang berlangsung"
	DeliveryStatusFinish      DeliveryStatus = "selesai"
	DeliveryStatusCancelled   DeliveryStatus = "dibatalkan"
)

type VolumeUnit string

const (
	UnitKilogram VolumeUnit = "kg"
	UnitTon      VolumeUnit = "ton"
	UnitCubicM   VolumeUnit = "m3"
	UnitLiter    VolumeUnit = "liter"
)

type Delivery struct {
	ID                   int            `json:"id" gorm:"primaryKey"`
	TransactionID        int            `json:"transaction_id"`
	DriverID             int            `json:"driver_id"`
	Driver               Driver         `json:"driver"`
	VehicleID            int            `json:"vehicle_id"`
	Vehicle              Vehicle        `json:"vehicle"`
	Content              string         `json:"content"`
	Volume               float64        `json:"volume"`
	VolumeUnit           VolumeUnit     `json:"volume_unit"`
	AddressFrom          string         `json:"address_from"`
	AddressTo            string         `json:"address_to"`
	DeliveryDate         time.Time      `json:"delivery_date"`
	DeliveryDeadlineDate time.Time      `json:"delivery_deadline_date"`
	DeliveryStatus       DeliveryStatus `json:"delivery_status"`
	Total                float64        `json:"total"`
	ApprovedAt           time.Time      `json:"approved_at"`
	CreatedAt            time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
}
