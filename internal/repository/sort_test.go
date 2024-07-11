package repository

import (
	"testing"

	"github.com/vnworkday/account/internal/fixture"
	"github.com/vnworkday/account/internal/model"
)

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
				Field:         "username",
				CaseSensitive: false,
				Order:         model.Asc,
			},
			want:    "username ASC",
			wantErr: false,
		},
		{
			name: "DescendingSortWithAlias",
			sort: model.Sort{
				Field:         "created_at",
				CaseSensitive: false,
				Order:         model.Desc,
			},
			optAlias: []string{"users"},
			want:     "users.created_at DESC",
			wantErr:  false,
		},
		{
			name: "CaseSensitiveSortWithoutAlias",
			sort: model.Sort{
				Field:         "email",
				CaseSensitive: true,
				Order:         model.Asc,
			},
			want:    "LOWER(email) ASC",
			wantErr: false,
		},
		{
			name: "CaseSensitiveSortWithAlias",
			sort: model.Sort{
				Field:         "name",
				CaseSensitive: true,
				Order:         model.Desc,
			},
			optAlias: []string{"people"},
			want:     "LOWER(people.name) DESC",
			wantErr:  false,
		},
		{
			name: "EmptyFieldName",
			sort: model.Sort{
				Field:         "",
				CaseSensitive: false,
				Order:         model.Asc,
			},
			wantErr: true,
		},
		{
			name: "InvalidOrder",
			sort: model.Sort{
				Field:         "username",
				CaseSensitive: false,
				Order:         model.SortOrder("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := StringifySort(tt.sort, tt.optAlias...)

			if err := fixture.ExpectationsWereMet(tt.want, got, tt.wantErr, gotErr); err != nil {
				t.Error(err)
			}
		})
	}
}
