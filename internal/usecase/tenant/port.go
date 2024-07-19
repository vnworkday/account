package tenant

import (
	"context"

	"github.com/vnworkday/account/internal/common/domain"
	"github.com/vnworkday/account/internal/common/port"

	"github.com/vnworkday/account/internal/domain/entity"

	"github.com/go-kit/kit/endpoint"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Port struct {
	DoListTenants  endpoint.Endpoint
	DoGetTenant    endpoint.Endpoint
	DoCreateTenant endpoint.Endpoint
	DoUpdateTenant endpoint.Endpoint
}

type PortParams struct {
	fx.In
	Logger  *zap.Logger
	Service Service `name:"tenant_service"`
}

func NewPort(params PortParams) Port {
	return Port{
		DoListTenants: port.MakeEndpoint[domain.ListRequest, domain.ListResponse[entity.Tenant]](
			params.Service.ListTenants,
			port.LoggingMiddleware(params.Logger.With(zap.String("method", "ListTenants"))),
		),
		DoGetTenant: port.MakeEndpoint[GetTenantRequest, entity.Tenant](
			params.Service.GetTenant,
			port.LoggingMiddleware(params.Logger.With(zap.String("method", "GetTenant"))),
		),
		DoCreateTenant: port.MakeEndpoint[CreateTenantRequest, entity.Tenant](
			params.Service.CreateTenant,
			port.LoggingMiddleware(params.Logger.With(zap.String("method", "CreateTenant"))),
		),
		DoUpdateTenant: port.MakeEndpoint[UpdateTenantRequest, entity.Tenant](
			params.Service.UpdateTenant,
			port.LoggingMiddleware(params.Logger.With(zap.String("method", "UpdateTenant"))),
		),
	}
}

func (t Port) ListTenants(
	ctx context.Context,
	request *domain.ListRequest,
) (*domain.ListResponse[entity.Tenant], error) {
	return port.Delegate[domain.ListRequest, domain.ListResponse[entity.Tenant]](ctx, request, t.DoListTenants)
}

func (t Port) GetTenant(
	ctx context.Context,
	request *GetTenantRequest,
) (*entity.Tenant, error) {
	return port.Delegate[GetTenantRequest, entity.Tenant](ctx, request, t.DoGetTenant)
}

func (t Port) CreateTenant(
	ctx context.Context,
	request *CreateTenantRequest,
) (*entity.Tenant, error) {
	return port.Delegate[CreateTenantRequest, entity.Tenant](ctx, request, t.DoCreateTenant)
}

func (t Port) UpdateTenant(
	ctx context.Context,
	request *UpdateTenantRequest,
) (*entity.Tenant, error) {
	return port.Delegate[UpdateTenantRequest, entity.Tenant](ctx, request, t.DoUpdateTenant)
}
