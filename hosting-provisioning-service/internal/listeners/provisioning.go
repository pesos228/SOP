package listeners

import (
	"context"
	"encoding/json"
	"fmt"
	events "hosting-events-contract"
	"hosting-kit/messaging"
	"hosting-provisioning-service/internal/service"
	"log"
	"math/rand"
	"time"
)

type ProvisioningListener struct {
	provisioningTime  time.Duration
	provessionService service.ProvisioningService
}

func NewProvisioningListener(provisioningTime time.Duration, provisioningService service.ProvisioningService) *ProvisioningListener {
	return &ProvisioningListener{
		provisioningTime:  provisioningTime,
		provessionService: provisioningService,
	}
}

func (pl *ProvisioningListener) Handle(ctx context.Context, body []byte) error {
	var cmd events.ProvisionServerCommand
	if err := json.Unmarshal(body, &cmd); err != nil {
		log.Printf("ERROR: failed to unmarshal command: %v. Message will be dropped.", err)
		return fmt.Errorf("%w: failed to unmarshal ServerProvisionedEvent: %v", messaging.ErrPermanentFailure, err)
	}

	log.Printf("Received provisioning request for server %s (%s)", cmd.Hostname, cmd.ServerID)
	log.Printf("Starting provisioning simulation for %v...", pl.provisioningTime)

	select {
	case <-time.After(pl.provisioningTime):
	case <-ctx.Done():
		log.Printf("Shutdown signal received during provisioning for server %s. Requeueing message.", cmd.ServerID)
		return fmt.Errorf("shutdown signal received")
	}

	if rand.Intn(5) == 0 {
		return pl.provessionService.HandleProvisionFailure(cmd.ServerID)
	}

	return pl.provessionService.HandleProvisionSuccess(cmd.ServerID)
}
