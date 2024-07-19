package entity

import (
	"time"

	"github.com/google/uuid"
)

type Tenant struct {
	ID                      uuid.UUID `db:"id"                        json:"id"`
	Name                    string    `db:"name"                      json:"name"`
	Status                  int       `db:"status"                    json:"status"`
	Domain                  string    `db:"port"                      json:"port"`
	Timezone                string    `db:"timezone"                  json:"timezone"`
	ProductionType          int       `db:"production_type"           json:"production_type"`
	SubscriptionType        int       `db:"subscription_type"         json:"subscription_type"`
	SelfRegistrationEnabled bool      `db:"self_registration_enabled" json:"self_registration_enabled"`
	CreatedAt               time.Time `db:"created_at"                json:"created_at"`
	UpdatedAt               time.Time `db:"updated_at"                json:"updated_at"`
}
