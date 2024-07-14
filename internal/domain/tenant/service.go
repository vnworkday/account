package tenant

import (
	"context"
	"time"

	"github.com/gookit/goutil/syncs"

	"github.com/google/uuid"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service interface {
	ListTenants(ctx context.Context, request *ListTenantsRequest) ([]*Tenant, int, error)
	GetTenant(ctx context.Context, request *GetTenantRequest) (*Tenant, error)
	CreateTenant(ctx context.Context, request *CreateTenantRequest) (*Tenant, error)
	UpdateTenant(ctx context.Context, request *UpdateTenantRequest) (*Tenant, error)
}

type ServiceParams struct {
	fx.In
	Logger *zap.Logger
	Store  Store `name:"tenant_store"`
}

func NewService(params ServiceParams) Service {
	return &service{
		logger: params.Logger,
		store:  params.Store,
	}
}

type service struct {
	logger *zap.Logger
	store  Store
}

func (s service) ListTenants(ctx context.Context, request *ListTenantsRequest) ([]*Tenant, int, error) {
	eg, egCtx := syncs.NewCtxErrGroup(ctx)

	var tenants []*Tenant
	var count int64

	eg.Go(func() error {
		var err error
		tenants, err = s.store.FindAll(egCtx, request)

		return err
	})

	eg.Go(func() error {
		var err error
		count, err = s.store.CountAll(egCtx, request)

		return err
	})

	err := eg.Wait()
	if err != nil {
		return nil, 0, err
	}

	return tenants, int(count), nil
}

func (s service) GetTenant(ctx context.Context, request *GetTenantRequest) (*Tenant, error) {
	tenant, err := s.store.FindByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s service) CreateTenant(ctx context.Context, request *CreateTenantRequest) (*Tenant, error) {
	now := time.Now()
	tenantID := uuid.New()
	tenant := &Tenant{
		ID:                      tenantID,
		Name:                    request.Name,
		State:                   1,
		Domain:                  request.Domain,
		Timezone:                request.Timezone,
		ProductionType:          1,
		SubscriptionType:        1,
		SelfRegistrationEnabled: request.SelfRegistrationEnabled,
		CreatedAt:               now,
		UpdatedAt:               now,
	}

	err := s.store.Save(ctx, tenant)
	if err != nil {
		return tenant, err
	}

	return s.store.FindByID(ctx, tenantID)
}

func (s service) UpdateTenant(ctx context.Context, request *UpdateTenantRequest) (*Tenant, error) {
	now := time.Now()
	tenantID := request.ID

	tenant, err := s.store.FindByID(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	tenant.Name = request.Name
	tenant.SubscriptionType = request.SubscriptionType
	tenant.SelfRegistrationEnabled = request.SelfRegistrationEnabled
	tenant.UpdatedAt = now

	err = s.store.Save(ctx, tenant)
	if err != nil {
		return tenant, err
	}

	return tenant, nil
}
