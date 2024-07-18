package transport

import (
	"context"

	"github.com/vnworkday/account/internal/domain/tenant"
	"github.com/vnworkday/account/internal/model"

	"buf.build/gen/go/ntduycs/vnworkday/grpc/go/account/tenant/v1/tenantv1grpc"
	tenantv1 "buf.build/gen/go/ntduycs/vnworkday/protocolbuffers/go/account/tenant/v1"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/vnworkday/account/pkg/converter"
	"github.com/vnworkday/account/pkg/endpoint"
	"go.uber.org/fx"
)

type tenantServer struct {
	listTenantHandler   grpc.Handler
	getTenantHandler    grpc.Handler
	createTenantHandler grpc.Handler
	updateTenantHandler grpc.Handler
}

type TenantServerParams struct {
	fx.In
	Endpoints endpoint.TenantEndpoints
}

func NewTenantGrpcServer(params TenantServerParams) tenantv1grpc.TenantServiceServer {
	return &tenantServer{
		listTenantHandler: grpc.NewServer(
			params.Endpoints.DoListTenants,
			func(ctx context.Context, in any) (any, error) {
				return converter.Convert[tenantv1.ListTenantsRequest, model.ListRequest](
					ctx, in, converter.ToListRequest)
			},
			func(ctx context.Context, out any) (any, error) {
				return converter.Convert[model.ListResponse[tenant.Tenant], tenantv1.ListTenantsResponse](
					ctx, out, converter.ToListResponse)
			},
		),
	}
}

func (s *tenantServer) CreateTenant(
	ctx context.Context,
	request *tenantv1.CreateTenantRequest,
) (*tenantv1.CreateTenantResponse, error) {
	return Serve[tenantv1.CreateTenantRequest, tenantv1.CreateTenantResponse](ctx, request, s.createTenantHandler)
}

func (s *tenantServer) GetTenant(
	ctx context.Context,
	request *tenantv1.GetTenantRequest,
) (*tenantv1.GetTenantResponse, error) {
	return Serve[tenantv1.GetTenantRequest, tenantv1.GetTenantResponse](ctx, request, s.getTenantHandler)
}

func (s *tenantServer) ListTenants(
	ctx context.Context,
	request *tenantv1.ListTenantsRequest,
) (*tenantv1.ListTenantsResponse, error) {
	return Serve[tenantv1.ListTenantsRequest, tenantv1.ListTenantsResponse](ctx, request, s.listTenantHandler)
}

func (s *tenantServer) UpdateTenant(
	ctx context.Context,
	request *tenantv1.UpdateTenantRequest,
) (*tenantv1.UpdateTenantResponse, error) {
	return Serve[tenantv1.UpdateTenantRequest, tenantv1.UpdateTenantResponse](ctx, request, s.updateTenantHandler)
}
