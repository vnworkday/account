package grpc

import (
	"context"

	"buf.build/gen/go/ntduycs/vnworkday/grpc/go/account/tenant/v1/tenantv1grpc"
	tenantv1 "buf.build/gen/go/ntduycs/vnworkday/protocolbuffers/go/account/tenant/v1"
	"github.com/go-kit/kit/transport/grpc"
	"github.com/vnworkday/account/internal/common/adapter"
	"github.com/vnworkday/account/internal/usecase/tenant"
	"go.uber.org/fx"
)

type TenantGRPCServer struct {
	listTenantHandler   grpc.Handler
	getTenantHandler    grpc.Handler
	createTenantHandler grpc.Handler
	updateTenantHandler grpc.Handler
}

type TenantGRPCServerParams struct {
	fx.In
	Port tenant.Port
}

func NewTenantGRPCServer(params TenantGRPCServerParams) tenantv1grpc.TenantServiceServer {
	return &TenantGRPCServer{
		listTenantHandler: adapter.NewGRPCServer(
			params.Port.DoListTenants,
			tenant.ToListRequest,
			tenant.ToListResponse,
		),
		getTenantHandler: adapter.NewGRPCServer(
			params.Port.DoGetTenant,
			tenant.ToGetRequest,
			tenant.ToGetResponse,
		),
		createTenantHandler: adapter.NewGRPCServer(
			params.Port.DoCreateTenant,
			tenant.ToCreateRequest,
			tenant.ToCreateResponse,
		),
		updateTenantHandler: adapter.NewGRPCServer(
			params.Port.DoUpdateTenant,
			tenant.ToUpdateRequest,
			tenant.ToUpdateResponse,
		),
	}
}

func (s *TenantGRPCServer) CreateTenant(
	ctx context.Context,
	request *tenantv1.CreateTenantRequest,
) (*tenantv1.CreateTenantResponse, error) {
	return adapter.ServeGRPC[tenantv1.CreateTenantRequest, tenantv1.CreateTenantResponse](
		ctx, request, s.createTenantHandler,
	)
}

func (s *TenantGRPCServer) GetTenant(
	ctx context.Context,
	request *tenantv1.GetTenantRequest,
) (*tenantv1.GetTenantResponse, error) {
	return adapter.ServeGRPC[tenantv1.GetTenantRequest, tenantv1.GetTenantResponse](
		ctx, request, s.getTenantHandler,
	)
}

func (s *TenantGRPCServer) ListTenants(
	ctx context.Context,
	request *tenantv1.ListTenantsRequest,
) (*tenantv1.ListTenantsResponse, error) {
	return adapter.ServeGRPC[tenantv1.ListTenantsRequest, tenantv1.ListTenantsResponse](
		ctx, request, s.listTenantHandler,
	)
}

func (s *TenantGRPCServer) UpdateTenant(
	ctx context.Context,
	request *tenantv1.UpdateTenantRequest,
) (*tenantv1.UpdateTenantResponse, error) {
	return adapter.ServeGRPC[tenantv1.UpdateTenantRequest, tenantv1.UpdateTenantResponse](
		ctx, request, s.updateTenantHandler,
	)
}
