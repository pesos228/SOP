package app

import (
	"context"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"hosting-provisioning-service/internal/config"
	"hosting-provisioning-service/internal/service"
)

type App struct {
	config *config.Config
}

func New(config *config.Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) Run() {
	var wg sync.WaitGroup

	rabbit, err := a.initRabbitManager(&wg)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ manager: %v", err)
	}

	provisioningService := service.NewProvisioningService(rabbit)

	a.runConsumers(rabbit, provisioningService)

	shutdownCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-shutdownCtx.Done()
	log.Println("Shutting down application gracefully...")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rabbit.Stop(timeoutCtx)

	wg.Wait()
}
