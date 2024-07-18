package endpoint

import (
	"context"

	"github.com/go-kit/kit/endpoint"
	"github.com/vnworkday/account/internal/domain/tenant"
	"github.com/vnworkday/account/internal/model"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type TenantEndpoints struct {
	DoListTenants  endpoint.Endpoint
	DoGetTenant    endpoint.Endpoint
	DoCreateTenant endpoint.Endpoint
	DoUpdateTenant endpoint.Endpoint
	service        tenant.Service
}

type TenantEndpointsParams struct {
	fx.In
	Logger  *zap.Logger
	Service tenant.Service `name:"tenant_service"`
}

func NewTenantEndpoints(params TenantEndpointsParams) TenantEndpoints {
	eps := TenantEndpoints{
		service: params.Service,
	}

	doListTenants := eps.makeListTenants(
		LoggingMiddleware(params.Logger.With(zap.String("method", "ListTenants"))),
	)
	doGetTenant := eps.makeGetTenant(
		LoggingMiddleware(params.Logger.With(zap.String("method", "GetTenant"))),
	)
	doCreateTenant := eps.makeCreateTenant(
		LoggingMiddleware(params.Logger.With(zap.String("method", "CreateTenant"))),
	)
	doUpdateTenant := eps.makeUpdateTenant(
		LoggingMiddleware(params.Logger.With(zap.String("method", "UpdateTenant"))),
	)

	eps.DoListTenants = doListTenants
	eps.DoGetTenant = doGetTenant
	eps.DoCreateTenant = doCreateTenant
	eps.DoUpdateTenant = doUpdateTenant

	return eps
}

func (t TenantEndpoints) makeListTenants(middlewares ...endpoint.Middleware) endpoint.Endpoint {
	return makeEndpoint[model.ListRequest, model.ListResponse[tenant.Tenant]](t.ListTenants, middlewares...)
}

func (t TenantEndpoints) makeGetTenant(middlewares ...endpoint.Middleware) endpoint.Endpoint {
	return makeEndpoint[tenant.GetTenantRequest, tenant.Tenant](t.GetTenant, middlewares...)
}

func (t TenantEndpoints) makeCreateTenant(middlewares ...endpoint.Middleware) endpoint.Endpoint {
	return makeEndpoint[tenant.CreateTenantRequest, tenant.Tenant](t.CreateTenant, middlewares...)
}

func (t TenantEndpoints) makeUpdateTenant(middlewares ...endpoint.Middleware) endpoint.Endpoint {
	return makeEndpoint[tenant.UpdateTenantRequest, tenant.Tenant](t.UpdateTenant, middlewares...)
}

func (t TenantEndpoints) ListTenants(
	ctx context.Context,
	request *model.ListRequest,
) (*model.ListResponse[tenant.Tenant], error) {
	return delegate[model.ListRequest, model.ListResponse[tenant.Tenant]](ctx, request, t.DoListTenants)
}

func (t TenantEndpoints) GetTenant(
	ctx context.Context,
	request *tenant.GetTenantRequest,
) (*tenant.Tenant, error) {
	return delegate[tenant.GetTenantRequest, tenant.Tenant](ctx, request, t.DoGetTenant)
}

func (t TenantEndpoints) CreateTenant(
	ctx context.Context,
	request *tenant.CreateTenantRequest,
) (*tenant.Tenant, error) {
	return delegate[tenant.CreateTenantRequest, tenant.Tenant](ctx, request, t.DoCreateTenant)
}

func (t TenantEndpoints) UpdateTenant(
	ctx context.Context,
	request *tenant.UpdateTenantRequest,
) (*tenant.Tenant, error) {
	return delegate[tenant.UpdateTenantRequest, tenant.Tenant](ctx, request, t.DoUpdateTenant)
}
