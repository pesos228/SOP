package app

import (
	"context"
	"fmt"
	events "hosting-events-contract"
	"hosting-service/internal/listeners"
	"hosting-service/internal/messaging"
	"hosting-service/internal/service"
	"log"
	"sync"
	"time"

	"github.com/wagslane/go-rabbitmq"
)

func (a *App) initRabbitMQ() (*rabbitmq.Conn, error) {
	conn, err := rabbitmq.NewConn(a.config.AMQP_URL, rabbitmq.WithConnectionOptionsLogging)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}
	return conn, nil
}

func (a *App) runMessageConsumers(ctx context.Context, wg *sync.WaitGroup, conn *rabbitmq.Conn, serverService service.ServerService) {
	resultsListener := listeners.NewProvisioningResultListener(serverService)
	resultsConsumer, err := messaging.NewManagedConsumer(
		conn,
		a.config.ResultsQueue,
		events.ProvisionResultKeyPattern,
		events.EventsExchange,
		resultsListener.Handle,
	)
	if err != nil {
		log.Fatalf("Failed to create results consumer: %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Println("Starting results listener...")
		if err := resultsConsumer.Run(); err != nil {
			log.Printf("Results listener stopped: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		log.Println("Closing RabbitMQ consumer...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resultsConsumer.CloseWithContext(shutdownCtx)
	}()

}
