package service

import (
	"context"
	"errors"
	"fmt"
	"hosting-service/internal/domain"
	"hosting-service/internal/dto"
	"hosting-service/internal/messaging"
	"hosting-service/internal/repository"
	"log"

	events "hosting-events-contract"

	"github.com/google/uuid"
)

type ActionType string

const (
	ActionStart  ActionType = "START"
	ActionStop   ActionType = "STOP"
	ActionReboot ActionType = "REBOOT"
	ActionDelete ActionType = "DELETE"
)

type CreateServerParams struct {
	Name   string
	PlanID uuid.UUID
}

type PerformActionParams struct {
	ServerID uuid.UUID
	Action   ActionType
}

var (
	ErrServerNotFound = errors.New("server not found")
	ErrInvalidAction  = errors.New("invalid action")
)

type ServerService interface {
	Save(ctx context.Context, params CreateServerParams) (*dto.ServerPreview, error)
	Search(ctx context.Context, page, pageSize int) (*dto.ServerSearch, error)
	FindByID(ctx context.Context, ID uuid.UUID) (*dto.ServerPreview, error)
	PerformAction(ctx context.Context, params PerformActionParams) (*dto.ServerPreview, error)
	HandleProvisionSuccess(ctx context.Context, event events.ServerProvisionedEvent) error
	HandleProvisionFailure(ctx context.Context, event events.ServerProvisionFailedEvent) error
}

type serverServiceImpl struct {
	serverRepository repository.ServerRepository
	planRepository   repository.PlanRepository
	publisher        *messaging.EventPublisher
}

func (s *serverServiceImpl) HandleProvisionFailure(ctx context.Context, event events.ServerProvisionFailedEvent) error {
	server, err := s.serverRepository.FindByID(ctx, event.ServerID)
	if err != nil {
		log.Printf("ERROR: received provision failure for non-existent server %s", event.ServerID)
		return err
	}

	if err := server.ProvisionFailed(); err != nil {
		log.Printf("WARN: could not apply provision failure to server %s (current status: %s): %v", server.ID, server.Status, err)
		return nil
	}

	log.Printf("Updating server %s status to PROVISION_FAILED. Reason: %s", server.ID, event.Reason)
	return s.serverRepository.Save(ctx, server)
}

func (s *serverServiceImpl) HandleProvisionSuccess(ctx context.Context, event events.ServerProvisionedEvent) error {
	server, err := s.serverRepository.FindByID(ctx, event.ServerID)
	if err != nil {
		log.Printf("ERROR: received provision success for non-existent server %s", event.ServerID)
		return err
	}

	if err := server.ProvisionSucceeded(event.IPv4Address); err != nil {
		log.Printf("WARN: could not apply provision success to server %s (current status: %s): %v", server.ID, server.Status, err)
		return nil
	}

	log.Printf("Updating server %s status to STOPPED with IP %s", server.ID, event.IPv4Address)
	return s.serverRepository.Save(ctx, server)
}

func (s *serverServiceImpl) FindByID(ctx context.Context, ID uuid.UUID) (*dto.ServerPreview, error) {
	server, err := s.serverRepository.FindByID(ctx, ID)
	if err != nil {
		if errors.Is(err, repository.ErrServerNotFound) {
			return nil, ErrServerNotFound
		}
		return nil, err
	}

	return dto.NewServerPreview(server), nil
}

func (s *serverServiceImpl) PerformAction(ctx context.Context, params PerformActionParams) (*dto.ServerPreview, error) {
	server, err := s.serverRepository.FindByID(ctx, params.ServerID)
	if err != nil {
		if errors.Is(err, repository.ErrServerNotFound) {
			return nil, ErrServerNotFound
		}
		return nil, err
	}

	switch params.Action {
	case ActionStart:
		err = server.Start()
	case ActionStop:
		err = server.Stop()
	case ActionReboot:
		err = server.Reboot()
	case ActionDelete:
		err = server.MarkForDeletion()
	default:
		return nil, ErrInvalidAction
	}

	if err != nil {
		return nil, err
	}

	err = s.serverRepository.Save(ctx, server)

	if err != nil {
		return nil, err
	}

	return dto.NewServerPreview(server), nil
}

func (s *serverServiceImpl) Save(ctx context.Context, params CreateServerParams) (*dto.ServerPreview, error) {
	plan, err := s.planRepository.FindByID(ctx, params.PlanID)
	if err != nil {
		return nil, fmt.Errorf("plan with id %s not found: %w", params.PlanID, err)
	}

	server, err := domain.NewServer(params.PlanID, params.Name)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			return nil, err
		}
		return nil, err
	}

	err = s.serverRepository.Save(ctx, server)
	if err != nil {
		return nil, err
	}

	command := events.ProvisionServerCommand{
		ServerID: server.ID,
		Hostname: server.Name,
		CPUCores: plan.CPUCores,
		RAMMB:    plan.RAMMB,
		DiskGB:   plan.DiskGB,
	}

	if err := s.publisher.Publish(command, events.ProvisionRequestKey); err != nil {
		log.Printf("CRITICAL: failed to publish provision command for server %s: %v", server.ID, err)
		return nil, fmt.Errorf("internal error: failed to queue server for provisioning")
	}

	return dto.NewServerPreview(server), nil
}

func (s *serverServiceImpl) Search(ctx context.Context, page int, pageSize int) (*dto.ServerSearch, error) {
	servers, count, err := s.serverRepository.FindAll(ctx, page, pageSize)

	if err != nil {
		return nil, err
	}

	data := make([]*dto.ServerPreview, len(servers))
	for i, server := range servers {
		data[i] = dto.NewServerPreview(server)
	}

	return &dto.ServerSearch{
		Data: data,
		Meta: repository.CalculatePaginationResult(page, pageSize, count),
	}, nil
}

func NewServerService(serverRepository repository.ServerRepository, planRepository repository.PlanRepository, publisher *messaging.EventPublisher) ServerService {
	return &serverServiceImpl{
		serverRepository: serverRepository,
		planRepository:   planRepository,
		publisher:        publisher}
}
