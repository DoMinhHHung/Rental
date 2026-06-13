package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/DoMinhHHung/Bee/notify-service/internal/application/port"
	"github.com/DoMinhHHung/Bee/notify-service/internal/domain"
)

type SendEmailOtpUseCase struct {
	emailSender port.EmailSender
	logger      *slog.Logger
}

func NewSendEmailOtpUseCase(emailSender port.EmailSender, logger *slog.Logger) *SendEmailOtpUseCase {
	return &SendEmailOtpUseCase{
		emailSender: emailSender,
		logger:      logger,
	}
}

func (u *SendEmailOtpUseCase) Execute(ctx context.Context, req domain.EmailOtpRequest) error {
	if req.Email == "" {
		return domain.ErrInvalidEmail
	}
	if req.Otp == "" {
		return domain.ErrInvalidOtp
	}
	if req.IsExpired() {
		return domain.ErrExpiredOtp
	}

	u.logger.Info("sending email otp", "email", req.Email, "type", req.Type)

	err := u.emailSender.SendOtpEmail(ctx, req)
	if err != nil {
		return fmt.Errorf("usecase failed: %w", err)
	}

	u.logger.Info("email otp sent successfully", "email", req.Email)
	return nil
}
