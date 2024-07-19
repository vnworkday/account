package domain

import (
	"testing"

	"github.com/vnworkday/account/internal/fixture"
)

func TestStructToTable(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   any
		want    *Table
		wantErr bool
	}{
		{
			name: "ValidStructPointerWithTags",
			input: &struct {
				ID        int    `db:"pid"`
				FirstName string `db:"firstname"`
			}{},
			want: &Table{
				Columns:    []string{"pid", "firstname"},
				Insertable: []string{"pid", "firstname"},
				Updatable:  []string{"pid", "firstname"},
			},
			wantErr: false,
		},
		{
			name: "ValidStructWithTags",
			input: struct {
				ID        int    `db:"pid"`
				FirstName string `db:"firstname"`
			}{},
			want: &Table{
				Columns:    []string{"pid", "firstname"},
				Insertable: []string{"pid", "firstname"},
				Updatable:  []string{"pid", "firstname"},
			},
			wantErr: false,
		},
		{
			name: "ValidStructWithTagsOption",
			input: struct {
				ID        int    `db:"pid,generated"`
				FirstName string `db:"firstname,immutable"`
				LastName  string `db:"lastname,generated,immutable"`
			}{},
			want: &Table{
				Columns:    []string{"pid", "firstname", "lastname"},
				Insertable: []string{"firstname"},
				Updatable:  []string{"pid"},
			},
			wantErr: false,
		},
		{
			name: "StructPointerWithoutTags",
			input: &struct {
				ID        int
				FirstName string
			}{},
			want:    nil,
			wantErr: true,
		},
		{
			name: "StructWithoutTags",
			input: struct {
				ID        int
				FirstName string
			}{},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "NilInput",
			input:   nil,
			want:    nil,
			wantErr: true,
		},
		{
			name:    "NeitherStructOrStructPointer",
			input:   123,
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, gotErr := StructToTable(tt.input, "table")

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, gotErr)
		})
	}
}
