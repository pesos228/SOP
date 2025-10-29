package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	"hosting-service/internal/service"
	"log"
	"time"

	events "hosting-events-contract"

	"github.com/wagslane/go-rabbitmq"
)

type eventHandlerFunc func(ctx context.Context, body []byte) error

type ProvisioningResultListener struct {
	serverService service.ServerService
	handlers      map[string]eventHandlerFunc
}

func NewProvisioningResultListener(svc service.ServerService) *ProvisioningResultListener {
	listener := &ProvisioningResultListener{
		serverService: svc,
	}

	listener.handlers = map[string]eventHandlerFunc{
		events.ProvisionSucceededKey: listener.handleSuccess,
		events.ProvisionFailedKey:    listener.handleFailure,
	}
	return listener
}

func (l *ProvisioningResultListener) Handle(d rabbitmq.Delivery) rabbitmq.Action {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	handler, ok := l.handlers[d.RoutingKey]
	if !ok {
		log.Printf("WARN: no handler registered for routing key: %s. Discarding message.", d.RoutingKey)
		return rabbitmq.NackDiscard
	}

	if err := handler(ctx, d.Body); err != nil {
		log.Printf("ERROR: handler for key '%s' failed: %v", d.RoutingKey, err)
		return rabbitmq.NackRequeue
	}

	log.Printf("Successfully handled event with key: %s", d.RoutingKey)
	return rabbitmq.Ack
}

func (l *ProvisioningResultListener) handleSuccess(ctx context.Context, body []byte) error {
	var event events.ServerProvisionedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshal ServerProvisionedEvent failed: %w", err)
	}
	return l.serverService.HandleProvisionSuccess(ctx, event)
}

func (l *ProvisioningResultListener) handleFailure(ctx context.Context, body []byte) error {
	var event events.ServerProvisionFailedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshal ServerProvisionFailedEvent failed: %w", err)
	}
	return l.serverService.HandleProvisionFailure(ctx, event)
}
