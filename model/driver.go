package model

type Driver struct {
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
	Status  bool   `json:"status"`
	Photo   string `json:"photo"`
}
