package endpoint

import (
	"context"

	"github.com/pkg/errors"

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
}

type TenantEndpointsParams struct {
	fx.In
	Logger  *zap.Logger
	Service tenant.Service `name:"tenant_service"`
}

func NewTenantEndpoints(params TenantEndpointsParams) TenantEndpoints {
	doListTenants := MakeListEndpoint(params.Service,
		LoggingMiddleware(params.Logger.With(zap.String("method", "ListTenants"))),
	)
	doGetTenant := MakeGetEndpoint(params.Service,
		LoggingMiddleware(params.Logger.With(zap.String("method", "GetTenant"))),
	)
	doCreateTenant := MakeCreateEndpoint(params.Service,
		LoggingMiddleware(params.Logger.With(zap.String("method", "CreateTenant"))),
	)
	doUpdateTenant := MakeUpdateEndpoint(params.Service,
		LoggingMiddleware(params.Logger.With(zap.String("method", "UpdateTenant"))),
	)

	return TenantEndpoints{
		DoListTenants:  doListTenants,
		DoGetTenant:    doGetTenant,
		DoCreateTenant: doCreateTenant,
		DoUpdateTenant: doUpdateTenant,
	}
}

func MakeListEndpoint(svc tenant.Service, middlewares ...endpoint.Middleware) endpoint.Endpoint {
	ep := func(ctx context.Context, request any) (any, error) {
		req, ok := request.(*model.ListRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.ListTenants(ctx, req)
	}

	return applyMiddleware(ep, middlewares...)
}

func MakeGetEndpoint(svc tenant.Service, middlewares ...endpoint.Middleware) endpoint.Endpoint {
	ep := func(ctx context.Context, request any) (any, error) {
		req, ok := request.(*tenant.GetTenantRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.GetTenant(ctx, req)
	}

	return applyMiddleware(ep, middlewares...)
}

func MakeCreateEndpoint(svc tenant.Service, middlewares ...endpoint.Middleware) endpoint.Endpoint {
	ep := func(ctx context.Context, request any) (any, error) {
		req, ok := request.(*tenant.CreateTenantRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.CreateTenant(ctx, req)
	}

	return applyMiddleware(ep, middlewares...)
}

func MakeUpdateEndpoint(svc tenant.Service, middlewares ...endpoint.Middleware) endpoint.Endpoint {
	ep := func(ctx context.Context, request any) (any, error) {
		req, ok := request.(*tenant.UpdateTenantRequest)
		if !ok {
			return nil, errors.New("invalid request")
		}

		return svc.UpdateTenant(ctx, req)
	}

	return applyMiddleware(ep, middlewares...)
}

func (t TenantEndpoints) ListTenants(
	ctx context.Context,
	request *model.ListRequest,
) (*model.ListResponse[tenant.Tenant], error) {
	return Do[model.ListRequest, model.ListResponse[tenant.Tenant]](ctx, request, t.DoListTenants)
}

func (t TenantEndpoints) GetTenant(
	ctx context.Context,
	request *tenant.GetTenantRequest,
) (*tenant.Tenant, error) {
	return Do[tenant.GetTenantRequest, tenant.Tenant](ctx, request, t.DoGetTenant)
}

func (t TenantEndpoints) CreateTenant(
	ctx context.Context,
	request *tenant.CreateTenantRequest,
) (*tenant.Tenant, error) {
	return Do[tenant.CreateTenantRequest, tenant.Tenant](ctx, request, t.DoCreateTenant)
}

func (t TenantEndpoints) UpdateTenant(
	ctx context.Context,
	request *tenant.UpdateTenantRequest,
) (*tenant.Tenant, error) {
	return Do[tenant.UpdateTenantRequest, tenant.Tenant](ctx, request, t.DoUpdateTenant)
}
