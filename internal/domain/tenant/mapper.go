package tenant

import (
	tenantv1 "buf.build/gen/go/ntduycs/vnworkday/protocolbuffers/go/account/tenant/v1"
	"github.com/google/uuid"
	"go.uber.org/fx"
)

type Mapper interface {
	ToTenantForCreate(request *tenantv1.CreateTenantRequest) *CreateTenantRequest
	ToTenantForUpdate(request *tenantv1.UpdateTenantRequest) *UpdateTenantRequest
	ToTenantForGet(request *tenantv1.GetTenantRequest) *GetTenantRequest
}

type MapperParams struct {
	fx.In
}

type mapper struct{}

func NewMapper(_ MapperParams) Mapper {
	return &mapper{}
}

func (m mapper) ToTenantForCreate(request *tenantv1.CreateTenantRequest) *CreateTenantRequest {
	return &CreateTenantRequest{
		Name:                    request.GetName(),
		Domain:                  request.GetDomain(),
		Timezone:                request.GetTimezone(),
		SubscriptionType:        int(request.GetSubscriptionType()),
		SelfRegistrationEnabled: request.GetSelfRegistrationEnabled(),
	}
}

func (m mapper) ToTenantForUpdate(request *tenantv1.UpdateTenantRequest) *UpdateTenantRequest {
	return &UpdateTenantRequest{
		ID:                      uuid.MustParse(request.GetId()),
		Name:                    request.GetName(),
		SubscriptionType:        int(request.GetSubscriptionType()),
		SelfRegistrationEnabled: request.GetSelfRegistrationEnabled(),
	}
}

func (m mapper) ToTenantForGet(request *tenantv1.GetTenantRequest) *GetTenantRequest {
	return &GetTenantRequest{
		ID: uuid.MustParse(request.GetId()),
	}
}
