package model

type Status string

const (
	Available    Status = "tersedia"
	NotAvailable Status = "tidak tersedia"
)

type Car struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	LicensePlat string  `json:"license_plat"`
	Price       float64 `json:"price"`
	Status      Status  `json:"status"`
}
