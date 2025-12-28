package commands

import "github.com/google/uuid"

const (
	ProvisionRequestKey = "server.provision.request"
)

type ProvisionServerCommand struct {
	ServerID uuid.UUID `json:"serverId"`
	Hostname string    `json:"hostname"`
}
