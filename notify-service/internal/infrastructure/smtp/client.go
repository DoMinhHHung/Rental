package smtp

import (
	"context"
	"fmt"
	"net/smtp"

	"github.com/DoMinhHHung/Bee/notify-service/internal/application/port"
	"github.com/DoMinhHHung/Bee/notify-service/internal/domain"
	"github.com/DoMinhHHung/Bee/notify-service/internal/infrastructure/config"
)

type smtpClient struct {
	config *config.Config
}

func NewSmtpClient(cfg *config.Config) port.EmailSender {
	return &smtpClient{config: cfg}
}

func (s *smtpClient) SendOtpEmail(ctx context.Context, req domain.EmailOtpRequest) error {
	auth := smtp.PlainAuth("", s.config.SmtpUser, s.config.SmtpPass, s.config.SmtpHost)

	to := []string{req.Email}
	msg := []byte(fmt.Sprintf("To: %s\r\n"+
		"Subject: Your OTP Code\r\n"+
		"\r\n"+
		"Your OTP code for %s is: %s. It will expire at %s.\r\n",
		req.Email, req.Type, req.Otp, req.ExpiredAt.Format("2006-01-02 15:04:05")))

	addr := fmt.Sprintf("%s:%d", s.config.SmtpHost, s.config.SmtpPort)
	err := smtp.SendMail(addr, auth, s.config.SmtpFrom, to, msg)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrFailedToSendEmail, err)
	}
	return nil
}
