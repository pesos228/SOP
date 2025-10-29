package dto

import (
	"hosting-service/internal/domain"
	"time"

	"github.com/google/uuid"
)

type ServerPreview struct {
	ID          uuid.UUID `json:"id"`
	IPv4Address *string   `json:"IPv4Address,omitempty"`
	PlanID      uuid.UUID `json:"planId"`
	Name        string    `json:"name"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ServerSearch struct {
	Data []*ServerPreview `json:"data"`
	Meta PaginationResult `json:"meta"`
}

func NewServerPreview(server *domain.Server) *ServerPreview {
	return &ServerPreview{
		ID:          server.ID,
		IPv4Address: server.IPv4Address,
		PlanID:      server.PlanID,
		Name:        server.Name,
		Status:      string(server.Status),
		CreatedAt:   server.CreatedAt,
	}
}
