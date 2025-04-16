package model

type Car struct {
	ID           int    `json:"ID"`
	Name         string `json:"name"`
	Type         string `json:"type"`
	LicensePlate string `json:"license_plate"`
	Price        int    `json:"price"`
	Status       bool   `json:"status"`
}
