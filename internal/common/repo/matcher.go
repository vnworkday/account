package repo

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/gookit/goutil/arrutil"
	"github.com/gookit/goutil/reflects"
	"github.com/gookit/goutil/structs"
	"github.com/gookit/goutil/strutil"

	"github.com/pkg/errors"
)

type MatcherBuilder[T any] struct {
	err     error
	mb      *MutationBuilder[T]
	matched bool
}

type Setter struct {
	Field string
	Value any
}

func ToSetters[T any](in T) ([]Setter, error) {
	temp, err := structs.StructToMap(in)
	if err != nil {
		return nil, err
	}

	setters := make([]Setter, 0, len(temp))

	for key, value := range temp {
		setters = append(setters, Setter{Field: key, Value: value})
	}

	return setters, nil
}

func NewErrorMatcher[T any](matched bool, err error) *MatcherBuilder[T] {
	return &MatcherBuilder[T]{
		err:     err,
		matched: matched,
	}
}

func NewMatcher[T any](mb *MutationBuilder[T], matched bool) *MatcherBuilder[T] {
	return &MatcherBuilder[T]{
		mb:      mb,
		matched: matched,
	}
}

func (m *MatcherBuilder[T]) ThenDoNothing() *MutationBuilder[T] {
	if m.err != nil {
		return m.mb
	}

	if m.matched {
		m.mb.matchClause.WriteString(" THEN DO NOTHING")
	} else {
		m.mb.notMatchClause.WriteString(" THEN DO NOTHING")
	}

	return m.mb
}

func (m *MatcherBuilder[T]) ThenDelete() *MutationBuilder[T] {
	if m.err != nil {
		return m.mb
	}

	if m.matched {
		m.mb.matchClause.WriteString(" THEN DELETE")
	} else {
		m.mb.notMatchClause.WriteString(" THEN DELETE")
	}

	return m.mb
}

func (m *MatcherBuilder[T]) ThenUpdate(columns ...string) *MutationBuilder[T] {
	if m.err != nil {
		return m.mb
	}

	setters := arrutil.Map[string, Setter](
		columns,
		func(col string) (Setter, bool) {
			return Setter{
				Field: col,
				Value: sourceAlias + "." + col,
			}, true
		})

	return m.thenUpdate(setters...)
}

func (m *MatcherBuilder[T]) thenUpdate(setters ...Setter) *MutationBuilder[T] {
	if m.err != nil {
		return m.mb
	}

	if len(setters) == 0 {
		return m.ThenDoNothing()
	}

	if m.matched {
		m.mb.matchClause.WriteString(" THEN UPDATE SET ")
	} else {
		err := errors.New("repository: cannot update when not matched")
		m.err = err
		m.mb.err = err

		return m.mb
	}

	for idx, setter := range setters {
		if reflects.IsNil(reflect.ValueOf(setter.Value)) {
			continue
		}

		if setter.Field == "" {
			err := errors.New("repository: field in setter is required")
			m.err = err
			m.mb.err = err

			return m.mb
		}

		if idx > 0 {
			m.mb.matchClause.WriteString(", ")
		}

		m.mb.matchClause.WriteString(setter.Field)
		m.mb.matchClause.WriteString(" = ")
		m.mb.matchClause.WriteString(strutil.SafeString(setter.Value))
	}

	return m.mb
}

func (m *MatcherBuilder[T]) ThenInsert(columns ...string) *MutationBuilder[T] {
	if m.err != nil {
		return m.mb
	}

	setters := arrutil.Map[string, Setter](
		columns,
		func(col string) (Setter, bool) {
			return Setter{
				Field: col,
				Value: sourceAlias + "." + col,
			}, true
		})

	return m.thenInsert(setters...)
}

func (m *MatcherBuilder[T]) thenInsert(setters ...Setter) *MutationBuilder[T] {
	if m.err != nil {
		return m.mb
	}

	if len(setters) == 0 {
		return m.ThenDoNothing()
	}

	if m.matched {
		err := errors.New("repository: cannot insert when matched")
		m.err = err
		m.mb.err = err

		return m.mb
	}

	var columns strings.Builder
	var values strings.Builder

	for idx, setter := range setters {
		if reflects.IsNil(reflect.ValueOf(setter.Value)) {
			continue
		}

		if setter.Field == "" {
			err := errors.New("repository: field in setter is required")
			m.err = err
			m.mb.err = err

			return m.mb
		}

		if idx > 0 {
			columns.WriteString(", ")
			values.WriteString(", ")
		}

		columns.WriteString(setter.Field)
		values.WriteString(strutil.SafeString(setter.Value))
	}

	m.mb.notMatchClause.WriteString(fmt.Sprintf(" THEN INSERT (%s) VALUES (%s)", columns.String(), values.String()))

	return m.mb
}

func (m *MatcherBuilder[T]) Exec(ctx context.Context, db *sql.DB) (int64, error) {
	if m.err != nil {
		return -1, m.err
	}

	return m.mb.Exec(ctx, db)
}
