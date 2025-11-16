package hostingeventscontract

import (
	"time"

	"github.com/google/uuid"
)

const (
	ProvisionRequestKey       = "server.provision.request"
	ProvisionSucceededKey     = "server.provision.succeeded"
	ProvisionFailedKey        = "server.provision.failed"
	ProvisionResultKeyPattern = "server.provision.*"
)

type ProvisionServerCommand struct {
	ServerID uuid.UUID `json:"serverId"`
	Hostname string    `json:"hostname"`
	CPUCores int       `json:"cpuCores"`
	RAMMB    int       `json:"ramMb"`
	DiskGB   int       `json:"diskGb"`
}

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
