package transport

import (
	"context"

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
		listTenantHandler: newGRPCServer(
			params.Endpoints.DoListTenants,
			converter.ToListRequest,
			converter.ToListResponse,
		),
		getTenantHandler: newGRPCServer(
			params.Endpoints.DoGetTenant,
			converter.ToGetRequest,
			converter.ToGetResponse,
		),
		createTenantHandler: newGRPCServer(
			params.Endpoints.DoCreateTenant,
			converter.ToCreateRequest,
			converter.ToCreateResponse,
		),
		updateTenantHandler: newGRPCServer(
			params.Endpoints.DoUpdateTenant,
			converter.ToUpdateRequest,
			converter.ToUpdateResponse,
		),
	}
}

func (s *tenantServer) CreateTenant(
	ctx context.Context,
	request *tenantv1.CreateTenantRequest,
) (*tenantv1.CreateTenantResponse, error) {
	return serveGRPC[tenantv1.CreateTenantRequest, tenantv1.CreateTenantResponse](ctx, request, s.createTenantHandler)
}

func (s *tenantServer) GetTenant(
	ctx context.Context,
	request *tenantv1.GetTenantRequest,
) (*tenantv1.GetTenantResponse, error) {
	return serveGRPC[tenantv1.GetTenantRequest, tenantv1.GetTenantResponse](ctx, request, s.getTenantHandler)
}

func (s *tenantServer) ListTenants(
	ctx context.Context,
	request *tenantv1.ListTenantsRequest,
) (*tenantv1.ListTenantsResponse, error) {
	return serveGRPC[tenantv1.ListTenantsRequest, tenantv1.ListTenantsResponse](ctx, request, s.listTenantHandler)
}

func (s *tenantServer) UpdateTenant(
	ctx context.Context,
	request *tenantv1.UpdateTenantRequest,
) (*tenantv1.UpdateTenantResponse, error) {
	return serveGRPC[tenantv1.UpdateTenantRequest, tenantv1.UpdateTenantResponse](ctx, request, s.updateTenantHandler)
}
