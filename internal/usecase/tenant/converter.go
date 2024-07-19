package tenant

import (
	tenantv1 "buf.build/gen/go/ntduycs/vnworkday/protocolbuffers/go/account/tenant/v1"
	sharedv1 "buf.build/gen/go/ntduycs/vnworkday/protocolbuffers/go/shared/v1"
	"github.com/google/uuid"
	"github.com/gookit/goutil/arrutil"
	model2 "github.com/vnworkday/account/internal/common/domain"
	"github.com/vnworkday/account/internal/domain/entity"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToCreateRequest(request *tenantv1.CreateTenantRequest) *CreateTenantRequest {
	return &CreateTenantRequest{
		Name:                    request.GetName(),
		Domain:                  request.GetDomain(),
		Timezone:                request.GetTimezone(),
		SubscriptionType:        int(request.GetSubscriptionType()),
		SelfRegistrationEnabled: request.GetSelfRegistrationEnabled(),
	}
}

func ToCreateResponse(response *entity.Tenant) *tenantv1.CreateTenantResponse {
	return &tenantv1.CreateTenantResponse{
		Tenant: toGrpcTenant(response),
	}
}

func ToUpdateRequest(request *tenantv1.UpdateTenantRequest) *UpdateTenantRequest {
	return &UpdateTenantRequest{
		ID:                      uuid.MustParse(request.GetId()),
		Name:                    request.GetName(),
		SubscriptionType:        int(request.GetSubscriptionType()),
		SelfRegistrationEnabled: request.GetSelfRegistrationEnabled(),
	}
}

func ToUpdateResponse(response *entity.Tenant) *tenantv1.UpdateTenantResponse {
	return &tenantv1.UpdateTenantResponse{
		Tenant: toGrpcTenant(response),
	}
}

func ToGetRequest(request *tenantv1.GetTenantRequest) *GetTenantRequest {
	return &GetTenantRequest{
		ID: uuid.MustParse(request.GetId()),
	}
}

func ToGetResponse(response *entity.Tenant) *tenantv1.GetTenantResponse {
	return &tenantv1.GetTenantResponse{
		Tenant: toGrpcTenant(response),
	}
}

func ToListRequest(_ *tenantv1.ListTenantsRequest) *model2.ListRequest {
	return &model2.ListRequest{
		Pagination: model2.Pagination{
			Offset: 0,
			Limit:  0,
		},
		Filters: nil,
		Sorts:   nil,
	}
}

func ToListResponse(response *model2.ListResponse[entity.Tenant]) *tenantv1.ListTenantsResponse {
	return &tenantv1.ListTenantsResponse{
		Pagination: &sharedv1.ResponsePagination{
			NextToken:     "",
			PreviousToken: "",
			Total:         int32(response.Count),
			TotalPages:    0,
		},
		Tenants: arrutil.Map(response.Items, func(input *entity.Tenant) (*tenantv1.Tenant, bool) {
			return toGrpcTenant(input), true
		}),
	}
}

func toGrpcTenant(from *entity.Tenant) *tenantv1.Tenant {
	return &tenantv1.Tenant{
		Id:                      from.ID.String(),
		Name:                    from.Name,
		Status:                  tenantv1.TenantStatus(from.Status),
		Domain:                  from.Domain,
		Timezone:                from.Timezone,
		ProductionType:          tenantv1.TenantProductionType(from.ProductionType),
		SubscriptionType:        tenantv1.TenantSubscriptionType(from.SubscriptionType),
		SelfRegistrationEnabled: from.SelfRegistrationEnabled,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: from.CreatedAt.Unix(),
			Nanos:   int32(from.CreatedAt.Nanosecond()),
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: from.UpdatedAt.Unix(),
			Nanos:   int32(from.UpdatedAt.Nanosecond()),
		},
	}
}
