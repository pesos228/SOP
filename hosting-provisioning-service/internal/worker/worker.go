package worker

import (
	"context"
	"fmt"
	events "hosting-events-contract"
	"hosting-provisioning-service/internal/config"
	"hosting-provisioning-service/internal/messaging"
	"log"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/wagslane/go-rabbitmq"
)

type Worker struct {
	cfg       *config.Config
	conn      *rabbitmq.Conn
	publisher *messaging.EventPublisher
}

func New(cfg *config.Config) (*Worker, error) {
	conn, err := rabbitmq.NewConn(cfg.AMQP_URL, rabbitmq.WithConnectionOptionsLogging)
	if err != nil {
		return nil, fmt.Errorf("failed to create RabbitMQ connection: %w", err)
	}

	publisher, err := messaging.NewEventPublisher(conn, events.EventsExchange)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &Worker{
		cfg:       cfg,
		conn:      conn,
		publisher: publisher,
	}, nil
}

func (w *Worker) Run() {
	var wg sync.WaitGroup
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	provisionHandler := NewCommandHandler(ctx, w.cfg, w.publisher)

	managedConsumer, err := messaging.NewManagedConsumer(
		w.conn,
		w.cfg.ProvisionQueue,
		events.ProvisionRequestKey,
		events.CommandsExchange,
		provisionHandler.Handle,
	)
	if err != nil {
		log.Fatalf("Failed to setup managed consumer: %v", err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := managedConsumer.Run(); err != nil {
			log.Printf("Managed consumer stopped with error: %v", err)
		}
	}()

	log.Println("All consumers started. Waiting for signals to shut down...")

	<-ctx.Done()

	log.Println("Shutting down worker gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	managedConsumer.CloseWithContext(shutdownCtx)

	w.publisher.Close()

	log.Println("Waiting for goroutines to finish...")
	wg.Wait()

	if err := w.conn.Close(); err != nil {
		log.Printf("ERROR: failed to close RabbitMQ connection: %v", err)
	}

	log.Println("All workers have stopped successfully.")
}
