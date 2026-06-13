package port

import (
	"context"

	"github.com/DoMinhHHung/Bee/notify-service/internal/domain"
)

type EmailSender interface {
	SendOtpEmail(ctx context.Context, req domain.EmailOtpRequest) error
}
