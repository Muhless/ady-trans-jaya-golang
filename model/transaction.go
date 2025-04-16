package model

import "time"

type PaymentStatus string

const (
	StatusPending PaymentStatus = "pending"
	StatusPaid    PaymentStatus = "lunas"
	StatusFaied   PaymentStatus = "gagal"
)

type Transaction struct {
	ID                int           `json:"id"`
	CustomerID        int           `json:"customer_id"`
	TotalDelivery     uint8         `json:"total_delivery"`
	Total             float64       `json:"total"`
	PaymentDeadline   time.Time     `json:"payment_deadline"`
	DownPayment       float64       `json:"down_payment"`
	DownPaymentTime   time.Time     `json:"down_payment_time"`
	DownPaymentStatus PaymentStatus `json:"down_payment_status"`
	DownPaymentProof  string        `json:"down_payment_proof"`
	FullPayment       float64       `json:"full_payment"`
	FullPaymentStatus PaymentStatus `json:"full_payment_status"`
	FullPaymentProof  string        `json:"full_payment_proof"`
	FullPaymentTime   *time.Time    `json:"full_payment_time"`
	Status            string        `json:"status"`
	DeliveryID        []Delivery    `json:"delivery,omitempty"`
}
