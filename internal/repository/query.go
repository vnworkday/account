package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/pkg/errors"

	"github.com/vnworkday/account/internal/model"
)

type QueryBuilder[T any] struct {
	query            string
	selectClause     string
	fromClause       string
	paginationClause string
	whereClause      strings.Builder
	whereArgs        []any
	sortClause       strings.Builder
	err              error
}

func NewQueryBuilder[T any]() *QueryBuilder[T] {
	return &QueryBuilder[T]{}
}

func (b *QueryBuilder[T]) Select(fields ...string) *QueryBuilder[T] {
	if b.err != nil {
		return b
	}

	if len(fields) == 0 {
		b.err = errors.New("repository: fields in select are required")

		return b
	}

	b.selectClause = "SELECT " + strings.Join(fields, ", ")

	return b
}

func (b *QueryBuilder[T]) From(tables ...string) *QueryBuilder[T] {
	if b.err != nil {
		return b
	}

	if len(tables) == 0 {
		b.err = errors.New("repository: tables in from are required")

		return b
	}

	b.fromClause += " FROM " + strings.Join(tables, ", ")

	return b
}

func (b *QueryBuilder[T]) WhereRaw(where string) *QueryBuilder[T] {
	if b.err != nil {
		return b
	}

	if b.whereClause.Len() > 0 {
		b.whereClause.WriteString(" AND ")
	} else {
		b.whereClause.WriteString(" WHERE ")
	}

	b.whereClause.WriteString(where)

	return b
}

func (b *QueryBuilder[T]) Where(filter model.Filter, optAlias ...string) *QueryBuilder[T] {
	if b.err != nil {
		return b
	}

	whereClause, err := StringifyFilter(filter, optAlias...)
	if err != nil {
		b.err = err

		return b
	}

	if b.whereClause.Len() > 0 {
		b.whereClause.WriteString(" AND ")
	} else {
		b.whereClause.WriteString(" WHERE ")
	}

	b.whereClause.WriteString(whereClause)
	b.whereArgs = append(b.whereArgs, filter.Value)

	return b
}

func (b *QueryBuilder[T]) OrderBy(sort model.Sort, optAlias ...string) *QueryBuilder[T] {
	if b.err != nil {
		return b
	}

	sortClause, err := StringifySort(sort, optAlias...)
	if err != nil {
		b.err = err

		return b
	}

	if b.sortClause.Len() > 0 {
		b.sortClause.WriteString(", ")
	} else {
		b.sortClause.WriteString(" ORDER BY ")
	}

	b.sortClause.WriteString(sortClause)

	return b
}

func (b *QueryBuilder[T]) Paginate(pagination model.Pagination) *QueryBuilder[T] {
	if b.err != nil {
		return b
	}

	if pagination.Offset < 0 {
		b.err = errors.New("repository: offset must be greater than or equal to 0")

		return b
	}

	if pagination.Limit < 0 {
		b.err = errors.New("repository: limit must be greater than 0")

		return b
	}

	var sb strings.Builder

	if pagination.Limit > 0 {
		sb.WriteString(fmt.Sprintf(" LIMIT %d", pagination.Limit))
	}

	if pagination.Offset > 0 {
		sb.WriteString(fmt.Sprintf(" OFFSET %d", pagination.Offset))
	}

	b.paginationClause = sb.String()

	return b
}

func (b *QueryBuilder[T]) build() (string, error) {
	if b.err != nil {
		return "", b.err
	}

	if b.selectClause == "" {
		return "", errors.New("repository: select clause is required")
	}

	if b.fromClause == "" {
		return "", errors.New("repository: from clause is required")
	}

	b.query = b.selectClause +
		b.fromClause +
		b.whereClause.String() +
		b.sortClause.String() +
		b.paginationClause

	return b.query, nil
}

func (b *QueryBuilder[T]) reset() {
	b.query = ""
	b.selectClause = ""
	b.fromClause = ""
	b.whereClause.Reset()
	b.whereArgs = nil
	b.sortClause.Reset()
	b.err = nil
}

func (b *QueryBuilder[T]) Query(ctx context.Context, db *sql.DB, scanner func(row *sql.Rows, out T) error) (*T, error) {
	var out T
	var rows *sql.Rows
	var query string
	var err error

	defer b.reset()

	if b.err != nil {
		return nil, b.err
	}

	query, err = b.build()
	if err != nil {
		return nil, err
	}

	rows, err = db.QueryContext(ctx, query, b.whereArgs...)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		if e := scanner(rows, out); e != nil {
			return nil, e
		}
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return &out, nil
}
