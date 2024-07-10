package repository

import (
	"reflect"
	"testing"

	"github.com/pkg/errors"
)

func TestMutationBuilder_MergeInto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		table     string
		wantError bool
		want      string
	}{
		{
			name:      "MergeInto With Valid Table",
			table:     "users",
			wantError: false,
			want:      "MERGE INTO users AS target",
		},
		{
			name:      "MergeInto With Empty Table",
			table:     "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := NewMutationBuilder[string]()
			gotBuilder := builder.MergeInto(tt.table)

			if tt.wantError {
				if gotBuilder.err == nil {
					t.Errorf("Expected an error but got nil")
				}
			} else {
				if gotBuilder.err != nil {
					t.Errorf("Did not expect an error but got one: %v", gotBuilder.err)
				}

				if gotBuilder.mergeClause != tt.want {
					t.Errorf("MergeClause = %v, want %v", gotBuilder.mergeClause, tt.want)
				}
			}
		})
	}
}

func TestMutationBuilder_UsingValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name            string
		setters         []Setter
		wantUsingClause string
		wantArgKeys     []string
		wantArgValues   []any
		wantError       bool
	}{
		{
			name:            "UsingValues With Single Setter",
			setters:         []Setter{{Field: "name", Value: "John Doe"}},
			wantUsingClause: "USING (VALUES (?)) AS source (name)",
			wantArgKeys:     []string{"name"},
			wantArgValues:   []any{"John Doe"},
			wantError:       false,
		},
		{
			name:            "UsingValues With Multiple Setters",
			setters:         []Setter{{Field: "name", Value: "John Doe"}, {Field: "age", Value: 30}},
			wantUsingClause: "USING (VALUES (?, ?)) AS source (name, age)",
			wantArgKeys:     []string{"name", "age"},
			wantArgValues:   []any{"John Doe", 30},
			wantError:       false,
		},
		{
			name:            "UsingValues With No Setters",
			setters:         []Setter{},
			wantUsingClause: "",
			wantArgKeys:     []string{},
			wantArgValues:   []any{},
			wantError:       true,
		},
		{
			name:            "UsingValues With Nil Setters",
			setters:         nil,
			wantUsingClause: "",
			wantArgKeys:     []string{},
			wantArgValues:   []any{},
			wantError:       true,
		},
		{
			name:            "UsingValues With Complex Value",
			setters:         []Setter{{Field: "data", Value: map[string]any{"key": "value"}}},
			wantUsingClause: "USING (VALUES (?)) AS source (data)",
			wantArgKeys:     []string{"data"},
			wantArgValues:   []any{map[string]any{"key": "value"}},
			wantError:       false,
		},
		{
			name:            "UsingValues With Setter Having Empty Field",
			setters:         []Setter{{Field: "", Value: "John Doe"}},
			wantUsingClause: "",
			wantArgKeys:     []string{},
			wantArgValues:   []any{},
			wantError:       true,
		},
		{
			name:            "UsingValues With Setter Having Nil Value",
			setters:         []Setter{{Field: "name", Value: nil}},
			wantUsingClause: "USING (VALUES (?)) AS source (name)",
			wantArgKeys:     []string{"name"},
			wantArgValues:   []any{nil},
			wantError:       false,
		},
		{
			name:            "UsingValues With Multiple Setters Including Empty Field",
			setters:         []Setter{{Field: "name", Value: "John Doe"}, {Field: "", Value: 30}},
			wantUsingClause: "",
			wantArgKeys:     []string{},
			wantArgValues:   []any{},
			wantError:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := NewMutationBuilder[string]()
			gotBuilder := builder.UsingValues(tt.setters...)

			if tt.wantError {
				if gotBuilder.err == nil {
					t.Errorf("Expected an error but got nil")
				}

				return
			}

			if gotBuilder.err != nil {
				t.Errorf("Did not expect an error but got one: %v", gotBuilder.err)
			}

			if gotBuilder.usingClause != tt.wantUsingClause {
				t.Errorf("UsingClause = %v, want %v", gotBuilder.usingClause, tt.wantUsingClause)
			}

			if !reflect.DeepEqual(gotBuilder.usingArgKeys, tt.wantArgKeys) {
				t.Errorf("UsingArgKeys = %v, want %v", gotBuilder.usingArgKeys, tt.wantArgKeys)
			}

			if !reflect.DeepEqual(gotBuilder.usingArgValues, tt.wantArgValues) {
				t.Errorf("UsingArgValues = %v, want %v", gotBuilder.usingArgValues, tt.wantArgValues)
			}
		})
	}
}

func TestMutationBuilder_OnRaw(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name         string
		mergeCond    string
		wantOnClause string
		wantError    bool
	}{
		{
			name:         "OnRaw With Valid Condition",
			mergeCond:    "target.id = source.id",
			wantOnClause: "ON target.id = source.id",
			wantError:    false,
		},
		{
			name:         "OnRaw With Empty Condition",
			mergeCond:    "",
			wantOnClause: "ON ",
			wantError:    true,
		},
		{
			name:         "OnRaw With Multiple Calls",
			mergeCond:    "target.id = source.id AND target.name = source.name",
			wantOnClause: "ON target.id = source.id AND target.name = source.name",
			wantError:    false,
		},
		{
			name:         "OnRaw With Special Characters",
			mergeCond:    "target.name = 'O\\'Reilly'",
			wantOnClause: "ON target.name = 'O\\'Reilly'",
			wantError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := NewMutationBuilder[string]()
			gotBuilder := builder.OnRaw(tt.mergeCond)

			if tt.wantError {
				if gotBuilder.err == nil {
					t.Error("Expected an error but got nil")
				}
			} else {
				if gotBuilder.err != nil {
					t.Errorf("Did not expect an error but got one: %v", gotBuilder.err)
				}

				if gotBuilder.onClause.String() != tt.wantOnClause {
					t.Errorf("OnClause = %v, want %v", gotBuilder.onClause.String(), tt.wantOnClause)
				}
			}
		})
	}
}

func TestMutationBuilder_WhenMatchedOrNot(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		conditions []string
		want       string
		wantError  bool
		matched    bool
	}{
		{
			name:       "WhenNotMatched Without Conditions",
			conditions: []string{},
			want:       "WHEN NOT MATCHED",
			wantError:  false,
			matched:    false,
		},
		{
			name:       "WhenNotMatched With Single Condition",
			conditions: []string{"target.id IS NULL"},
			want:       "WHEN NOT MATCHED AND target.id IS NULL",
			wantError:  false,
			matched:    false,
		},
		{
			name:       "WhenNotMatched With Multiple Conditions",
			conditions: []string{"target.id IS NULL", "target.name IS NULL"},
			want:       "WHEN NOT MATCHED AND target.id IS NULL AND target.name IS NULL",
			wantError:  false,
			matched:    false,
		},
		{
			name:       "WhenNotMatched Preceded By Error",
			conditions: []string{"target.id IS NULL"},
			want:       "",
			wantError:  true,
			matched:    false,
		},
		{
			name:       "WhenMatched Without Conditions",
			conditions: []string{},
			want:       "WHEN MATCHED",
			wantError:  false,
			matched:    true,
		},
		{
			name:       "WhenMatched With Single Condition",
			conditions: []string{"target.id = source.id"},
			want:       "WHEN MATCHED AND target.id = source.id",
			wantError:  false,
			matched:    true,
		},
		{
			name:       "WhenMatched With Multiple Conditions",
			conditions: []string{"target.id = source.id", "target.name = source.name"},
			want:       "WHEN MATCHED AND target.id = source.id AND target.name = source.name",
			wantError:  false,
			matched:    true,
		},
		{
			name:       "WhenMatched Preceded By Error",
			conditions: []string{"target.id = source.id"},
			want:       "",
			wantError:  true,
			matched:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := NewMutationBuilder[string]()
			if tt.wantError {
				builder.err = errors.New("forced error")
			}

			var gotBuilder *MatcherBuilder[string]

			if !tt.matched {
				gotBuilder = builder.WhenNotMatched(tt.conditions...)
			} else {
				gotBuilder = builder.WhenMatched(tt.conditions...)
			}

			if tt.wantError {
				if gotBuilder.err == nil {
					t.Error("Expected an error but got nil")
				}

				return
			}

			if gotBuilder.err != nil {
				t.Errorf("Did not expect an error but got one: %v", gotBuilder.err)
			}

			var got string

			if tt.matched {
				got = gotBuilder.mb.matchClause.String()
			} else {
				got = gotBuilder.mb.notMatchClause.String()
			}

			if got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMutationBuilder_Build(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setup     func(b *MutationBuilder[string])
		want      string
		wantError bool
	}{
		{
			name: "Build Without When Clauses",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					OnRaw("target.id = source.id")
			},
			want:      "",
			wantError: true,
		},
		{
			name: "Build Without MergeInto Clause",
			setup: func(b *MutationBuilder[string]) {
				b.UsingValues(Setter{Field: "id", Value: 1}).
					OnRaw("target.id = source.id")
			},
			want:      "",
			wantError: true,
		},
		{
			name: "Build Without UsingValues Clause",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					OnRaw("target.id = source.id")
			},
			want:      "",
			wantError: true,
		},
		{
			name: "Build Without OnRaw Clause",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1})
			},
			want:      "",
			wantError: true,
		},
		{
			name: "Build With All Clauses",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					OnRaw("target.id = source.id").
					WhenMatched("target.name = source.name").
					ThenDelete().
					WhenNotMatched().
					ThenDoNothing()
			},
			want: "MERGE INTO users AS target " +
				"USING (VALUES (?)) AS source (id) " +
				"ON target.id = source.id " +
				"WHEN NOT MATCHED THEN DO NOTHING " +
				"WHEN MATCHED AND target.name = source.name THEN DELETE",
			wantError: false,
		},
		{
			name: "Build With Preceding Error",
			setup: func(b *MutationBuilder[string]) {
				b.err = errors.New("forced error")
			},
			want:      "",
			wantError: true,
		},
		{
			name: "Build With Complex UsingValues",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "data", Value: map[string]any{"key": "value"}}).
					OnRaw("target.id = source.id").
					WhenMatched("target.name = source.name").
					ThenDelete().
					WhenNotMatched().
					ThenDoNothing()
			},
			want: "MERGE INTO users AS target " +
				"USING (VALUES (?)) AS source (data) " +
				"ON target.id = source.id " +
				"WHEN NOT MATCHED THEN DO NOTHING " +
				"WHEN MATCHED AND target.name = source.name THEN DELETE",
			wantError: false,
		},
		{
			name: "Build With Multiple OnRaw Conditions",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					OnRaw("target.id = source.id").
					OnRaw("target.name = source.name").
					WhenMatched("target.name = source.name").
					ThenDelete().
					WhenNotMatched().
					ThenDoNothing()
			},
			want: "MERGE INTO users AS target " +
				"USING (VALUES (?)) AS source (id) " +
				"ON target.id = source.id AND target.name = source.name " +
				"WHEN NOT MATCHED THEN DO NOTHING " +
				"WHEN MATCHED AND target.name = source.name THEN DELETE",
			wantError: false,
		},
		{
			name: "Build With Multiple WhenMatched Conditions",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					OnRaw("target.id = source.id").
					WhenMatched("target.name = source.name", "target.age = source.age").
					ThenDelete()
			},
			want: "MERGE INTO users AS target " +
				"USING (VALUES (?)) AS source (id) " +
				"ON target.id = source.id " +
				"WHEN MATCHED AND target.name = source.name AND target.age = source.age " +
				"THEN DELETE",
			wantError: false,
		},
		{
			name: "Build With WhenNotMatched And WhenMatched",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					OnRaw("target.id = source.id").
					WhenNotMatched().
					ThenDoNothing().
					WhenMatched("target.name = source.name").
					ThenDoNothing()
			},
			want: "MERGE INTO users AS target " +
				"USING (VALUES (?)) AS source (id) " +
				"ON target.id = source.id " +
				"WHEN NOT MATCHED THEN DO NOTHING " +
				"WHEN MATCHED AND target.name = source.name THEN DO NOTHING",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := NewMutationBuilder[string]()
			tt.setup(builder)

			query, err := builder.build()

			if tt.wantError {
				if err == nil {
					t.Errorf("Expected error %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}

				if query != tt.want {
					t.Errorf("Expected query to be:\n%v\ngot:\n%v", tt.want, query)
				}
			}
		})
	}
}
