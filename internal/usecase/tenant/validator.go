package tenant

import (
	"context"

	"github.com/vnworkday/account/internal/common/util"

	validator2 "github.com/vnworkday/account/internal/common/validator"
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
	Repo repository.TenantRepo `name:"tenant_store"`
}

type validator struct {
	repo repository.TenantRepo
}

func NewValidator(params ValidatorParams) Validator {
	return &validator{
		repo: params.Repo,
	}
}

func (v validator) ValidateCreateTenant(ctx context.Context, request *CreateTenantRequest) error {
	validations := []validator2.ValidationFunc{
		v.validateNameNotExists,
		v.validateDomainNotExists,
	}

	return validator2.Validate(ctx, request, validations...)
}

func (v validator) ValidateUpdateTenant(ctx context.Context, request *UpdateTenantRequest) error {
	validations := []validator2.ValidationFunc{
		v.validateNameNotExists,
	}

	return validator2.Validate(ctx, request, validations...)
}

// validateNameNotExists checks if the tenant name already exists.
func (v validator) validateNameNotExists(ctx context.Context, request any) error {
	var exist bool
	var err error

	switch util.Type(request) {
	case "*CreateTenantRequest":
		req := util.SafeCast[*CreateTenantRequest](request)
		exist, err = v.repo.ExistByName(ctx, req.Name)
	case "*UpdateTenantRequest":
		req := util.SafeCast[*UpdateTenantRequest](request)
		exist, err = v.repo.ExistByNameAndIDNot(ctx, req.Name, req.ID)
	default:
		return errors.New("validator: unrecognized request")
	}

	if err != nil {
		return errors.Wrap(err, "validator: cannot validate tenant name existence")
	}

	if exist {
		return errors.New("validator: tenant name already exists")
	}

	return nil
}

// validateDomainNotExists checks if the tenant domain already exists.
func (v validator) validateDomainNotExists(ctx context.Context, request any) error {
	req, ok := request.(*CreateTenantRequest)
	if !ok {
		return errors.New("validator: invalid request")
	}

	exist, err := v.repo.ExistByDomain(ctx, req.Domain)
	if err != nil {
		return errors.Wrap(err, "validator: cannot validate tenant domain existence")
	}

	if exist {
		return errors.New("validator: tenant domain already exists")
	}

	return nil
}
