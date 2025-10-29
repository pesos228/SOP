package messaging

import (
	"context"
	"fmt"
	"log"

	"github.com/wagslane/go-rabbitmq"
)

type ManagedConsumer struct {
	consumer *rabbitmq.Consumer
	handler  rabbitmq.Handler
}

func NewManagedConsumer(
	conn *rabbitmq.Conn,
	queueName,
	routingKey,
	exchangeName string,
	handler rabbitmq.Handler,
) (*ManagedConsumer, error) {
	log.Printf("Setting up consumer for queue '%s' with key '%s'", queueName, routingKey)

	consumer, err := rabbitmq.NewConsumer(
		conn,
		queueName,
		rabbitmq.WithConsumerOptionsRoutingKey(routingKey),
		rabbitmq.WithConsumerOptionsExchangeName(exchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsQueueDurable,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &ManagedConsumer{
		consumer: consumer,
		handler:  handler,
	}, nil
}

func (mc *ManagedConsumer) Run() error {
	log.Println("Consumer run loop started")
	err := mc.consumer.Run(mc.handler)
	log.Printf("Consumer run loop finished. Reason: %v", err)
	return err
}

func (mc *ManagedConsumer) CloseWithContext(ctx context.Context) {
	log.Println("Closing consumer with context...")
	mc.consumer.CloseWithContext(ctx)
}
