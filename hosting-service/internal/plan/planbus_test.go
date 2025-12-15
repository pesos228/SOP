package plan_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"hosting-service/internal/plan"
	"hosting-service/internal/platform/page"
)

type mockStorer struct {
	CreateFunc   func(ctx context.Context, p plan.Plan) error
	FindByIDFunc func(ctx context.Context, ID uuid.UUID) (plan.Plan, error)
	FindAllFunc  func(ctx context.Context, pg page.Page) ([]plan.Plan, int, error)
}

func (m *mockStorer) Create(ctx context.Context, p plan.Plan) error {
	if m.CreateFunc != nil {
		return m.CreateFunc(ctx, p)
	}
	return nil
}

func (m *mockStorer) FindByID(ctx context.Context, ID uuid.UUID) (plan.Plan, error) {
	if m.FindByIDFunc != nil {
		return m.FindByIDFunc(ctx, ID)
	}
	return plan.Plan{}, nil
}

func (m *mockStorer) FindAll(ctx context.Context, pg page.Page) ([]plan.Plan, int, error) {
	if m.FindAllFunc != nil {
		return m.FindAllFunc(ctx, pg)
	}
	return nil, 0, nil
}

func Test_Create(t *testing.T) {
	validParams := plan.CreatePlanParams{
		Name:     "Premium",
		CPUCores: 4,
		RAMMB:    8192,
		DiskGB:   100,
	}

	type testCase struct {
		name      string
		params    plan.CreatePlanParams
		mockSetup func() *mockStorer
		wantErr   error
	}

	table := []testCase{
		{
			name:   "success",
			params: validParams,
			mockSetup: func() *mockStorer {
				return &mockStorer{
					CreateFunc: func(ctx context.Context, p plan.Plan) error {
						if p.Name != "Premium" {
							return errors.New("data corrupted before save")
						}
						return nil
					},
				}
			},
			wantErr: nil,
		},
		{
			name: "validation_error_empty_name",
			params: plan.CreatePlanParams{
				Name:     "",
				CPUCores: 2,
				RAMMB:    1024,
				DiskGB:   10,
			},
			mockSetup: func() *mockStorer { return &mockStorer{} },
			wantErr:   plan.ErrValidation,
		},
		{
			name: "validation_error_zero_cpu",
			params: plan.CreatePlanParams{
				Name:     "Bad CPU",
				CPUCores: 0,
				RAMMB:    1024,
				DiskGB:   10,
			},
			mockSetup: func() *mockStorer { return &mockStorer{} },
			wantErr:   plan.ErrValidation,
		},
		{
			name:   "storage_error",
			params: validParams,
			mockSetup: func() *mockStorer {
				return &mockStorer{
					CreateFunc: func(ctx context.Context, p plan.Plan) error {
						return errors.New("db connection lost")
					},
				}
			},
			wantErr: errors.New("db connection lost"),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			bus := plan.NewBusiness(tt.mockSetup())

			gotPlan, err := bus.Create(context.Background(), tt.params)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) && err.Error() != tt.wantErr.Error() {
					t.Errorf("got error %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if gotPlan.ID == uuid.Nil {
				t.Error("plan ID was not generated")
			}
			if gotPlan.Name != tt.params.Name {
				t.Errorf("got name %s, want %s", gotPlan.Name, tt.params.Name)
			}
		})
	}
}

func Test_FindByID(t *testing.T) {
	targetID := uuid.New()
	foundPlan := plan.Plan{ID: targetID, Name: "Existing Plan"}

	type testCase struct {
		name      string
		id        uuid.UUID
		mockSetup func() *mockStorer
		wantPlan  plan.Plan
		wantErr   error
	}

	table := []testCase{
		{
			name: "success",
			id:   targetID,
			mockSetup: func() *mockStorer {
				return &mockStorer{
					FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (plan.Plan, error) {
						if ID == targetID {
							return foundPlan, nil
						}
						return plan.Plan{}, plan.ErrPlanNotFound
					},
				}
			},
			wantPlan: foundPlan,
			wantErr:  nil,
		},
		{
			name: "not_found",
			id:   uuid.New(),
			mockSetup: func() *mockStorer {
				return &mockStorer{
					FindByIDFunc: func(ctx context.Context, ID uuid.UUID) (plan.Plan, error) {
						return plan.Plan{}, plan.ErrPlanNotFound
					},
				}
			},
			wantPlan: plan.Plan{},
			wantErr:  plan.ErrPlanNotFound,
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			bus := plan.NewBusiness(tt.mockSetup())

			got, err := bus.FindByID(context.Background(), tt.id)

			if tt.wantErr != nil {
				if !errors.Is(err, tt.wantErr) {
					t.Errorf("got error %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if !reflect.DeepEqual(got, tt.wantPlan) {
				t.Errorf("got plan %+v, want %+v", got, tt.wantPlan)
			}
		})
	}
}

func Test_Search(t *testing.T) {
	pageParams := page.Parse(1, 10)

	plansList := []plan.Plan{
		{ID: uuid.New(), Name: "Plan A", CPUCores: 1},
		{ID: uuid.New(), Name: "Plan B", CPUCores: 2},
	}

	type testCase struct {
		name      string
		pg        page.Page
		mockSetup func() *mockStorer
		wantTotal int
		wantLen   int
		wantErr   error
	}

	table := []testCase{
		{
			name: "success_found_2_plans",
			pg:   pageParams,
			mockSetup: func() *mockStorer {
				return &mockStorer{
					FindAllFunc: func(ctx context.Context, pg page.Page) ([]plan.Plan, int, error) {
						if pg.Number() != 1 {
							return nil, 0, errors.New("wrong page number")
						}
						return plansList, 10, nil
					},
				}
			},
			wantTotal: 10,
			wantLen:   2,
			wantErr:   nil,
		},
		{
			name: "db_error",
			pg:   pageParams,
			mockSetup: func() *mockStorer {
				return &mockStorer{
					FindAllFunc: func(ctx context.Context, pg page.Page) ([]plan.Plan, int, error) {
						return nil, 0, errors.New("db connection failed")
					},
				}
			},
			wantTotal: 0,
			wantLen:   0,
			wantErr:   errors.New("db connection failed"),
		},
	}

	for _, tt := range table {
		t.Run(tt.name, func(t *testing.T) {
			bus := plan.NewBusiness(tt.mockSetup())

			gotPlans, gotTotal, err := bus.Search(context.Background(), tt.pg)

			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("got error %v, want %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if gotTotal != tt.wantTotal {
				t.Errorf("got total %d, want %d", gotTotal, tt.wantTotal)
			}
			if len(gotPlans) != tt.wantLen {
				t.Errorf("got plans length %d, want %d", len(gotPlans), tt.wantLen)
			}
		})
	}
}
