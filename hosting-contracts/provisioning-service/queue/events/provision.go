package events

import (
	"time"

	"github.com/google/uuid"
)

const (
	ProvisionSucceededKey     = "server.provision.succeeded"
	ProvisionFailedKey        = "server.provision.failed"
	ProvisionResultKeyPattern = "server.provision.*"
)

type ServerProvisionedEvent struct {
	ServerID      uuid.UUID `json:"serverId"`
	IPv4Address   string    `json:"ipv4Address"`
	ProvisionedAt time.Time `json:"provisionedAt"`
}

type ServerProvisionFailedEvent struct {
	ServerID uuid.UUID `json:"serverId"`
	Reason   string    `json:"reason"`
	FailedAt time.Time `json:"failedAt"`
}
