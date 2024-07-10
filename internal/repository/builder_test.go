package repository

import (
	"testing"

	"github.com/vnworkday/account/internal/model"
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

			if tt.wantErr && qb.err == nil {
				t.Errorf("Select() error = %v, wantErr %v", qb.err, tt.wantErr)
			} else if !tt.wantErr && qb.err != nil {
				t.Errorf("From() unexpected error: %v", qb.err)
			}

			if got != tt.want {
				t.Errorf("Select() got = %v, want %v", got, tt.want)
			}
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

			if tt.wantErr && qb.err == nil {
				t.Errorf("From() want error but got none")
			} else if !tt.wantErr && qb.err != nil {
				t.Errorf("From() unexpected error: %v", qb.err)
			}

			if got := qb.fromClause; got != tt.want && !tt.wantErr {
				t.Errorf("From() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryBuilder_Where(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		filter  model.Filter
		want    string
		wantErr bool
	}{
		{
			name: "ValidSingleCondition",
			filter: model.Filter{
				Field:    "name",
				Operator: model.Eq,
				Value:    "John Doe",
			},
			want:    " WHERE name = ?",
			wantErr: false,
		},
		{
			name: "ValidMultipleConditions",
			filter: model.Filter{
				Field:    "age",
				Operator: model.Gt,
				Value:    30,
			},
			want:    " WHERE age > ?",
			wantErr: false,
		},
		{
			name:    "InvalidFilterNoField",
			filter:  model.Filter{},
			wantErr: true,
		},
		{
			name: "ValidWithAlias",
			filter: model.Filter{
				Field:    "users.name",
				Operator: model.Eq,
				Value:    "Jane Doe",
			},
			want:    " WHERE users.name = ?",
			wantErr: false,
		},
		{
			name: "FieldWithAlias",
			filter: model.Filter{
				Field: "u.name", Operator: model.Eq, Value: "John Doe",
			},

			want:    " WHERE u.name = ?",
			wantErr: false,
		},
		{
			name: "UsingLessThanOperator",
			filter: model.Filter{
				Field: "salary", Operator: model.Lt, Value: 50000,
			},

			want:    " WHERE salary < ?",
			wantErr: false,
		},
		{
			name: "UsingInOperatorWithSlice",
			filter: model.Filter{
				Field: "department", Operator: model.In, Value: []string{"HR", "Engineering", "Marketing"},
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

			if tt.wantErr && qb.err == nil {
				t.Errorf("Where() want error but got none")
			} else if !tt.wantErr && qb.err != nil {
				t.Errorf("Where() unexpected error: %v", qb.err)
			}

			if got := qb.whereClause.String(); got != tt.want && !tt.wantErr {
				t.Errorf("Where() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryBuilder_OrderBy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		sort    model.Sort
		want    string
		wantErr bool
	}{
		{
			name: "OrderBySingleField",
			sort: model.Sort{Field: "name", Order: model.Asc},
			want: " ORDER BY name ASC",
		},
		{
			name: "OrderByMultipleFields",
			sort: model.Sort{Field: "name, age", Order: model.Desc},
			want: " ORDER BY name, age DESC",
		},
		{
			name:    "InvalidSortOption",
			sort:    model.Sort{Field: "", Order: model.Asc},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			qb := &QueryBuilder[any]{}
			qb.OrderBy(tt.sort)

			if tt.wantErr && qb.err == nil {
				t.Errorf("OrderBy() want error but got none")
			} else if !tt.wantErr && qb.err != nil {
				t.Errorf("OrderBy() unexpected error: %v", qb.err)
			}

			if got := qb.sortClause.String(); got != tt.want && !tt.wantErr {
				t.Errorf("OrderBy() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryBuilder_Paginate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		pagination model.Pagination
		want       string
		wantErr    bool
	}{
		{
			name:       "PaginateWithLimitAndOffset",
			pagination: model.Pagination{Limit: 10, Offset: 20},
			want:       " LIMIT 10 OFFSET 20",
		},
		{
			name:       "PaginateWithOnlyLimit",
			pagination: model.Pagination{Limit: 10},
			want:       " LIMIT 10",
		},
		{
			name:       "PaginateWithOnlyOffset",
			pagination: model.Pagination{Offset: 20},
			want:       " OFFSET 20",
		},
		{
			name:       "InvalidPaginationNegativeLimit",
			pagination: model.Pagination{Limit: -10, Offset: 20},
			want:       "",
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			qb := &QueryBuilder[any]{}
			qb.Paginate(tt.pagination)

			if tt.wantErr && qb.err == nil {
				t.Errorf("Paginate() want error but got none")
			} else if !tt.wantErr && qb.err != nil {
				t.Errorf("Paginate() unexpected error: %v", qb.err)
			}

			if got := qb.paginationClause; got != tt.want {
				t.Errorf("Paginate() got = %v, want %v", got, tt.want)
			}
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
					Where(model.Filter{Field: "age", Operator: model.Gt, Value: 18}).
					OrderBy(model.Sort{Field: "name", Order: model.Asc}).
					Paginate(model.Pagination{Limit: 10, Offset: 20})
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
				qb.Where(model.Filter{})
			},
			wantQuery: "",
			wantErr:   true,
		},
		{
			name: "InvalidOrderBy",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.OrderBy(model.Sort{})
			},
			wantQuery: "",
			wantErr:   true,
		},
		{
			name: "InvalidPaginateNegativeLimit",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Paginate(model.Pagination{Limit: -1})
			},
			wantQuery: "",
			wantErr:   true,
		},
		{
			name: "InvalidPaginateNegativeOffset",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Paginate(model.Pagination{Offset: -1})
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
					Where(model.Filter{Field: "department", Operator: model.In, Value: []string{"HR", "Engineering"}})
			},
			wantQuery: "SELECT id, name FROM employees WHERE department IN (?)",
			wantErr:   false,
		},
		{
			name: "OrderByDescending",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Select("name").From("projects").OrderBy(model.Sort{Field: "deadline", Order: model.Desc})
			},
			wantQuery: "SELECT name FROM projects ORDER BY deadline DESC",
			wantErr:   false,
		},
		{
			name: "PaginateWithoutLimit",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Select("name").From("tasks").Paginate(model.Pagination{Offset: 15})
			},
			wantQuery: "SELECT name FROM tasks OFFSET 15",
			wantErr:   false,
		},
		{
			name: "ComplexQuery",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Select("id", "name", "status").
					From("tickets").
					Where(model.Filter{Field: "status", Operator: model.Eq, Value: "open"}).
					OrderBy(model.Sort{Field: "priority", Order: model.Asc}).
					Paginate(model.Pagination{Limit: 5, Offset: 10})
			},
			wantQuery: "SELECT id, name, status FROM tickets WHERE status = ? ORDER BY priority ASC LIMIT 5 OFFSET 10",
			wantErr:   false,
		},
		{
			name: "MultipleChainedWhereAndOrderBy",
			setupFunc: func(qb *QueryBuilder[any]) {
				qb.Select("id", "name", "email").
					From("employees").
					Where(model.Filter{Field: "department", Operator: model.Eq, Value: "Engineering"}).
					Where(model.Filter{Field: "location", Operator: model.Eq, Value: "New York"}).
					OrderBy(model.Sort{Field: "name", Order: model.Asc}).
					OrderBy(model.Sort{Field: "id", Order: model.Desc})
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
			gotQuery, err := qb.build()

			if (err != nil) != tt.wantErr {
				t.Errorf("build() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && gotQuery != tt.wantQuery {
				t.Errorf("build() gotQuery = %v, want %v", gotQuery, tt.wantQuery)
			}
		})
	}
}
