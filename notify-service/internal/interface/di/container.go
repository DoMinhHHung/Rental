package di

import (
	"log/slog"

	"github.com/DoMinhHHung/Bee/notify-service/internal/application/port"
	"github.com/DoMinhHHung/Bee/notify-service/internal/application/usecase"
	"github.com/DoMinhHHung/Bee/notify-service/internal/infrastructure/config"
	"github.com/DoMinhHHung/Bee/notify-service/internal/infrastructure/logger"
	"github.com/DoMinhHHung/Bee/notify-service/internal/infrastructure/rabbitmq"
	"github.com/DoMinhHHung/Bee/notify-service/internal/infrastructure/smtp"
	"github.com/DoMinhHHung/Bee/notify-service/internal/interface/http"
)

type Container struct {
	Config         *config.Config
	Logger         *slog.Logger
	EmailSender    port.EmailSender
	UseCase        *usecase.SendEmailOtpUseCase
	RabbitConsumer *rabbitmq.Consumer
	HealthHandler  *http.HealthHandler
}

func NewContainer() (*Container, error) {
	cfg := config.LoadConfig()
	log := logger.NewLogger()

	smtpClient := smtp.NewSmtpClient(cfg)
	sendEmailOtpUseCase := usecase.NewSendEmailOtpUseCase(smtpClient, log)

	rabbitConsumer, err := rabbitmq.NewConsumer(cfg, log)
	if err != nil {
		return nil, err
	}

	healthHandler := http.NewHealthHandler()

	return &Container{
		Config:         cfg,
		Logger:         log,
		EmailSender:    smtpClient,
		UseCase:        sendEmailOtpUseCase,
		RabbitConsumer: rabbitConsumer,
		HealthHandler:  healthHandler,
	}, nil
}
