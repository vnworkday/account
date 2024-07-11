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
		Name:                    request.Name,
		Domain:                  request.Domain,
		Timezone:                request.Timezone,
		SubscriptionType:        int(request.SubscriptionType),
		SelfRegistrationEnabled: request.SelfRegistrationEnabled,
	}
}

func (m mapper) ToTenantForUpdate(request *tenantv1.UpdateTenantRequest) *UpdateTenantRequest {
	return &UpdateTenantRequest{
		ID:                      uuid.MustParse(request.Id),
		Name:                    request.Name,
		SubscriptionType:        int(request.SubscriptionType),
		SelfRegistrationEnabled: request.SelfRegistrationEnabled,
	}
}

func (m mapper) ToTenantForGet(request *tenantv1.GetTenantRequest) *GetTenantRequest {
	return &GetTenantRequest{
		ID: uuid.MustParse(request.Id),
	}
}
