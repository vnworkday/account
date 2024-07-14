package tenant

import (
	"context"
	"database/sql"

	"github.com/pkg/errors"

	"go.uber.org/fx"

	"github.com/google/uuid"

	"github.com/vnworkday/account/internal/model"
	"github.com/vnworkday/account/internal/repository"
)

type Store interface {
	FindByID(ctx context.Context, id uuid.UUID) (*Tenant, error)
	FindByPublicID(ctx context.Context, publicID string) (*Tenant, error)
	FindAll(ctx context.Context, request *ListTenantsRequest) ([]*Tenant, error)

	ExistByName(ctx context.Context, name string) (bool, error)
	ExistByDomain(ctx context.Context, domain string) (bool, error)
	ExistByNameAndIDNot(ctx context.Context, name string, id uuid.UUID) (bool, error)

	CountAll(ctx context.Context, request *ListTenantsRequest) (int64, error)

	Save(ctx context.Context, tenant *Tenant) error
}

type StoreParams struct {
	fx.In
	DB *sql.DB
}

func NewStore(params StoreParams) (Store, error) {
	table, err := model.StructToTable(Tenant{}, "tenant")
	if err != nil {
		return nil, err
	}

	return &store{
		db:    params.DB,
		table: table,
	}, nil
}

type store struct {
	db    *sql.DB
	table *model.Table
}

func (r store) ExistByNameAndIDNot(ctx context.Context, name string, id uuid.UUID) (bool, error) {
	return repository.NewQueryBuilder[Tenant]().
		SelectExists().
		From(r.table.Name).
		Where(model.Filter{
			Field: "name",
			Op:    model.Eq,
			Value: name,
		}).
		Where(model.Filter{
			Field: "id",
			Op:    model.Ne,
			Value: id,
		}).
		Exist(ctx, r.db)
}

func (r store) ExistByDomain(ctx context.Context, domain string) (bool, error) {
	return repository.NewQueryBuilder[Tenant]().
		SelectExists().
		From(r.table.Name).
		Where(model.Filter{
			Field: "domain",
			Op:    model.Eq,
			Value: domain,
		}).
		Exist(ctx, r.db)
}

func (r store) ExistByName(ctx context.Context, name string) (bool, error) {
	return repository.NewQueryBuilder[Tenant]().
		SelectExists().
		From(r.table.Name).
		Where(model.Filter{
			Field: "name",
			Op:    model.Eq,
			Value: name,
		}).
		Exist(ctx, r.db)
}

func (r store) CountAll(ctx context.Context, request *ListTenantsRequest) (int64, error) {
	countBuilder := repository.NewQueryBuilder[Tenant]().
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

func (r store) FindAll(ctx context.Context, request *ListTenantsRequest) ([]*Tenant, error) {
	queryBuilder := repository.NewQueryBuilder[Tenant]().
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

func (r store) FindByID(ctx context.Context, id uuid.UUID) (*Tenant, error) {
	return repository.NewQueryBuilder[Tenant]().
		Select(r.table.Columns...).
		From(r.table.Name).
		Where(model.Filter{
			Field: "id",
			Op:    model.Eq,
			Value: id,
		}).
		Query(ctx, r.db, r.scanTo)
}

func (r store) FindByPublicID(ctx context.Context, publicID string) (*Tenant, error) {
	return repository.NewQueryBuilder[Tenant]().
		Select(r.table.Columns...).
		From(r.table.Name).
		Where(model.Filter{
			Field: "public_id",
			Op:    model.Eq,
			Value: publicID,
		}).
		Query(ctx, r.db, r.scanTo)
}

func (r store) Save(ctx context.Context, tenant *Tenant) error {
	_, err := repository.NewMutationBuilder[Tenant]().
		MergeInto(r.table.Name).
		Using(tenant).
		On(repository.MergeCondition{
			SourceCol: "id",
			TargetCol: "id",
			Op:        model.Eq,
		}).
		WhenMatched().
		ThenUpdate(r.table.Updatable...).
		WhenNotMatched().
		ThenInsert(r.table.Insertable...).
		Exec(ctx, r.db)

	return err
}

func (r store) scanTo(rows *sql.Rows, tenant Tenant) error {
	if err := rows.Scan(
		&tenant.ID,
		&tenant.Name,
		&tenant.State,
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
