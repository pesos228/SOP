package servergrpc

import (
	"context"
	"fmt"
	"hosting-service/internal/server"
	"hosting-service/internal/server/stores/servergrpc/gen"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ResourcesManager struct {
	client  gen.ResourcesClient
	timeOut time.Duration
}

func NewGrpc(client *grpc.ClientConn, timeOut time.Duration) *ResourcesManager {
	return &ResourcesManager{client: gen.NewResourcesClient(client), timeOut: timeOut}
}

func (r *ResourcesManager) Consume(ctx context.Context, resources server.Resources) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeOut)
	defer cancel()

	resp, err := r.client.ConsumeResource(ctx, &gen.ConsumeRequest{
		Resource: &gen.Resource{
			CpuCores: int32(resources.CPUCores),
			RamMb:    int32(resources.RAMMB),
			DiskGb:   int32(resources.DiskGB),
			IpCount:  int32(resources.IPCount),
		},
	})

	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.FailedPrecondition:
				return uuid.Nil, server.ErrNoResources
			case codes.InvalidArgument:
				return uuid.Nil, server.ErrValidation
			}
		}
		return uuid.Nil, fmt.Errorf("grpc: %w", err)
	}

	poolID, err := uuid.Parse(resp.GetPoolId())
	if err != nil {
		return uuid.Nil, fmt.Errorf("grpc: %w", err)
	}

	return poolID, nil
}

func (r *ResourcesManager) Return(ctx context.Context, resources server.Resources, poolID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeOut)
	defer cancel()

	_, err := r.client.ReturnResource(ctx, &gen.ReturnRequest{
		Resource: &gen.Resource{
			CpuCores: int32(resources.CPUCores),
			RamMb:    int32(resources.RAMMB),
			DiskGb:   int32(resources.DiskGB),
			IpCount:  int32(resources.IPCount),
		},
		PoolId: poolID.String(),
	})

	if err != nil {
		if st, ok := status.FromError(err); ok {
			switch st.Code() {
			case codes.NotFound:
				return server.ErrInvalidPlan
			case codes.InvalidArgument:
				return server.ErrValidation
			}
		}
		return fmt.Errorf("grpc: %w", err)
	}
	return nil
}
