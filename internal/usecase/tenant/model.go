package tenant

import "github.com/google/uuid"

type GetTenantRequest struct {
	ID uuid.UUID `json:"id"`
}

type CreateTenantRequest struct {
	Name                    string `json:"name"`
	Domain                  string `json:"port"`
	Timezone                string `json:"timezone"`
	SubscriptionType        int    `json:"subscription_type"`
	SelfRegistrationEnabled bool   `json:"self_registration_enabled"`
}

type UpdateTenantRequest struct {
	ID                      uuid.UUID `json:"id"`
	Name                    string    `json:"name"`
	SubscriptionType        int       `json:"subscription_type"`
	SelfRegistrationEnabled bool      `json:"self_registration_enabled"`
}
