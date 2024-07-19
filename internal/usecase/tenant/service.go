package tenant

import (
	"context"
	"time"


	model2 "github.com/vnworkday/account/internal/common/domain"

	"github.com/vnworkday/account/internal/domain/entity"
	"github.com/vnworkday/account/internal/domain/repository"

	"github.com/gookit/goutil/syncs"

	"github.com/google/uuid"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Service interface {
	ListTenants(ctx context.Context, request *model2.ListRequest) (*model2.ListResponse[entity.Tenant], error)
	GetTenant(ctx context.Context, request *GetTenantRequest) (*entity.Tenant, error)
	CreateTenant(ctx context.Context, request *CreateTenantRequest) (*entity.Tenant, error)
	UpdateTenant(ctx context.Context, request *UpdateTenantRequest) (*entity.Tenant, error)
}

type ServiceParams struct {
	fx.In
	Logger *zap.Logger
	Store  repository.TenantRepo `name:"tenant_store"`
}

func NewService(params ServiceParams) Service {
	return &service{
		logger: params.Logger,
		store:  params.Store,
	}
}

type service struct {
	logger *zap.Logger
	store  repository.TenantRepo
}

func (s service) ListTenants(
	ctx context.Context,
	request *model2.ListRequest,
) (*model2.ListResponse[entity.Tenant], error) {
	eg, egCtx := syncs.NewCtxErrGroup(ctx)

	var tenants []*entity.Tenant
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
		return nil, err
	}

	return &model2.ListResponse[entity.Tenant]{
		Items: tenants,
		Count: int(count),
	}, nil
}

func (s service) GetTenant(
	ctx context.Context,
	request *GetTenantRequest,
) (*entity.Tenant, error) {
	tenant, err := s.store.FindByID(ctx, request.ID)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (s service) CreateTenant(
	ctx context.Context,
	request *CreateTenantRequest,
) (*entity.Tenant, error) {
	now := time.Now()
	tenantID := uuid.New()
	tenant := &entity.Tenant{
		ID:                      tenantID,
		Name:                    request.Name,
		Status:                  1,
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

func (s service) UpdateTenant(
	ctx context.Context,
	request *UpdateTenantRequest,
) (*entity.Tenant, error) {
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
