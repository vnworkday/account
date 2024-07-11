package repository

import (
	"testing"

	"github.com/vnworkday/account/internal/fixture"

	"github.com/vnworkday/account/internal/model"
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

			if err := fixture.ExpectationsWereMet(tt.want, got, tt.wantErr, gotErr); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestStringifyOp(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		operator model.Op
		want     string
		wantErr  bool
	}{
		{
			name:     "EqualityOperator",
			operator: model.Eq,
			want:     "=",
			wantErr:  false,
		},
		{
			name:     "NotEqualOperator",
			operator: model.Ne,
			want:     "<>",
			wantErr:  false,
		},
		{
			name:     "GreaterThanOperator",
			operator: model.Gt,
			want:     ">",
			wantErr:  false,
		},
		{
			name:     "LessThanOperator",
			operator: model.Lt,
			want:     "<",
			wantErr:  false,
		},
		{
			name:     "GreaterThanOrEqualOperator",
			operator: model.Ge,
			want:     ">=",
			wantErr:  false,
		},
		{
			name:     "LessThanOrEqualOperator",
			operator: model.Le,
			want:     "<=",
			wantErr:  false,
		},
		{
			name:     "InOperator",
			operator: model.In,
			want:     "IN",
			wantErr:  false,
		},
		{
			name:     "NotInOperator",
			operator: model.NotIn,
			want:     "NOT IN",
			wantErr:  false,
		},
		{
			name:     "ContainsOperator",
			operator: model.Contains,
			want:     "LIKE",
			wantErr:  false,
		},
		{
			name:     "NotContainsOperator",
			operator: model.NotContains,
			want:     "NOT LIKE",
			wantErr:  false,
		},
		{
			name:     "StartsWithOperator",
			operator: model.StartsWith,
			want:     "LIKE",
			wantErr:  false,
		},
		{
			name:     "EndsWithOperator",
			operator: model.EndsWith,
			want:     "LIKE",
			wantErr:  false,
		},
		{
			name:     "NullOperator",
			operator: model.Null,
			want:     "IS NULL",
			wantErr:  false,
		},
		{
			name:     "NotNullOperator",
			operator: model.NotNull,
			want:     "IS NOT NULL",
			wantErr:  false,
		},
		{
			name:     "BetweenOperator",
			operator: model.Between,
			want:     "BETWEEN",
			wantErr:  false,
		},
		{
			name:     "UnsupportedOperator",
			operator: model.Op(999),
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := stringifyOp(tt.operator)

			if err := fixture.ExpectationsWereMet(tt.want, got, tt.wantErr, gotErr); err != nil {
				t.Error(err)
			}
		})
	}
}
