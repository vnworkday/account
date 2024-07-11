package repository

import (
	"testing"

	"github.com/vnworkday/account/internal/fixture"
	"github.com/vnworkday/account/internal/model"

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

			got := gotBuilder.mergeClause
			gotErr := gotBuilder.err

			if err := fixture.ExpectationsWereMet(tt.want, got, tt.wantError, gotErr); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMutationBuilder_UsingValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		setters   []Setter
		want      map[string]any
		wantError bool
	}{
		{
			name:    "UsingValues With Single Setter",
			setters: []Setter{{Field: "name", Value: "John Doe"}},
			want: map[string]any{
				"clause": "USING (VALUES (?)) AS source (name)",
				"keys":   []string{"name"},
				"values": []any{"John Doe"},
			},
			wantError: false,
		},
		{
			name:    "UsingValues With Multiple Setters",
			setters: []Setter{{Field: "name", Value: "John Doe"}, {Field: "age", Value: 30}},
			want: map[string]any{
				"clause": "USING (VALUES (?, ?)) AS source (name, age)",
				"keys":   []string{"name", "age"},
				"values": []any{"John Doe", 30},
			},
			wantError: false,
		},
		{
			name:    "UsingValues With No Setters",
			setters: []Setter{},
			want: map[string]any{
				"clause": "",
				"keys":   []string{},
				"values": []any{},
			},
			wantError: true,
		},
		{
			name:    "UsingValues With Nil Setters",
			setters: nil,
			want: map[string]any{
				"clause": "",
				"keys":   []string{},
				"values": []any{},
			},
			wantError: true,
		},
		{
			name:    "UsingValues With Complex Value",
			setters: []Setter{{Field: "data", Value: map[string]any{"key": "value"}}},
			want: map[string]any{
				"clause": "USING (VALUES (?)) AS source (data)",
				"keys":   []string{"data"},
				"values": []any{map[string]any{"key": "value"}},
			},
			wantError: false,
		},
		{
			name:    "UsingValues With Setter Having Empty Field",
			setters: []Setter{{Field: "", Value: "John Doe"}},
			want: map[string]any{
				"clause": "",
				"keys":   []string{},
				"values": []any{},
			},
			wantError: true,
		},
		{
			name:    "UsingValues With Setter Having Nil Value",
			setters: []Setter{{Field: "name", Value: nil}},
			want: map[string]any{
				"clause": "USING (VALUES (?)) AS source (name)",
				"keys":   []string{"name"},
				"values": []any{nil},
			},
			wantError: false,
		},
		{
			name:    "UsingValues With Multiple Setters Including Empty Field",
			setters: []Setter{{Field: "name", Value: "John Doe"}, {Field: "", Value: 30}},
			want: map[string]any{
				"clause": "",
				"keys":   []string{},
				"values": []any{},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := NewMutationBuilder[string]()
			gotBuilder := builder.UsingValues(tt.setters...)

			got := map[string]any{
				"clause": gotBuilder.usingClause,
				"keys":   gotBuilder.usingArgKeys,
				"values": gotBuilder.usingArgValues,
			}
			gotErr := gotBuilder.err

			if err := fixture.ExpectationsWereMet(tt.want, got, tt.wantError, gotErr); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMutationBuilder_On(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mergeCond MergeCondition
		want      string
		wantError bool
	}{
		{
			name: "On With Valid Condition",
			mergeCond: MergeCondition{
				SourceCol:     "id",
				TargetCol:     "id",
				Op:            model.Eq,
				CaseSensitive: false,
			},
			want:      "ON source.id = target.id",
			wantError: false,
		},
		{
			name: "On With Case Sensitive Condition",
			mergeCond: MergeCondition{
				SourceCol:     "name",
				TargetCol:     "name",
				Op:            model.Eq,
				CaseSensitive: true,
			},
			want:      "ON LOWER(source.name) = LOWER(target.name)",
			wantError: false,
		},
		{
			name: "On With Invalid Operator",
			mergeCond: MergeCondition{
				SourceCol:     "id",
				TargetCol:     "id",
				Op:            model.Op(999),
				CaseSensitive: false,
			},
			want:      "",
			wantError: true,
		},
		{
			name: "On With Empty Source Column",
			mergeCond: MergeCondition{
				SourceCol:     "",
				TargetCol:     "id",
				Op:            model.Eq,
				CaseSensitive: false,
			},
			want:      "",
			wantError: true,
		},
		{
			name: "On With Empty Target Column",
			mergeCond: MergeCondition{
				SourceCol:     "id",
				TargetCol:     "",
				Op:            model.Eq,
				CaseSensitive: false,
			},
			want:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := NewMutationBuilder[string]()
			gotBuilder := builder.On(tt.mergeCond)

			got := gotBuilder.onClause.String()
			gotErr := gotBuilder.err

			if err := fixture.ExpectationsWereMet(tt.want, got, tt.wantError, gotErr); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMutationBuilder_OnRaw(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mergeCond string
		want      string
		wantError bool
	}{
		{
			name:      "onRaw With Valid Condition",
			mergeCond: "target.id = source.id",
			want:      "ON target.id = source.id",
			wantError: false,
		},
		{
			name:      "onRaw With Empty Condition",
			mergeCond: "",
			want:      "ON ",
			wantError: true,
		},
		{
			name:      "onRaw With Multiple Calls",
			mergeCond: "target.id = source.id AND target.name = source.name",
			want:      "ON target.id = source.id AND target.name = source.name",
			wantError: false,
		},
		{
			name:      "onRaw With Special Characters",
			mergeCond: "target.name = 'O\\'Reilly'",
			want:      "ON target.name = 'O\\'Reilly'",
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := NewMutationBuilder[string]()
			gotBuilder := builder.onRaw(tt.mergeCond)

			got := gotBuilder.onClause.String()
			gotErr := gotBuilder.err

			if err := fixture.ExpectationsWereMet(tt.want, got, tt.wantError, gotErr); err != nil {
				t.Error(err)
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

			var gotBuilder *MatcherBuilder[string]

			if !tt.matched {
				gotBuilder = builder.WhenNotMatched(tt.conditions...)
			} else {
				gotBuilder = builder.WhenMatched(tt.conditions...)
			}

			if tt.wantError {
				gotBuilder.err = errors.New("forced error")
			}

			var got string
			gotErr := gotBuilder.err

			if !tt.matched {
				got = gotBuilder.mb.notMatchClause.String()
			} else {
				got = gotBuilder.mb.matchClause.String()
			}

			if err := fixture.ExpectationsWereMet(tt.want, got, tt.wantError, gotErr); err != nil {
				t.Error(err)
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
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        model.Eq,
					})
			},
			want:      "",
			wantError: true,
		},
		{
			name: "Build Without MergeInto Clause",
			setup: func(b *MutationBuilder[string]) {
				b.UsingValues(Setter{Field: "id", Value: 1}).
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        model.Eq,
					})
			},
			want:      "",
			wantError: true,
		},
		{
			name: "Build Without UsingValues Clause",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        model.Eq,
					})
			},
			want:      "",
			wantError: true,
		},
		{
			name: "Build Without onRaw Clause",
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
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        model.Eq,
					}).
					WhenMatched("target.name = source.name").
					ThenDelete().
					WhenNotMatched().
					ThenDoNothing()
			},
			want: "MERGE INTO users AS target " +
				"USING (VALUES (?)) AS source (id) " +
				"ON source.id = target.id " +
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
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        model.Eq,
					}).
					WhenMatched("target.name = source.name").
					ThenDelete().
					WhenNotMatched().
					ThenDoNothing()
			},
			want: "MERGE INTO users AS target " +
				"USING (VALUES (?)) AS source (data) " +
				"ON source.id = target.id " +
				"WHEN NOT MATCHED THEN DO NOTHING " +
				"WHEN MATCHED AND target.name = source.name THEN DELETE",
			wantError: false,
		},
		{
			name: "Build With Multiple onRaw Conditions",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        model.Eq,
					}).
					On(MergeCondition{
						SourceCol: "name",
						TargetCol: "name",
						Op:        model.Eq,
					}).
					WhenMatched("target.name = source.name").
					ThenDelete().
					WhenNotMatched().
					ThenDoNothing()
			},
			want: "MERGE INTO users AS target " +
				"USING (VALUES (?)) AS source (id) " +
				"ON source.id = target.id AND source.name = target.name " +
				"WHEN NOT MATCHED THEN DO NOTHING " +
				"WHEN MATCHED AND target.name = source.name THEN DELETE",
			wantError: false,
		},
		{
			name: "Build With Multiple WhenMatched Conditions",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        model.Eq,
					}).
					WhenMatched("target.name = source.name", "target.age = source.age").
					ThenDelete()
			},
			want: "MERGE INTO users AS target " +
				"USING (VALUES (?)) AS source (id) " +
				"ON source.id = target.id " +
				"WHEN MATCHED AND target.name = source.name AND target.age = source.age " +
				"THEN DELETE",
			wantError: false,
		},
		{
			name: "Build With WhenNotMatched And WhenMatched",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        model.Eq,
					}).
					WhenNotMatched().
					ThenDoNothing().
					WhenMatched("target.name = source.name").
					ThenDoNothing()
			},
			want: "MERGE INTO users AS target " +
				"USING (VALUES (?)) AS source (id) " +
				"ON source.id = target.id " +
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

			got, gotErr := builder.build()

			if err := fixture.ExpectationsWereMet(tt.want, got, tt.wantError, gotErr); err != nil {
				t.Error(err)
			}
		})
	}
}
