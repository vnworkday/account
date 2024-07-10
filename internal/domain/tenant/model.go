package tenant

import "github.com/vnworkday/account/internal/model"

type Tenant struct {
	ID                      string `json:"id"`
	PublicID                string `json:"public_id"`
	Name                    string `json:"name"`
	State                   int    `json:"state"`
	Domain                  string `json:"domain"`
	Timezone                string `json:"timezone"`
	ProductionType          int    `json:"production_type"`
	SubscriptionType        int    `json:"subscription_type"`
	SelfRegistrationEnabled bool   `json:"self_registration_enabled"`
	CreatedAt               string `json:"created_at"`
	UpdatedAt               string `json:"updated_at"`
}

type ListTenantsRequest struct {
	Pagination model.Pagination
	Filters    []model.Filter
	Sorts      []model.Sort
}

type ListTenantsResponse struct {
	model.Page
	Tenants []Tenant `json:"tenants"`
}

type GetTenantRequest struct {
	ID string `json:"id"`
}

type CreateTenantRequest struct {
	Name                    string `json:"name"`
	Domain                  string `json:"domain"`
	Timezone                string `json:"timezone"`
	SubscriptionType        int    `json:"subscription_type"`
	SelfRegistrationEnabled bool   `json:"self_registration_enabled"`
}

type UpdateTenantRequest struct {
	ID                      string `json:"id"`
	Name                    string `json:"name"`
	SubscriptionType        int    `json:"subscription_type"`
	SelfRegistrationEnabled bool   `json:"self_registration_enabled"`
}
