package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/wagslane/go-rabbitmq"
)

type MessageHandler func(ctx context.Context, body []byte, routingKey string) error

type MessageManager struct {
	conn           *rabbitmq.Conn
	consumers      []*rabbitmq.Consumer
	publisher      *rabbitmq.Publisher
	wg             sync.WaitGroup
	handlerTimeout time.Duration
}

func NewMessageManager(url string, exchanges []ExchangeConfig, handlerTimeout time.Duration) (*MessageManager, error) {
	conn, err := rabbitmq.NewConn(url, rabbitmq.WithConnectionOptionsLogging)
	if err != nil {
		return nil, err
	}

	for _, ex := range exchanges {
		declarer, err := rabbitmq.NewPublisher(
			conn,
			rabbitmq.WithPublisherOptionsExchangeName(ex.Name),
			rabbitmq.WithPublisherOptionsExchangeKind(string(ex.Type)),
			rabbitmq.WithPublisherOptionsExchangeDeclare,
			rabbitmq.WithPublisherOptionsExchangeDurable,
		)
		if err != nil {
			conn.Close()
			return nil, err
		}
		declarer.Close()
	}

	publisher, err := rabbitmq.NewPublisher(
		conn,
		rabbitmq.WithPublisherOptionsConfirm,
	)
	if err != nil {
		conn.Close()
		return nil, err
	}

	return &MessageManager{
		conn:           conn,
		consumers:      []*rabbitmq.Consumer{},
		publisher:      publisher,
		handlerTimeout: handlerTimeout,
	}, nil
}

func (m *MessageManager) Subscribe(queueName, routingKey, exchangeName string, handler MessageHandler, dlq *DLQConfig,
) error {

	opts := []func(*rabbitmq.ConsumerOptions){
		rabbitmq.WithConsumerOptionsRoutingKey(routingKey),
		rabbitmq.WithConsumerOptionsExchangeName(exchangeName),
		rabbitmq.WithConsumerOptionsQueueDurable,
	}

	if dlq != nil {
		opts = append(opts, rabbitmq.WithConsumerOptionsQueueArgs(map[string]interface{}{
			"x-dead-letter-exchange":    dlq.ExchangeName,
			"x-dead-letter-routing-key": dlq.RoutingKey,
		}))
	}

	consumer, err := rabbitmq.NewConsumer(
		m.conn,
		queueName,
		opts...,
	)
	if err != nil {
		return err
	}

	m.consumers = append(m.consumers, consumer)

	rabbitHandler := func(d rabbitmq.Delivery) rabbitmq.Action {
		ctx, cancel := context.WithTimeout(context.Background(), m.handlerTimeout)
		defer cancel()

		err := handler(ctx, d.Body, d.RoutingKey)

		if err == nil {
			return rabbitmq.Ack
		}

		if errors.Is(err, ErrPermanentFailure) {
			return rabbitmq.NackDiscard
		}

		return rabbitmq.NackRequeue
	}

	m.wg.Add(1)
	go func() {
		defer m.wg.Done()
		if err := consumer.Run(rabbitHandler); err != nil {
			log.Printf("Consumer stopped with error: %v", err)
		}
	}()

	return nil
}

func (m *MessageManager) Publish(exchangeName, routingKey string, data interface{}) error {
	eventBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	return m.publisher.Publish(
		eventBytes,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(exchangeName),
		rabbitmq.WithPublishOptionsPersistentDelivery,
	)
}

func (m *MessageManager) Stop(ctx context.Context) error {
	for _, consumer := range m.consumers {
		consumer.CloseWithContext(ctx)
	}

	m.wg.Wait()

	m.publisher.Close()
	if err := m.conn.Close(); err != nil {
		return err
	}
	return nil
}
