package model

import "time"

type OTP struct {
	ID           int       `json:"id" gorm:"primaryKey"`
	Phone        string    `json:"phone"`
	OtpCode      string    `json:"otp_code"`
	OtpExpiredAT time.Time `json:"otp_expired_at"`
	IsVerified   bool      `json:"is_verified"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
