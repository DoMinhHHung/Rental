package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/DoMinhHHung/Rental/notify-service/internal/interface/di"
)

func main() {
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to initialize container: %v", err)
	}

	// Setup RabbitMQ Topology
	if err := container.RabbitConsumer.SetupTopology(); err != nil {
		container.Logger.Error("Failed to setup RabbitMQ topology", "error", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start RabbitMQ Consumer
	err = container.RabbitConsumer.Listen(ctx, container.UseCase.Execute)
	if err != nil {
		container.Logger.Error("Failed to start RabbitMQ listener", "error", err)
		os.Exit(1)
	}

	// Start HTTP Server for Health Check
	mux := http.NewServeMux()
	mux.HandleFunc("/health", container.HealthHandler.HealthCheck)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", container.Config.AppPort),
		Handler: mux,
	}

	go func() {
		container.Logger.Info("Starting HTTP server", "port", container.Config.AppPort)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			container.Logger.Error("HTTP server failed", "error", err)
		}
	}()

	// Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop
	container.Logger.Info("Shutting down gracefully...")

	cancel() // Stop consumer

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		container.Logger.Error("HTTP server shutdown failed", "error", err)
	}

	container.RabbitConsumer.Close()
	container.Logger.Info("Service stopped")
}
