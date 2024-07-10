package repository

import "testing"

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

			got, err := stringifyField(tt.field, tt.sensitive, tt.alias)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringifyField() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if got != tt.want {
				t.Errorf("stringifyField() got = %v, want %v", got, tt.want)
			}
		})
	}
}
