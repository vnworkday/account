package tenant

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service interface {
	ListTenants(ctx context.Context, request ListTenantsRequest) (ListTenantsResponse, error)
	GetTenant(ctx context.Context, request GetTenantRequest) (Tenant, error)
	CreateTenant(ctx context.Context, request CreateTenantRequest) (Tenant, error)
	UpdateTenant(ctx context.Context, request UpdateTenantRequest) (Tenant, error)
}

type Params struct {
	fx.In

	Logger *zap.Logger
}

func NewService(params Params) Service {
	return &service{
		logger: params.Logger,
	}
}

type service struct {
	logger *zap.Logger
}

func (s service) ListTenants(_ context.Context, _ ListTenantsRequest) (ListTenantsResponse, error) {
	panic("implement me")
}

func (s service) GetTenant(_ context.Context, _ GetTenantRequest) (Tenant, error) {
	panic("implement me")
}

func (s service) CreateTenant(_ context.Context, _ CreateTenantRequest) (Tenant, error) {
	panic("implement me")
}

func (s service) UpdateTenant(_ context.Context, _ UpdateTenantRequest) (Tenant, error) {
	panic("implement me")
}
