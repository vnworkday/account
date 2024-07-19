package repo

import (
	"testing"


	"github.com/vnworkday/account/internal/common/domain"

	"github.com/vnworkday/account/internal/fixture"
)

func TestStringifySort(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		sort     domain.Sort
		optAlias []string
		want     string
		wantErr  bool
	}{
		{
			name: "AscendingSortWithoutAlias",
			sort: domain.Sort{
				Field:         "username",
				CaseSensitive: false,
				Order:         domain.Asc,
			},
			want:    "username ASC",
			wantErr: false,
		},
		{
			name: "DescendingSortWithAlias",
			sort: domain.Sort{
				Field:         "created_at",
				CaseSensitive: false,
				Order:         domain.Desc,
			},
			optAlias: []string{"users"},
			want:     "users.created_at DESC",
			wantErr:  false,
		},
		{
			name: "CaseSensitiveSortWithoutAlias",
			sort: domain.Sort{
				Field:         "email",
				CaseSensitive: true,
				Order:         domain.Asc,
			},
			want:    "LOWER(email) ASC",
			wantErr: false,
		},
		{
			name: "CaseSensitiveSortWithAlias",
			sort: domain.Sort{
				Field:         "name",
				CaseSensitive: true,
				Order:         domain.Desc,
			},
			optAlias: []string{"people"},
			want:     "LOWER(people.name) DESC",
			wantErr:  false,
		},
		{
			name: "EmptyFieldName",
			sort: domain.Sort{
				Field:         "",
				CaseSensitive: false,
				Order:         domain.Asc,
			},
			wantErr: true,
		},
		{
			name: "InvalidOrder",
			sort: domain.Sort{
				Field:         "username",
				CaseSensitive: false,
				Order:         domain.SortOrder("invalid"),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := StringifySort(tt.sort, tt.optAlias...)

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, gotErr)
		})
	}
}
