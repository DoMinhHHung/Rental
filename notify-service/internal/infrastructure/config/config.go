package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppPort int

	SmtpHost     string
	SmtpPort     int
	SmtpUser     string
	SmtpPass     string
	SmtpFrom     string

	RabbitMQURL  string
	RabbitMQExchange string
	RabbitMQQueue    string
	RabbitMQRoutingKey string
	RabbitMQDLX      string
	RabbitMQDLQ      string
}

func LoadConfig() *Config {
	return &Config{
		AppPort: getEnvAsInt("APP_PORT", 8081),

		SmtpHost: os.Getenv("SMTP_HOST"),
		SmtpPort: getEnvAsInt("SMTP_PORT", 587),
		SmtpUser: os.Getenv("SMTP_USER"),
		SmtpPass: os.Getenv("SMTP_PASS"),
		SmtpFrom: os.Getenv("SMTP_FROM"),

		RabbitMQURL:  os.Getenv("RABBITMQ_URL"),
		RabbitMQExchange: os.Getenv("RABBITMQ_EXCHANGE"),
		RabbitMQQueue:    os.Getenv("RABBITMQ_QUEUE"),
		RabbitMQRoutingKey: os.Getenv("RABBITMQ_ROUTING_KEY"),
		RabbitMQDLX:      os.Getenv("RABBITMQ_DLX"),
		RabbitMQDLQ:      os.Getenv("RABBITMQ_DLQ"),
	}
}

func getEnvAsInt(name string, defaultVal int) int {
	valueStr := os.Getenv(name)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultVal
}
