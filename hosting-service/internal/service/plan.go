package service

import (
	"context"
	"errors"
	"hosting-service/internal/domain"
	"hosting-service/internal/dto"
	"hosting-service/internal/repository"

	"github.com/google/uuid"
)

var (
	ErrPlanNotFound = errors.New("plan not found")
)

type CreatePlanParams struct {
	Name     string
	CPUCores int
	RAMMB    int
	DiskGB   int
}

type PlanService interface {
	Save(ctx context.Context, params CreatePlanParams) (*dto.PlanPreview, error)
	Search(ctx context.Context, page, pageSize int) (*dto.PlanSearch, error)
	FindByID(ctx context.Context, ID uuid.UUID) (*dto.PlanPreview, error)
}

type planServiceImpl struct {
	planRepository repository.PlanRepository
}

func (p *planServiceImpl) FindByID(ctx context.Context, ID uuid.UUID) (*dto.PlanPreview, error) {
	plan, err := p.planRepository.FindByID(ctx, ID)
	if err != nil {
		if errors.Is(err, repository.ErrPlanNotFound) {
			return nil, ErrPlanNotFound
		}
		return nil, err
	}

	return dto.NewPlanPreview(plan), nil
}

func (p *planServiceImpl) Save(ctx context.Context, params CreatePlanParams) (*dto.PlanPreview, error) {
	plan, err := domain.NewPlan(domain.NewPlanParams{
		Name:     params.Name,
		CpuCores: params.CPUCores,
		RamMb:    params.RAMMB,
		DiskGb:   params.DiskGB,
	})

	if err != nil {
		return nil, err
	}

	err = p.planRepository.Save(ctx, plan)

	if err != nil {
		return nil, err
	}

	return dto.NewPlanPreview(plan), nil
}

func (p *planServiceImpl) Search(ctx context.Context, page int, pageSize int) (*dto.PlanSearch, error) {
	plans, count, err := p.planRepository.FindAll(ctx, page, pageSize)

	if err != nil {
		return nil, err
	}

	data := make([]*dto.PlanPreview, len(plans))
	for i, plan := range plans {
		data[i] = dto.NewPlanPreview(plan)
	}

	return &dto.PlanSearch{
		Data: data,
		Meta: repository.CalculatePaginationResult(page, pageSize, count),
	}, nil
}

func NewPlanService(planRepository repository.PlanRepository) PlanService {
	return &planServiceImpl{planRepository: planRepository}
}
