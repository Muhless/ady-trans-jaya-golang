package model

type Status string

const (
	Available    Status = "tersedia"
	NotAvailable Status = "tidak tersedia"
)

type Car struct {
	ID           int     `json:"ID"`
	Name         string  `json:"name"`
	Type         string  `json:"type"`
	LicensePlate string  `json:"license_plate"`
	Price        float64 `json:"price"`
	Status       Status  `json:"status"`
}
