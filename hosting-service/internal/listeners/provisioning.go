package listeners

import (
	"context"
	"encoding/json"
	"fmt"

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
