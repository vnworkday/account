package tenant

import (
	"context"

	"github.com/vnworkday/account/internal/domain/repository"

	"github.com/pkg/errors"
	"go.uber.org/fx"
)

type Validator interface {
	ValidateCreateTenant(ctx context.Context, request *CreateTenantRequest) error
	ValidateUpdateTenant(ctx context.Context, request *UpdateTenantRequest) error
}

type ValidatorParams struct {
	fx.In
	Store repository.TenantRepo `name:"tenant_store"`
}

type validator struct {
	store repository.TenantRepo
}

func NewValidator(params ValidatorParams) Validator {
	return &validator{
		store: params.Store,
	}
}

func (v validator) ValidateCreateTenant(ctx context.Context, request *CreateTenantRequest) error {
	var err error
	var exist bool

	exist, err = v.store.ExistByName(ctx, request.Name)
	if err != nil {
		return errors.Wrap(err, "validator: cannot validate create tenant request")
	}

	if exist {
		return errors.New("validator: tenant is already existing with given name")
	}

	exist, err = v.store.ExistByDomain(ctx, request.Domain)
	if err != nil {
		return errors.Wrap(err, "validator: cannot validate create tenant request")
	}

	if exist {
		return errors.New("validator: tenant is already existing with given port")
	}

	return nil
}

func (v validator) ValidateUpdateTenant(ctx context.Context, request *UpdateTenantRequest) error {
	var err error
	var exist bool

	exist, err = v.store.ExistByNameAndIDNot(ctx, request.Name, request.ID)
	if err != nil {
		return errors.Wrap(err, "validator: cannot validate update tenant request")
	}

	if exist {
		return errors.New("validator: tenant is already existing with given name")
	}

	return nil
}
