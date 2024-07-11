package tenant

import (
	"context"
	"database/sql"

	"github.com/gookit/goutil/arrutil"
	"github.com/vnworkday/account/internal/model"
	"github.com/vnworkday/account/internal/repository"
)

type Store interface {
	FindByPublicID(ctx context.Context, publicID string) (*Tenant, error)
	Save(ctx context.Context, tenant *Tenant) error
}

func NewDataStore(db *sql.DB) Store {
	allColumns := []string{
		"id",
		"public_id",
		"name",
		"state",
		"domain",
		"timezone",
		"production_type",
		"subscription_type",
		"self_registration_enabled",
		"created_at",
		"updated_at",
	}

	insertableColumns := arrutil.TakeWhile(allColumns, func(col string) bool {
		return col != "id" && col != "public_id"
	})

	updatableColumns := arrutil.TakeWhile(allColumns, func(col string) bool {
		return col != "id" && col != "public_id" && col != "created_at"
	})

	return &store{
		db:                db,
		table:             "tenant",
		allColumns:        allColumns,
		insertableColumns: insertableColumns,
		updatableColumns:  updatableColumns,
	}
}

type store struct {
	db                *sql.DB
	table             string
	allColumns        []string
	insertableColumns []string
	updatableColumns  []string
}

func (r store) Save(ctx context.Context, tenant *Tenant) error {
	_, err := repository.NewMutationBuilder[any]().
		MergeInto(r.table).
		Using(tenant).
		On(repository.MergeCondition{
			SourceCol: "id",
			TargetCol: "id",
			Op:        model.Eq,
		}).
		WhenMatched().
		ThenUpdate(r.updatableColumns...).
		WhenNotMatched().
		ThenInsert(r.insertableColumns...).
		Exec(ctx, r.db)
	if err != nil {
		return err
	}

	return nil
}

func (r store) FindByPublicID(ctx context.Context, publicID string) (*Tenant, error) {
	tenant, err := repository.NewQueryBuilder[Tenant]().
		Select(r.allColumns...).
		From(r.table).
		Where(model.Filter{
			Field: "public_id",
			Op:    model.Eq,
			Value: publicID,
		}).
		Query(ctx, r.db, r.scanTo)
	if err != nil {
		return nil, err
	}

	return &tenant, nil
}

func (r store) scanTo(rows *sql.Rows, tenant Tenant) error {
	if err := rows.Scan(
		&tenant.ID,
		&tenant.PublicID,
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
