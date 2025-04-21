package model

type DriverStatus string

const (
	DriverAvailable    DriverStatus = "tersedia"
	DriverNotAvailable DriverStatus = "tidak tersedia"
)

type Driver struct {
	ID      int          `json:"id"`
	Name    string       `json:"name"`
	Phone   string       `json:"phone"`
	Address string       `json:"address"`
	Status  DriverStatus `json:"status"`
}
