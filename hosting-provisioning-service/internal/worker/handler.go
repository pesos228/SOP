package worker

import (
	"context"
	"encoding/json"
	"fmt"
	"hosting-provisioning-service/internal/config"
	"hosting-provisioning-service/internal/messaging"
	"log"
	"math/rand"
	"time"

	events "hosting-events-contract"

	"github.com/wagslane/go-rabbitmq"
)

type CommandHandler struct {
	cfg       *config.Config
	publisher *messaging.EventPublisher
	ctx       context.Context
}

func NewCommandHandler(ctx context.Context, cfg *config.Config, publisher *messaging.EventPublisher) *CommandHandler {
	return &CommandHandler{
		cfg:       cfg,
		publisher: publisher,
		ctx:       ctx,
	}
}

func (h *CommandHandler) Handle(d rabbitmq.Delivery) rabbitmq.Action {
	var cmd events.ProvisionServerCommand
	if err := json.Unmarshal(d.Body, &cmd); err != nil {
		log.Printf("ERROR: failed to unmarshal command: %v. Message will be dropped.", err)
		return rabbitmq.NackDiscard
	}

	log.Printf("Received provisioning request for server %s (%s)", cmd.Hostname, cmd.ServerID)
	log.Printf("Starting provisioning simulation for %v...", h.cfg.ProvisioningTime)

	select {
	case <-time.After(h.cfg.ProvisioningTime):
	case <-h.ctx.Done():
		log.Printf("Shutdown signal received during provisioning for server %s. Requeueing message.", cmd.ServerID)
		return rabbitmq.NackRequeue
	}

	if rand.Intn(5) == 0 {
		return h.handleFailure(cmd)
	}

	return h.handleSuccess(cmd)
}

func (h *CommandHandler) handleSuccess(cmd events.ProvisionServerCommand) rabbitmq.Action {
	successEvent := events.ServerProvisionedEvent{
		ServerID:      cmd.ServerID,
		IPv4Address:   fmt.Sprintf("192.168.1.%d", 100+rand.Intn(100)),
		ProvisionedAt: time.Now().UTC(),
	}

	if err := h.publisher.Publish(successEvent, events.ProvisionSucceededKey); err != nil {
		log.Printf("ERROR: failed to publish success event for server %s: %v. Retrying...", cmd.ServerID, err)
		return rabbitmq.NackRequeue
	}

	log.Printf("Successfully provisioned server %s and published success event.", cmd.ServerID)
	return rabbitmq.Ack
}

func (h *CommandHandler) handleFailure(cmd events.ProvisionServerCommand) rabbitmq.Action {
	failedEvent := events.ServerProvisionFailedEvent{
		ServerID: cmd.ServerID,
		Reason:   "failed to allocate resources on host (simulated)",
		FailedAt: time.Now().UTC(),
	}

	if err := h.publisher.Publish(failedEvent, events.ProvisionFailedKey); err != nil {
		log.Printf("ERROR: failed to publish failure event for server %s: %v. Retrying...", cmd.ServerID, err)
		return rabbitmq.NackRequeue
	}

	log.Printf("Provisioning failed for server %s and published failure event.", cmd.ServerID)
	return rabbitmq.Ack
}
