package model

import "time"

type Transaction struct {
	ID                int        `json:"id" gorm:"primaryKey"`
	CustomerID        int        `json:"customer_id"`
	Customer          Customer   `json:"customer" gorm:"foreignKey:CustomerID"`
	Delivery          []Delivery `json:"delivery,omitempty" gorm:"foreignKey:TransactionID"`
	TotalDelivery     uint       `json:"total_delivery"`
	Cost              float64    `json:"cost"`
	PaymentDeadline   time.Time  `json:"payment_deadline"`
	DownPayment       float64    `json:"down_payment"`
	DownPaymentStatus string     `json:"down_payment_status"`
	DownPaymentTime   time.Time  `json:"down_payment_time"`
	FullPayment       float64    `json:"full_payment"`
	FullPaymentStatus string     `json:"full_payment_status"`
	FullPaymentTime   *time.Time `json:"full_payment_time"`
	TransactionStatus string     `json:"transaction_status"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
}
