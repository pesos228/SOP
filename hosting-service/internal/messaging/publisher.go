package messaging

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/wagslane/go-rabbitmq"
)

type EventPublisher struct {
	publisher *rabbitmq.Publisher
	exchange  string
}

func NewEventPublisher(conn *rabbitmq.Conn, exchangeName string) (*EventPublisher, error) {
	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsLogging,
		rabbitmq.WithPublisherOptionsExchangeName(exchangeName),
		rabbitmq.WithPublisherOptionsExchangeDeclare,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	return &EventPublisher{
		publisher: publisher,
		exchange:  exchangeName,
	}, nil
}

func (p *EventPublisher) Publish(event interface{}, routingKey string) error {
	eventBytes, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("could not marshal event %+v: %w", event, err)
	}

	err = p.publisher.Publish(
		eventBytes,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(p.exchange),
	)
	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("Published event with key %s to exchange %s", routingKey, p.exchange)
	return nil
}

func (p *EventPublisher) Close() {
	p.publisher.Close()
}
