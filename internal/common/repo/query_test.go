package repo

import (
	"github.com/vnworkday/account/internal/common/domain"
	"github.com/vnworkday/account/internal/common/fixture"
	"testing"
)

func TestQueryBuilder_Select(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		fields  []string
		want    string
		wantErr bool
	}{
		{
			name:   "SingleField",
			fields: []string{"id"},
			want:   "SELECT id",
		},
		{
			name:   "MultipleFields",
			fields: []string{"id", "name", "email"},
			want:   "SELECT id, name, email",
		},
		{
			name:    "NoFields",
			fields:  []string{},
			want:    "",
			wantErr: true,
		},
		{
			name:    "NilFields",
			fields:  nil,
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			qb := &QueryBuilder[any]{}
			got := qb.Select(tt.fields...).selectClause

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, qb.err)
		})
	}
}

func TestQueryBuilder_From(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		tables  []string
		want    string
		wantErr bool
	}{
		{
			name:   "SingleTable",
			tables: []string{"users"},
			want:   " FROM users",
		},
		{
			name:   "MultipleTables",
			tables: []string{"users", "orders"},
			want:   " FROM users, orders",
		},
		{
			name:    "NoTables",
			tables:  []string{},
			wantErr: true,
		},
		{
			name:    "NilTables",
			tables:  nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			qb := &QueryBuilder[any]{}
			qb.From(tt.tables...)

			got := qb.fromClause

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, qb.err)
		})
	}
}

func TestQueryBuilder_Where(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		filter  domain.Filter
		want    string
		wantErr bool
	}{
		{
			name: "ValidSingleCondition",
			filter: domain.Filter{
				Field: "name",
				Op:    domain.Eq,
				Value: "John Doe",
			},
			want:    " WHERE name = ?",
			wantErr: false,
		},
		{
			name: "ValidMultipleConditions",
			filter: domain.Filter{
				Field: "age",
				Op:    domain.Gt,
				Value: 30,
			},
			want:    " WHERE age > ?",
			wantErr: false,
		},
		{
			name:    "InvalidFilterNoField",
			filter:  domain.Filter{},
			wantErr: true,
		},
		{
			name: "ValidWithAlias",
			filter: domain.Filter{
				Field: "users.name",
				Op:    domain.Eq,
				Value: "Jane Doe",
			},
			want:    " WHERE users.name = ?",
			wantErr: false,
		},
		{
			name: "FieldWithAlias",
			filter: domain.Filter{
				Field: "u.name", Op: domain.Eq, Value: "John Doe",
			},

			want:    " WHERE u.name = ?",
			wantErr: false,
		},
		{
			name: "UsingLessThanOperator",
			filter: domain.Filter{
				Field: "salary", Op: domain.Lt, Value: 50000,
			},

			want:    " WHERE salary < ?",
			wantErr: false,
		},
		{
			name: "UsingInOperatorWithSlice",
			filter: domain.Filter{
				Field: "department", Op: domain.In, Value: []string{"HR", "Engineering", "Marketing"},
			},
			want:    " WHERE department IN (?)",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			qb := &QueryBuilder[any]{}
			qb.Where(tt.filter)

			got := qb.whereClause.String()

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, qb.err)
		})
	}
}

func TestQueryBuilder_OrderBy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		sort    domain.Sort
		want    string
		wantErr bool
	}{
		{
			name: "OrderBySingleField",
			sort: domain.Sort{Field: "name", Order: domain.Asc},
			want: " ORDER BY name ASC",
		},
		{
			name: "OrderByMultipleFields",
			sort: domain.Sort{Field: "name, age", Order: domain.Desc},
			want: " ORDER BY name, age DESC",
		},
		{
			name:    "InvalidSortOption",
			sort:    domain.Sort{Field: "", Order: domain.Asc},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			qb := &QueryBuilder[any]{}
			qb.OrderBy(tt.sort)

			got := qb.sortClause.String()

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, qb.err)
		})
	}
}

func TestQueryBuilder_Paginate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		pagination domain.Pagination
		want       string
		wantErr    bool
	}{
		{
			name:       "PaginateWithLimitAndOffset",
			pagination: domain.Pagination{Limit: 10, Offset: 20},
			want:       " LIMIT 10 OFFSET 20",
		},
		{
			name:       "PaginateWithOnlyLimit",
			pagination: domain.Pagination{Limit: 10},
			want:       " LIMIT 10",
		},
		{
			name:       "PaginateWithOnlyOffset",
			pagination: domain.Pagination{Offset: 20},
			want:       " OFFSET 20",
		},
		{
			name:       "InvalidPaginationNegativeLimit",
			pagination: domain.Pagination{Limit: -10, Offset: 20},
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			qb := &QueryBuilder[any]{}
			qb.Paginate(tt.pagination)

			got := qb.paginationClause

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, qb.err)
		})
	}
}

func TestQueryBuilder_Build(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setupFunc func(qb *QueryBuilder[any])
		wantQuery string
		wantErr   bool
	}{
		{
			name: "CompleteQuery",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Select("id", "name").
					From("users").
					Where(domain.Filter{Field: "age", Op: domain.Gt, Value: 18}).
					OrderBy(domain.Sort{Field: "name", Order: domain.Asc}).
					Paginate(domain.Pagination{Limit: 10, Offset: 20})
			},
			wantQuery: "SELECT id, name FROM users WHERE age > ? ORDER BY name ASC LIMIT 10 OFFSET 20",
			wantErr:   false,
		},
		{
			name: "OnlySelectAndFrom",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Select("id", "name").From("users")
			},
			wantQuery: "SELECT id, name FROM users",
			wantErr:   false,
		},
		{
			name: "InvalidSelect",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Select()
			},
			wantQuery: "",
			wantErr:   true,
		},
		{
			name: "InvalidFrom",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.From()
			},
			wantQuery: "",
			wantErr:   true,
		},
		{
			name: "InvalidWhere",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Where(domain.Filter{})
			},
			wantQuery: "",
			wantErr:   true,
		},
		{
			name: "InvalidOrderBy",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.OrderBy(domain.Sort{})
			},
			wantQuery: "",
			wantErr:   true,
		},
		{
			name: "InvalidPaginateNegativeLimit",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Paginate(domain.Pagination{Limit: -1})
			},
			wantQuery: "",
			wantErr:   true,
		},
		{
			name: "InvalidPaginateNegativeOffset",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Paginate(domain.Pagination{Offset: -1})
			},
			wantQuery: "",
			wantErr:   true,
		},
		{
			name: "EmptyQuery",
			setupFunc: func(_ *QueryBuilder[any]) {
				// No setup actions, testing the builder with no inputs
			},
			wantQuery: "",
			wantErr:   true, // Expecting an error because no select or from clause is specified
		},
		{
			name: "SelectFromWhereWithInOperator",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Select("id", "name").
					From("employees").
					Where(domain.Filter{Field: "department", Op: domain.In, Value: []string{"HR", "Engineering"}})
			},
			wantQuery: "SELECT id, name FROM employees WHERE department IN (?)",
			wantErr:   false,
		},
		{
			name: "OrderByDescending",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Select("name").From("projects").OrderBy(domain.Sort{Field: "deadline", Order: domain.Desc})
			},
			wantQuery: "SELECT name FROM projects ORDER BY deadline DESC",
			wantErr:   false,
		},
		{
			name: "PaginateWithoutLimit",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Select("name").From("tasks").Paginate(domain.Pagination{Offset: 15})
			},
			wantQuery: "SELECT name FROM tasks OFFSET 15",
			wantErr:   false,
		},
		{
			name: "ComplexQuery",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Select("id", "name", "status").
					From("tickets").
					Where(domain.Filter{Field: "status", Op: domain.Eq, Value: "open"}).
					OrderBy(domain.Sort{Field: "priority", Order: domain.Asc}).
					Paginate(domain.Pagination{Limit: 5, Offset: 10})
			},
			wantQuery: "SELECT id, name, status FROM tickets WHERE status = ? ORDER BY priority ASC LIMIT 5 OFFSET 10",
			wantErr:   false,
		},
		{
			name: "MultipleChainedWhereAndOrderBy",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Select("id", "name", "email").
					From("employees").
					Where(domain.Filter{Field: "department", Op: domain.Eq, Value: "Engineering"}).
					Where(domain.Filter{Field: "location", Op: domain.Eq, Value: "New York"}).
					OrderBy(domain.Sort{Field: "name", Order: domain.Asc}).
					OrderBy(domain.Sort{Field: "id", Order: domain.Desc})
			},
			wantQuery: "SELECT id, name, email FROM employees WHERE department = ? AND location = ? ORDER BY name ASC, id DESC",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			qb := &QueryBuilder[any]{}
			tt.setupFunc(qb)
			gotQuery, gotErr := qb.build()

			fixture.ExpectationsWereMet(t, tt.wantQuery, gotQuery, tt.wantErr, gotErr)
		})
	}
}
