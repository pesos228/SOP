package listeners

import (
	"context"
	"encoding/json"
	"fmt"

	"hosting-kit/messaging"
	"hosting-service/internal/service"

	events "hosting-events-contract"
)

type ProvisioningResultListener struct {
	serverService service.ServerService
}

func NewProvisioningResultListener(svc service.ServerService) *ProvisioningResultListener {
	listener := &ProvisioningResultListener{
		serverService: svc,
	}

	return listener
}

func (l *ProvisioningResultListener) Handle(ctx context.Context, body []byte, routingKey string) error {
	switch routingKey {
	case events.ProvisionSucceededKey:
		return l.HandleSuccess(ctx, body)
	case events.ProvisionFailedKey:
		return l.HandleFailure(ctx, body)
	default:
		return fmt.Errorf("%w: unknown routing key: %s", messaging.ErrPermanentFailure, routingKey)
	}
}

func (l *ProvisioningResultListener) HandleSuccess(ctx context.Context, body []byte) error {
	var event events.ServerProvisionedEvent

	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshal ServerProvisionedEvent failed: %w", err)
	}

	return l.serverService.HandleProvisionSuccess(ctx, event)
}

func (l *ProvisioningResultListener) HandleFailure(ctx context.Context, body []byte) error {
	var event events.ServerProvisionFailedEvent

	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("unmarshal ServerProvisionFailedEvent failed: %w", err)
	}

	return l.serverService.HandleProvisionFailure(ctx, event)
}
