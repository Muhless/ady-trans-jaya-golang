package model

import "time"

type PaymentStatus string
type TransactionStatus string

const (
	PaymentStatusPending PaymentStatus = "pending"
	PaymentStatusPaid    PaymentStatus = "lunas"
	PaymentStatusFailed  PaymentStatus = "gagal"
)

const (
	TransactionStatusWaiting   TransactionStatus = "menunggu persetujuan"
	TransactionStatusOngoing   TransactionStatus = "sedang berlangsung"
	TransactionStatusCompleted TransactionStatus = "selesai"
	TransactionStatusFailed    TransactionStatus = "dibatalkan"
)

type Transaction struct {
	ID                int               `json:"id" gorm:"primaryKey"`
	CustomerID        int               `json:"customer_id"`
	Customer          Customer          `json:"customer"`
	TotalDelivery     uint8             `json:"total_delivery"`
	Total             float64           `json:"total"`
	PaymentDeadline   time.Time         `json:"payment_deadline"`
	DownPayment       float64           `json:"down_payment"`
	DownPaymentStatus PaymentStatus     `json:"down_payment_status"`
	DownPaymentTime   time.Time         `json:"down_payment_time"`
	DownPaymentProof  string            `json:"down_payment_proof"`
	FullPayment       float64           `json:"full_payment"`
	FullPaymentStatus PaymentStatus     `json:"full_payment_status"`
	FullPaymentTime   *time.Time        `json:"full_payment_time"`
	FullPaymentProof  string            `json:"full_payment_proof"`
	TransactionStatus TransactionStatus `json:"transaction_status"`
	DeliveryID        []Delivery        `json:"delivery,omitempty"`
	CreatedAt         time.Time         `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time         `json:"updated_at" gorm:"autoUpdateTime"`
}
