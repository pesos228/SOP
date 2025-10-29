package dto

import (
	"hosting-service/internal/domain"

	"github.com/google/uuid"
)

type PlanPreview struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	CPUCores int       `json:"cpuCores"`
	RAMMB    int       `json:"ramMb"`
	DiskGB   int       `json:"diskGb"`
}

type PlanSearch struct {
	Data []*PlanPreview   `json:"data"`
	Meta PaginationResult `json:"meta"`
}

func NewPlanPreview(plan *domain.Plan) *PlanPreview {
	return &PlanPreview{
		ID:       plan.ID,
		Name:     plan.Name,
		CPUCores: plan.CPUCores,
		RAMMB:    plan.RAMMB,
		DiskGB:   plan.DiskGB,
	}
}
