package repository

import (
	"context"
	"database/sql"

	"github.com/vnworkday/account/internal/common/domain"

	"github.com/vnworkday/account/internal/common/repo"
	"github.com/vnworkday/account/internal/domain/entity"

	"github.com/pkg/errors"

	"go.uber.org/fx"

	"github.com/google/uuid"
)

type TenantRepo interface {
	FindByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error)
	FindByPublicID(ctx context.Context, publicID string) (*entity.Tenant, error)
	FindAll(ctx context.Context, request *domain.ListRequest) ([]*entity.Tenant, error)

	ExistByName(ctx context.Context, name string) (bool, error)
	ExistByDomain(ctx context.Context, domain string) (bool, error)
	ExistByNameAndIDNot(ctx context.Context, name string, id uuid.UUID) (bool, error)

	CountAll(ctx context.Context, request *domain.ListRequest) (int64, error)

	Save(ctx context.Context, tenant *entity.Tenant) error
}

type TenantRepoParams struct {
	fx.In
	DB *sql.DB
}

func NewTenantRepo(params TenantRepoParams) (TenantRepo, error) {
	table, err := domain.StructToTable(entity.Tenant{}, "tenant")
	if err != nil {
		return nil, err
	}

	return &tenantRepo{
		db:    params.DB,
		table: table,
	}, nil
}

type tenantRepo struct {
	db    *sql.DB
	table *domain.Table
}

func (r tenantRepo) ExistByNameAndIDNot(ctx context.Context, name string, id uuid.UUID) (bool, error) {
	return repo.NewQueryBuilder[entity.Tenant]().
		SelectExists().
		From(r.table.Name).
		Where(domain.Filter{
			Field: "name",
			Op:    domain.Eq,
			Value: name,
		}).
		Where(domain.Filter{
			Field: "id",
			Op:    domain.Ne,
			Value: id,
		}).
		Exist(ctx, r.db)
}

func (r tenantRepo) ExistByDomain(ctx context.Context, domainStr string) (bool, error) {
	return repo.NewQueryBuilder[entity.Tenant]().
		SelectExists().
		From(r.table.Name).
		Where(domain.Filter{
			Field: "port",
			Op:    domain.Eq,
			Value: domainStr,
		}).
		Exist(ctx, r.db)
}

func (r tenantRepo) ExistByName(ctx context.Context, name string) (bool, error) {
	return repo.NewQueryBuilder[entity.Tenant]().
		SelectExists().
		From(r.table.Name).
		Where(domain.Filter{
			Field: "name",
			Op:    domain.Eq,
			Value: name,
		}).
		Exist(ctx, r.db)
}

func (r tenantRepo) CountAll(ctx context.Context, request *domain.ListRequest) (int64, error) {
	countBuilder := repo.NewQueryBuilder[entity.Tenant]().
		SelectCount().
		From(r.table.Name)

	for _, filter := range request.Filters {
		countBuilder = countBuilder.Where(filter)
	}

	count, err := countBuilder.Count(ctx, r.db)
	if err != nil {
		return 0, errors.Wrap(err, "repository: failed to count tenants")
	}

	return count, nil
}

func (r tenantRepo) FindAll(ctx context.Context, request *domain.ListRequest) ([]*entity.Tenant, error) {
	queryBuilder := repo.NewQueryBuilder[entity.Tenant]().
		Select(r.table.Columns...).
		From(r.table.Name)

	for _, filter := range request.Filters {
		queryBuilder = queryBuilder.Where(filter)
	}

	for _, sort := range request.Sorts {
		queryBuilder = queryBuilder.OrderBy(sort)
	}

	tenants, err := queryBuilder.QueryAll(ctx, r.db, r.scanTo)
	if err != nil {
		return nil, errors.Wrap(err, "repository: failed to find tenants")
	}

	return tenants, nil
}

func (r tenantRepo) FindByID(ctx context.Context, id uuid.UUID) (*entity.Tenant, error) {
	return repo.NewQueryBuilder[entity.Tenant]().
		Select(r.table.Columns...).
		From(r.table.Name).
		Where(domain.Filter{
			Field: "id",
			Op:    domain.Eq,
			Value: id,
		}).
		Query(ctx, r.db, r.scanTo)
}

func (r tenantRepo) FindByPublicID(ctx context.Context, publicID string) (*entity.Tenant, error) {
	return repo.NewQueryBuilder[entity.Tenant]().
		Select(r.table.Columns...).
		From(r.table.Name).
		Where(domain.Filter{
			Field: "public_id",
			Op:    domain.Eq,
			Value: publicID,
		}).
		Query(ctx, r.db, r.scanTo)
}

func (r tenantRepo) Save(ctx context.Context, tenant *entity.Tenant) error {
	_, err := repo.NewMutationBuilder[entity.Tenant]().
		MergeInto(r.table.Name).
		Using(tenant).
		On(repo.MergeCondition{
			SourceCol: "id",
			TargetCol: "id",
			Op:        domain.Eq,
		}).
		WhenMatched().
		ThenUpdate(r.table.Updatable...).
		WhenNotMatched().
		ThenInsert(r.table.Insertable...).
		Exec(ctx, r.db)

	return err
}

func (r tenantRepo) scanTo(rows *sql.Rows, tenant entity.Tenant) error {
	if err := rows.Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.Status,
		&tenant.Domain,
		&tenant.Timezone,
		&tenant.ProductionType,
		&tenant.SubscriptionType,
		&tenant.SelfRegistrationEnabled,
		&tenant.CreatedAt,
		&tenant.UpdatedAt,
	); err != nil {
		return err
	}

	return nil
}
