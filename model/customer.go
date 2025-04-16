package model

type Customer struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	Company      string        `json:"company"`
	Email        string        `json:"email"`
	Phone        string        `json:"phone"`
	Address      string        `json:"string"`
	Transactions []Transaction `json:"transactions,omitempty"`
}
