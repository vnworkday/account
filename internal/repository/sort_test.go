package repository

import (
	"strings"
	"testing"

	"github.com/vnworkday/account/internal/model"
)

func TestAppendSortClause(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		sorts []string
		want  string
	}{
		{
			name:  "NoSorts",
			sorts: []string{},
			want:  "",
		},
		{
			name:  "SingleSort",
			sorts: []string{"username ASC"},
			want:  " ORDER BY username ASC",
		},
		{
			name:  "MultipleSorts",
			sorts: []string{"username ASC", "created_at DESC"},
			want:  " ORDER BY username ASC, created_at DESC",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var query strings.Builder
			AppendSortClause(&query, tt.sorts...)
			got := query.String()

			if got != tt.want {
				t.Errorf("AppendSortClause() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringifySort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		sort     model.Sort
		optAlias []string
		want     string
		wantErr  bool
	}{
		{
			name: "AscendingSortWithoutAlias",
			sort: model.Sort{
				Field:           "username",
				IsCaseSensitive: false,
				Order:           model.Asc,
			},
			want:    "username ASC",
			wantErr: false,
		},
		{
			name: "DescendingSortWithAlias",
			sort: model.Sort{
				Field:           "created_at",
				IsCaseSensitive: false,
				Order:           model.Desc,
			},
			optAlias: []string{"users"},
			want:     "users.created_at DESC",
			wantErr:  false,
		},
		{
			name: "CaseSensitiveSortWithoutAlias",
			sort: model.Sort{
				Field:           "email",
				IsCaseSensitive: true,
				Order:           model.Asc,
			},
			want:    "LOWER(email) ASC",
			wantErr: false,
		},
		{
			name: "CaseSensitiveSortWithAlias",
			sort: model.Sort{
				Field:           "name",
				IsCaseSensitive: true,
				Order:           model.Desc,
			},
			optAlias: []string{"people"},
			want:     "LOWER(people.name) DESC",
			wantErr:  false,
		},
		{
			name: "EmptyFieldName",
			sort: model.Sort{
				Field:           "",
				IsCaseSensitive: false,
				Order:           model.Asc,
			},
			wantErr: true,
		},
		{
			name: "InvalidOrder",
			sort: model.Sort{
				Field:           "username",
				IsCaseSensitive: false,
				Order:           model.SortOrder("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := StringifySort(tt.sort, tt.optAlias...)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringifySort() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if got != tt.want {
				t.Errorf("StringifySort() got = %v, want %v", got, tt.want)
			}
		})
	}
}
