package model

import "time"

type DeliveryStatus string

const (
	DeliveryStatusPending       DeliveryStatus = "menunggu persetujuan"
	DeliveryStatusApproved      DeliveryStatus = "disetujui"
	DeliveryStatusNotApproved   DeliveryStatus = "ditolak"
	DeliveryStatusWaitingDriver DeliveryStatus = "menunggu pengemudi"
	DeliveryStatusOnDelivery    DeliveryStatus = "dalam pengiriman"
	DeliveryStatusFinish        DeliveryStatus = "selesai"
	DeliveryStatusCancelled     DeliveryStatus = "dibatalkan"
)

type Delivery struct {
	ID                   int                    `json:"id" gorm:"primaryKey"`
	TransactionID        int                    `json:"transaction_id"`
	Transaction          Transaction            `json:"transaction" gorm:"foreignKey:TransactionID"`
	DriverID             int                    `json:"driver_id"`
	Driver               Driver                 `json:"driver" gorm:"foreignKey:DriverID"`
	VehicleID            int                    `json:"vehicle_id"`
	Vehicle              Vehicle                `json:"vehicle" gorm:"foreignKey:VehicleID"`
	DeliveryCode         string                 `json:"delivery_code"`
	TotalItem            int                    `json:"total_item"`
	TotalWeight          int                    `json:"total_weight"`
	PickupAddress        string                 `json:"pickup_address"`
	PickupAddressLat     float64                `json:"pickup_address_lat"`
	PickupAddressLang    float64                `json:"pickup_address_lang"`
	DeliveryDate         *time.Time             `json:"delivery_date"`
	DeliveryStatus       DeliveryStatus         `json:"delivery_status"`
	DeliveryCost         float64                `json:"delivery_cost"`
	ApprovedAt           *time.Time             `json:"approved_at"`
	CreatedAt            time.Time              `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt            time.Time              `json:"updated_at" gorm:"autoUpdateTime"`
	Items                []DeliveryItem         `json:"items" gorm:"foreignKey:DeliveryID"`
	DeliveryDestinations []DeliveryDestinations `json:"delivery_destinations" gorm:"foreignKey:DeliveryID"`
}
