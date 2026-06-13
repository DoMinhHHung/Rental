package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/DoMinhHHung/Bee/notify-service/internal/domain"
	"github.com/DoMinhHHung/Bee/notify-service/internal/infrastructure/config"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	config  *config.Config
	logger  *slog.Logger
	handler func(ctx context.Context, req domain.EmailOtpRequest) error
}

func NewConsumer(cfg *config.Config, logger *slog.Logger) (*Consumer, error) {
	conn, err := amqp.Dial(cfg.RabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}

	return &Consumer{
		conn:    conn,
		channel: ch,
		config:  cfg,
		logger:  logger,
	}, nil
}

func (c *Consumer) SetupTopology() error {
	// Declare DLX
	err := c.channel.ExchangeDeclare(
		c.config.RabbitMQDLX,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare DLQ
	_, err = c.channel.QueueDeclare(
		c.config.RabbitMQDLQ,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Bind DLQ to DLX
	err = c.channel.QueueBind(
		c.config.RabbitMQDLQ,
		"email.otp.dead",
		c.config.RabbitMQDLX,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare Main Exchange
	err = c.channel.ExchangeDeclare(
		c.config.RabbitMQExchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare Main Queue with DLX configuration
	args := amqp.Table{
		"x-dead-letter-exchange":    c.config.RabbitMQDLX,
		"x-dead-letter-routing-key": "email.otp.dead",
	}
	_, err = c.channel.QueueDeclare(
		c.config.RabbitMQQueue,
		true,
		false,
		false,
		false,
		args,
	)
	if err != nil {
		return err
	}

	// Bind Main Queue to Main Exchange
	err = c.channel.QueueBind(
		c.config.RabbitMQQueue,
		c.config.RabbitMQRoutingKey,
		c.config.RabbitMQExchange,
		false,
		nil,
	)
	return err
}

func (c *Consumer) Listen(ctx context.Context, handler func(ctx context.Context, req domain.EmailOtpRequest) error) error {
	c.handler = handler

	msgs, err := c.channel.Consume(
		c.config.RabbitMQQueue,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case d, ok := <-msgs:
				if !ok {
					return
				}
				c.processMessage(ctx, d)
			}
		}
	}()

	return nil
}

func (c *Consumer) processMessage(ctx context.Context, d amqp.Delivery) {
	var req domain.EmailOtpRequest
	if err := json.Unmarshal(d.Body, &req); err != nil {
		c.logger.Error("failed to unmarshal message", "error", err)
		d.Nack(false, false) // don't requeue, send to DLQ
		return
	}

	c.logger.Info("processing email otp request", "email", req.Email, "type", req.Type)

	err := c.handler(ctx, req)
	if err != nil {
		c.logger.Error("failed to handle message", "error", err, "email", req.Email)

		// Simple retry logic: if x-death count < 3, requeue. Otherwise, send to DLQ.
		// For simplicity in Phase 1, we can use a basic Nack without requeue for fatal errors,
		// and RabbitMQ's DLX will handle it.
		// If we want actual retries with delay, we'd need a more complex setup.
		// Here we'll just nack without requeue if it fails, which moves it to DLQ.
		d.Nack(false, false)
		return
	}

	d.Ack(false)
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}
