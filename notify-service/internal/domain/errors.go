package domain

import "errors"

var (
	ErrInvalidEmail    = errors.New("invalid email address")
	ErrInvalidOtp      = errors.New("invalid otp code")
	ErrExpiredOtp      = errors.New("otp has expired")
	ErrFailedToSendEmail = errors.New("failed to send email")
)
