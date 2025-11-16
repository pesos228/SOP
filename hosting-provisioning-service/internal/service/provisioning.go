package service

import (
	"fmt"
	events "hosting-events-contract"
	"hosting-kit/messaging"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

type ProvisioningService interface {
	HandleProvisionSuccess(serverID uuid.UUID) error
	HandleProvisionFailure(serverID uuid.UUID) error
}

type provisioningServiceImpl struct {
	messageManager *messaging.MessageManager
}

func NewProvisioningService(messageManager *messaging.MessageManager) ProvisioningService {
	return &provisioningServiceImpl{
		messageManager: messageManager,
	}
}

func (s *provisioningServiceImpl) HandleProvisionSuccess(serverID uuid.UUID) error {
	successEvent := events.ServerProvisionedEvent{
		ServerID:      serverID,
		IPv4Address:   fmt.Sprintf("192.168.1.%d", 100+rand.Intn(100)),
		ProvisionedAt: time.Now().UTC(),
	}

	if err := s.messageManager.Publish(events.EventsExchange, events.ProvisionSucceededKey, successEvent); err != nil {
		return fmt.Errorf("failed to publish ServerProvisionedEvent: %w", err)
	}

	return nil
}

func (s *provisioningServiceImpl) HandleProvisionFailure(serverID uuid.UUID) error {
	failedEvent := events.ServerProvisionFailedEvent{
		ServerID: serverID,
		Reason:   "failed to allocate resources on host (simulated)",
		FailedAt: time.Now().UTC(),
	}

	if err := s.messageManager.Publish(events.EventsExchange, events.ProvisionFailedKey, failedEvent); err != nil {
		return fmt.Errorf("failed to publish ServerProvisionFailedEvent: %w", err)
	}

	return nil
}
