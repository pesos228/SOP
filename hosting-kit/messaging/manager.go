package messaging

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wagslane/go-rabbitmq"
)

type MessageHandler func(ctx context.Context, body []byte) error

type MessageManager struct {
	conn           *rabbitmq.Conn
	consumers      []*rabbitmq.Consumer
	publisher      *rabbitmq.Publisher
	wg             *sync.WaitGroup
	handlerTimeout time.Duration
}

func NewMessageManager(url string, exchanges []ExchangeConfig, wg *sync.WaitGroup, handlerTimeout time.Duration) (*MessageManager, error) {
	conn, err := rabbitmq.NewConn(url, rabbitmq.WithConnectionOptionsLogging)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	for _, ex := range exchanges {
		log.Printf("Declaring exchange '%s' of type '%s'", ex.Name, ex.Type)

		declarer, err := rabbitmq.NewPublisher(
			conn,
			rabbitmq.WithPublisherOptionsExchangeName(ex.Name),
			rabbitmq.WithPublisherOptionsExchangeKind(string(ex.Type)),
			rabbitmq.WithPublisherOptionsExchangeDeclare,
			rabbitmq.WithPublisherOptionsExchangeDurable,
		)
		if err != nil {
			conn.Close()
			return nil, fmt.Errorf("failed to declare exchange %s: %w", ex.Name, err)
		}
		declarer.Close()
	}

	publisher, err := rabbitmq.NewPublisher(
		conn,
	)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	return &MessageManager{
		conn:           conn,
		consumers:      []*rabbitmq.Consumer{},
		publisher:      publisher,
		wg:             wg,
		handlerTimeout: handlerTimeout,
	}, nil
}

func (m *MessageManager) Subscribe(queueName, routingKey, exchangeName string,
	exchangeType ExchangeType, handler MessageHandler,
) error {
	consumer, err := rabbitmq.NewConsumer(
		m.conn,
		queueName,
		rabbitmq.WithConsumerOptionsRoutingKey(routingKey),
		rabbitmq.WithConsumerOptionsExchangeName(exchangeName),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeKind(string(exchangeType)),
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsExchangeDurable,
	)
	if err != nil {
		return fmt.Errorf("failed to create consumer: %w", err)
	}

	m.consumers = append(m.consumers, consumer)

	log.Printf("Consumer registered for queue: %s", queueName)

	rabbitHandler := func(d rabbitmq.Delivery) rabbitmq.Action {
		ctx, cancel := context.WithTimeout(context.Background(), m.handlerTimeout)
		defer cancel()

		err := handler(ctx, d.Body)

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
		return fmt.Errorf("could not marshal event %+v: %w", data, err)
	}

	return m.publisher.Publish(
		eventBytes,
		[]string{routingKey},
		rabbitmq.WithPublishOptionsContentType("application/json"),
		rabbitmq.WithPublishOptionsExchange(exchangeName),
	)
}

func (m *MessageManager) Stop(ctx context.Context) {
	for _, consumer := range m.consumers {
		consumer.CloseWithContext(ctx)
	}

	m.wg.Wait()

	m.publisher.Close()
	m.conn.Close()
}
