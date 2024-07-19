package repo

import (
	"testing"

	"github.com/vnworkday/account/internal/common/domain"
	"github.com/vnworkday/account/internal/common/fixture"
)

func TestStringifyField(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		field     string
		sensitive bool
		alias     string
		want      string
		wantErr   bool
	}{
		{
			name:      "SimpleFieldWithoutAlias",
			field:     "username",
			sensitive: false,
			alias:     "",
			want:      "username",
			wantErr:   false,
		},
		{
			name:      "SimpleFieldWithAlias",
			field:     "username",
			sensitive: false,
			alias:     "users",
			want:      "users.username",
			wantErr:   false,
		},
		{
			name:      "SensitiveFieldWithoutAlias",
			field:     "username",
			sensitive: true,
			alias:     "",
			want:      "LOWER(username)",
			wantErr:   false,
		},
		{
			name:      "SensitiveFieldWithAlias",
			field:     "username",
			sensitive: true,
			alias:     "users",
			want:      "LOWER(users.username)",
			wantErr:   false,
		},
		{
			name:      "EmptyFieldName",
			field:     "",
			sensitive: false,
			alias:     "",
			want:      "",
			wantErr:   true,
		},
		{
			name:      "EmptyFieldNameWithAlias",
			field:     "",
			sensitive: false,
			alias:     "users",
			want:      "",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := stringifyField(tt.field, tt.sensitive, tt.alias)

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, gotErr)
		})
	}
}

func TestStringifyOp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		operator domain.Op
		want     string
		wantErr  bool
	}{
		{
			name:     "EqualityOperator",
			operator: domain.Eq,
			want:     "=",
			wantErr:  false,
		},
		{
			name:     "NotEqualOperator",
			operator: domain.Ne,
			want:     "<>",
			wantErr:  false,
		},
		{
			name:     "GreaterThanOperator",
			operator: domain.Gt,
			want:     ">",
			wantErr:  false,
		},
		{
			name:     "LessThanOperator",
			operator: domain.Lt,
			want:     "<",
			wantErr:  false,
		},
		{
			name:     "GreaterThanOrEqualOperator",
			operator: domain.Ge,
			want:     ">=",
			wantErr:  false,
		},
		{
			name:     "LessThanOrEqualOperator",
			operator: domain.Le,
			want:     "<=",
			wantErr:  false,
		},
		{
			name:     "InOperator",
			operator: domain.In,
			want:     "IN",
			wantErr:  false,
		},
		{
			name:     "NotInOperator",
			operator: domain.NotIn,
			want:     "NOT IN",
			wantErr:  false,
		},
		{
			name:     "ContainsOperator",
			operator: domain.Contains,
			want:     "LIKE",
			wantErr:  false,
		},
		{
			name:     "NotContainsOperator",
			operator: domain.NotContains,
			want:     "NOT LIKE",
			wantErr:  false,
		},
		{
			name:     "StartsWithOperator",
			operator: domain.StartsWith,
			want:     "LIKE",
			wantErr:  false,
		},
		{
			name:     "EndsWithOperator",
			operator: domain.EndsWith,
			want:     "LIKE",
			wantErr:  false,
		},
		{
			name:     "NullOperator",
			operator: domain.Null,
			want:     "IS NULL",
			wantErr:  false,
		},
		{
			name:     "NotNullOperator",
			operator: domain.NotNull,
			want:     "IS NOT NULL",
			wantErr:  false,
		},
		{
			name:     "BetweenOperator",
			operator: domain.Between,
			want:     "BETWEEN",
			wantErr:  false,
		},
		{
			name:     "UnsupportedOperator",
			operator: domain.Op(999),
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := stringifyOp(tt.operator)

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, gotErr)
		})
	}
}
