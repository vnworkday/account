package repository

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strings"

	"github.com/vnworkday/account/internal/model"

	"github.com/gookit/goutil/reflects"

	"github.com/pkg/errors"
)

const (
	sourceAlias = "source"
	targetAlias = "target"

	matchActionUpdate = "update"
	matchActionInsert = "insert"
)

type MergeCondition struct {
	SourceCol     string
	TargetCol     string
	Op            model.Op
	CaseSensitive bool
}

type MutationBuilder[T any] struct {
	query          string
	mergeClause    string
	usingClause    string
	usingArgValues []any
	usingArgKeys   []string
	onClause       strings.Builder
	notMatchClause strings.Builder
	matchClause    strings.Builder
	err            error
}

func NewMutationBuilder[T any]() *MutationBuilder[T] {
	return &MutationBuilder[T]{}
}

func (b *MutationBuilder[T]) MergeInto(table string) *MutationBuilder[T] {
	if b.err != nil {
		return b
	}

	if table == "" {
		b.err = errors.New("repository: target table is required")

		return b
	}

	targetTable := table

	b.mergeClause = fmt.Sprintf("MERGE INTO %s AS %s", targetTable, targetAlias)

	return b
}

func (b *MutationBuilder[T]) Using(values T) *MutationBuilder[T] {
	if b.err != nil {
		return b
	}

	setters, err := ToSetters(values)
	if err != nil {
		b.err = errors.Wrap(err, "repository: failed to convert values to setters")
	}

	return b.UsingValues(setters...)
}

func (b *MutationBuilder[T]) UsingValues(values ...Setter) *MutationBuilder[T] {
	if b.err != nil {
		return b
	}

	if len(values) == 0 {
		b.err = errors.New("repository: values are required")

		return b
	}

	wildcards := make([]string, 0, len(values))
	keys := make([]string, 0, len(values))
	vals := make([]any, 0, len(values))

	for _, setter := range values {
		if reflects.IsNil(reflect.ValueOf(setter.Value)) {
			continue
		}

		if setter.Field == "" {
			b.err = errors.New("repository: field in setter is required")

			return b
		}

		wildcards = append(wildcards, "?")
		keys = append(keys, setter.Field)
		vals = append(vals, setter.Value)
	}

	b.usingClause = fmt.Sprintf("USING (VALUES (%s)) AS %s (%s)",
		strings.Join(wildcards, ", "),
		sourceAlias,
		strings.Join(keys, ", "),
	)

	b.usingArgKeys = keys
	b.usingArgValues = vals

	return b
}

func (b *MutationBuilder[T]) On(mergeCond MergeCondition) *MutationBuilder[T] {
	if b.err != nil {
		return b
	}

	var src, op, target string
	var srcErr, opErr, targetErr error

	src, srcErr = stringifyField(mergeCond.SourceCol, mergeCond.CaseSensitive, sourceAlias)
	if srcErr != nil {
		b.err = errors.Wrap(srcErr, "repository: failed to stringify merge condition source")

		return b
	}

	target, targetErr = stringifyField(mergeCond.TargetCol, mergeCond.CaseSensitive, "target")
	if targetErr != nil {
		b.err = errors.Wrap(targetErr, "repository: failed to stringify merge condition target")

		return b
	}

	op, opErr = stringifyOp(mergeCond.Op)
	if opErr != nil {
		b.err = errors.Wrap(opErr, "repository: failed to stringify merge condition operator")

		return b
	}

	raw := fmt.Sprintf("%s %s %s", src, op, target)

	return b.onRaw(raw)
}

func (b *MutationBuilder[T]) onRaw(mergeCond string) *MutationBuilder[T] {
	if b.err != nil {
		return b
	}

	if mergeCond == "" {
		b.err = errors.New("repository: cannot concat empty merge condition")

		return b
	}

	if b.onClause.Len() > 0 {
		b.onClause.WriteString(" AND ")
	} else {
		b.onClause.WriteString("ON ")
	}

	b.onClause.WriteString(mergeCond)

	return b
}

func (b *MutationBuilder[T]) WhenNotMatched(cond ...string) *MatcherBuilder[T] {
	if b.err != nil {
		return NewErrorMatcher[T](true, b.err)
	}

	b.notMatchClause.WriteString("WHEN NOT MATCHED")

	if len(cond) > 0 {
		b.notMatchClause.WriteString(" AND ")
		b.notMatchClause.WriteString(strings.Join(cond, " AND "))
	}

	return NewMatcher[T](b, false)
}

func (b *MutationBuilder[T]) WhenMatched(cond ...string) *MatcherBuilder[T] {
	if b.err != nil {
		return NewErrorMatcher[T](true, b.err)
	}

	b.matchClause.WriteString("WHEN MATCHED")

	if len(cond) > 0 {
		b.matchClause.WriteString(" AND ")
		b.matchClause.WriteString(strings.Join(cond, " AND "))
	}

	return NewMatcher[T](b, true)
}

func (b *MutationBuilder[T]) build() (string, error) {
	if b.err != nil {
		return "", b.err
	}

	if b.mergeClause == "" {
		return "", errors.New("repository: merge clause is required")
	}

	if b.usingClause == "" {
		return "", errors.New("repository: using clause is required")
	}

	if b.onClause.Len() == 0 {
		return "", errors.New("repository: on clause is required")
	}

	if b.notMatchClause.Len() == 0 && b.matchClause.Len() == 0 {
		return "", errors.New("repository: at least one of not match or match clause is required")
	}

	query := fmt.Sprintf("%s %s %s",
		b.mergeClause,
		b.usingClause,
		b.onClause.String(),
	)

	if b.notMatchClause.Len() > 0 {
		query += " " + b.notMatchClause.String()
	}

	if b.matchClause.Len() > 0 {
		query += " " + b.matchClause.String()
	}

	return strings.TrimSpace(query), nil
}

func (b *MutationBuilder[T]) reset() {
	b.query = ""
	b.mergeClause = ""
	b.usingClause = ""
	b.usingArgValues = nil
	b.usingArgKeys = nil
	b.onClause.Reset()
	b.notMatchClause.Reset()
	b.matchClause.Reset()
	b.err = nil
}

func (b *MutationBuilder[T]) Exec(ctx context.Context, db *sql.DB) (int64, error) {
	var out sql.Result
	var query string
	var err error
	var rowsAffected int64

	defer b.reset()

	if b.err != nil {
		return -1, b.err
	}

	query, err = b.build()
	if err != nil {
		return -1, err
	}

	out, err = db.ExecContext(ctx, query, b.usingArgValues...)
	if err != nil {
		return -1, err
	}

	rowsAffected, err = out.RowsAffected()
	if err != nil {
		return -1, err
	}

	return rowsAffected, nil
}
