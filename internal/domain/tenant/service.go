package tenant

import (
	"context"
	"time"

	"github.com/google/uuid"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service interface {
	ListTenants(ctx context.Context, request *ListTenantsRequest) (*ListTenantsResponse, error)
	GetTenant(ctx context.Context, request *GetTenantRequest) (*Tenant, error)
	CreateTenant(ctx context.Context, request *CreateTenantRequest) (*Tenant, error)
	UpdateTenant(ctx context.Context, request *UpdateTenantRequest) (*Tenant, error)
}

type ServiceParams struct {
	fx.In
	Logger *zap.Logger
	Store  Store `name:"tenant_store"`
}

func NewService(params ServiceParams) Service {
	return &service{
		logger: params.Logger,
		store:  params.Store,
	}
}

type service struct {
	logger *zap.Logger
	store  Store
}

func (s service) ListTenants(_ context.Context, _ *ListTenantsRequest) (*ListTenantsResponse, error) {
	panic("implement me")
}

func (s service) GetTenant(_ context.Context, _ *GetTenantRequest) (*Tenant, error) {
	panic("implement me")
}

func (s service) CreateTenant(ctx context.Context, request *CreateTenantRequest) (*Tenant, error) {
	now := time.Now()
	tenantID := uuid.New()
	tenant := &Tenant{
		ID:                      tenantID,
		Name:                    request.Name,
		State:                   1,
		Domain:                  request.Domain,
		Timezone:                request.Timezone,
		ProductionType:          1,
		SubscriptionType:        1,
		SelfRegistrationEnabled: request.SelfRegistrationEnabled,
		CreatedAt:               now,
		UpdatedAt:               now,
	}

	err := s.store.Save(ctx, tenant)
	if err != nil {
		return tenant, err
	}

	return s.store.FindByID(ctx, tenantID)
}

func (s service) UpdateTenant(_ context.Context, _ *UpdateTenantRequest) (*Tenant, error) {
	panic("implement me")
}
