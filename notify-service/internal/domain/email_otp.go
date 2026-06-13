package domain

import (
	"time"
)

type OtpType string

const (
	OtpTypeVerifyEmail    OtpType = "verify_email"
	OtpTypeForgotPassword OtpType = "forgot_password"
)

type EmailOtpRequest struct {
	Email     string    `json:"email"`
	Otp       string    `json:"otp"`
	Type      OtpType   `json:"type"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (r *EmailOtpRequest) IsExpired() bool {
	return time.Now().After(r.ExpiredAt)
}
