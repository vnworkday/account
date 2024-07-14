package repository

import (
	"context"
	"database/sql"

	"github.com/vnworkday/account/internal/model"
)

func ExistByFields[T any](ctx context.Context, db *sql.DB, tbName string, fieldValues []model.Filter) (bool, error) {
	builder := NewQueryBuilder[T]().
		SelectOne().
		From(tbName)

	for _, fieldValue := range fieldValues {
		builder = builder.Where(fieldValue)
	}

	return builder.Exist(ctx, db)
}
