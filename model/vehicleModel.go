package model

type Status string

const (
	Available    Status = "tersedia"
	NotAvailable Status = "tidak tersedia"
)

type Vehicle struct {
	ID          int     `json:"id" gorm:"primaryKey"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	LicensePlat string  `json:"license_plat"`
	Price       float64 `json:"price"`
	Status      Status  `json:"status"`
}
