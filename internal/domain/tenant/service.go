package tenant

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type IService interface {
	ListTenants(ctx context.Context, request ListTenantsRequest) (ListTenantsResponse, error)
	GetTenant(ctx context.Context, request GetTenantRequest) (Tenant, error)
	CreateTenant(ctx context.Context, request CreateTenantRequest) (Tenant, error)
	UpdateTenant(ctx context.Context, request UpdateTenantRequest) (Tenant, error)
}

type Params struct {
	fx.In

	Logger *zap.Logger
}

func NewService(params Params) *Service {
	return &Service{
		logger: params.Logger,
	}
}

type Service struct {
	logger *zap.Logger
}
