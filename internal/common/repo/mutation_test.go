package repo

import (
	"testing"


	"github.com/vnworkday/account/internal/common/domain"

	"github.com/pkg/errors"
	"github.com/vnworkday/account/internal/fixture"
)

func TestMutationBuilder_MergeInto(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		table   string
		wantErr bool
		want    string
	}{
		{
			name:    "MergeInto With Valid Table",
			table:   "users",
			wantErr: false,
			want:    "MERGE INTO users AS target",
		},
		{
			name:    "MergeInto With Empty Table",
			table:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := NewMutationBuilder[string]()
			gotBuilder := builder.MergeInto(tt.table)

			got := gotBuilder.mergeClause
			gotErr := gotBuilder.err

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, gotErr)
		})
	}
}

func TestMutationBuilder_UsingValues(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setters []Setter
		want    map[string]any
		wantErr bool
	}{
		{
			name:    "UsingValues With Single Setter",
			setters: []Setter{{Field: "name", Value: "John Doe"}},
			want: map[string]any{
				"clause": "USING (VALUES (?)) AS source (name)",
				"keys":   []string{"name"},
				"values": []any{"John Doe"},
			},
			wantErr: false,
		},
		{
			name:    "UsingValues With Multiple Setters",
			setters: []Setter{{Field: "name", Value: "John Doe"}, {Field: "age", Value: 30}},
			want: map[string]any{
				"clause": "USING (VALUES (?, ?)) AS source (name, age)",
				"keys":   []string{"name", "age"},
				"values": []any{"John Doe", 30},
			},
			wantErr: false,
		},
		{
			name:    "UsingValues With No Setters",
			setters: []Setter{},
			want: map[string]any{
				"clause": "",
				"keys":   []string{},
				"values": []any{},
			},
			wantErr: true,
		},
		{
			name:    "UsingValues With Nil Setters",
			setters: nil,
			want: map[string]any{
				"clause": "",
				"keys":   []string{},
				"values": []any{},
			},
			wantErr: true,
		},
		{
			name:    "UsingValues With Complex Value",
			setters: []Setter{{Field: "data", Value: map[string]any{"key": "value"}}},
			want: map[string]any{
				"clause": "USING (VALUES (?)) AS source (data)",
				"keys":   []string{"data"},
				"values": []any{map[string]any{"key": "value"}},
			},
			wantErr: false,
		},
		{
			name:    "UsingValues With Setter Having Empty Field",
			setters: []Setter{{Field: "", Value: "John Doe"}},
			want: map[string]any{
				"clause": "",
				"keys":   []string{},
				"values": []any{},
			},
			wantErr: true,
		},
		{
			name:    "UsingValues With Setter Having Nil Value",
			setters: []Setter{{Field: "name", Value: nil}},
			want: map[string]any{
				"clause": "USING (VALUES (?)) AS source (name)",
				"keys":   []string{"name"},
				"values": []any{nil},
			},
			wantErr: false,
		},
		{
			name:    "UsingValues With Multiple Setters Including Empty Field",
			setters: []Setter{{Field: "name", Value: "John Doe"}, {Field: "", Value: 30}},
			want: map[string]any{
				"clause": "",
				"keys":   []string{},
				"values": []any{},
			},
			wantErr: true,
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

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, gotErr)
		})
	}
}

func TestMutationBuilder_On(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mergeCond MergeCondition
		want      string
		wantErr   bool
	}{
		{
			name: "On With Valid Condition",
			mergeCond: MergeCondition{
				SourceCol:     "id",
				TargetCol:     "id",
				Op:            domain.Eq,
				CaseSensitive: false,
			},
			want:    "ON source.id = target.id",
			wantErr: false,
		},
		{
			name: "On With Case Sensitive Condition",
			mergeCond: MergeCondition{
				SourceCol:     "name",
				TargetCol:     "name",
				Op:            domain.Eq,
				CaseSensitive: true,
			},
			want:    "ON LOWER(source.name) = LOWER(target.name)",
			wantErr: false,
		},
		{
			name: "On With Invalid Operator",
			mergeCond: MergeCondition{
				SourceCol:     "id",
				TargetCol:     "id",
				Op:            domain.Op(999),
				CaseSensitive: false,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "On With Empty Source Column",
			mergeCond: MergeCondition{
				SourceCol:     "",
				TargetCol:     "id",
				Op:            domain.Eq,
				CaseSensitive: false,
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "On With Empty Target Column",
			mergeCond: MergeCondition{
				SourceCol:     "id",
				TargetCol:     "",
				Op:            domain.Eq,
				CaseSensitive: false,
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := NewMutationBuilder[string]()
			gotBuilder := builder.On(tt.mergeCond)

			got := gotBuilder.onClause.String()
			gotErr := gotBuilder.err

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, gotErr)
		})
	}
}

func TestMutationBuilder_OnRaw(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mergeCond string
		want      string
		wantErr   bool
	}{
		{
			name:      "onRaw With Valid Condition",
			mergeCond: "target.id = source.id",
			want:      "ON target.id = source.id",
			wantErr:   false,
		},
		{
			name:      "onRaw With Empty Condition",
			mergeCond: "",
			want:      "ON ",
			wantErr:   true,
		},
		{
			name:      "onRaw With Multiple Calls",
			mergeCond: "target.id = source.id AND target.name = source.name",
			want:      "ON target.id = source.id AND target.name = source.name",
			wantErr:   false,
		},
		{
			name:      "onRaw With Special Characters",
			mergeCond: "target.name = 'O\\'Reilly'",
			want:      "ON target.name = 'O\\'Reilly'",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := NewMutationBuilder[string]()
			gotBuilder := builder.onRaw(tt.mergeCond)

			got := gotBuilder.onClause.String()
			gotErr := gotBuilder.err

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, gotErr)
		})
	}
}

func TestMutationBuilder_WhenMatchedOrNot(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		conditions []string
		want       string
		wantErr    bool
		matched    bool
	}{
		{
			name:       "WhenNotMatched Without Conditions",
			conditions: []string{},
			want:       "WHEN NOT MATCHED",
			wantErr:    false,
			matched:    false,
		},
		{
			name:       "WhenNotMatched With Single Condition",
			conditions: []string{"target.id IS NULL"},
			want:       "WHEN NOT MATCHED AND target.id IS NULL",
			wantErr:    false,
			matched:    false,
		},
		{
			name:       "WhenNotMatched With Multiple Conditions",
			conditions: []string{"target.id IS NULL", "target.name IS NULL"},
			want:       "WHEN NOT MATCHED AND target.id IS NULL AND target.name IS NULL",
			wantErr:    false,
			matched:    false,
		},
		{
			name:       "WhenNotMatched Preceded By Error",
			conditions: []string{"target.id IS NULL"},
			want:       "",
			wantErr:    true,
			matched:    false,
		},
		{
			name:       "WhenMatched Without Conditions",
			conditions: []string{},
			want:       "WHEN MATCHED",
			wantErr:    false,
			matched:    true,
		},
		{
			name:       "WhenMatched With Single Condition",
			conditions: []string{"target.id = source.id"},
			want:       "WHEN MATCHED AND target.id = source.id",
			wantErr:    false,
			matched:    true,
		},
		{
			name:       "WhenMatched With Multiple Conditions",
			conditions: []string{"target.id = source.id", "target.name = source.name"},
			want:       "WHEN MATCHED AND target.id = source.id AND target.name = source.name",
			wantErr:    false,
			matched:    true,
		},
		{
			name:       "WhenMatched Preceded By Error",
			conditions: []string{"target.id = source.id"},
			want:       "",
			wantErr:    true,
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

			if tt.wantErr {
				gotBuilder.err = errors.New("forced error")
			}

			var got string
			gotErr := gotBuilder.err

			if !tt.matched {
				got = gotBuilder.mb.notMatchClause.String()
			} else {
				got = gotBuilder.mb.matchClause.String()
			}

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, gotErr)
		})
	}
}

func TestMutationBuilder_Build(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		setup   func(b *MutationBuilder[string])
		want    string
		wantErr bool
	}{
		{
			name: "Build Without When Clauses",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        domain.Eq,
					})
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Build Without MergeInto Clause",
			setup: func(b *MutationBuilder[string]) {
				b.UsingValues(Setter{Field: "id", Value: 1}).
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        domain.Eq,
					})
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Build Without UsingValues Clause",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        domain.Eq,
					})
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Build Without onRaw Clause",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1})
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Build With All Clauses",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        domain.Eq,
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
			wantErr: false,
		},
		{
			name: "Build With Preceding Error",
			setup: func(b *MutationBuilder[string]) {
				b.err = errors.New("forced error")
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Build With Complex UsingValues",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "data", Value: map[string]any{"key": "value"}}).
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        domain.Eq,
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
			wantErr: false,
		},
		{
			name: "Build With Multiple onRaw Conditions",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        domain.Eq,
					}).
					On(MergeCondition{
						SourceCol: "name",
						TargetCol: "name",
						Op:        domain.Eq,
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
			wantErr: false,
		},
		{
			name: "Build With Multiple WhenMatched Conditions",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        domain.Eq,
					}).
					WhenMatched("target.name = source.name", "target.age = source.age").
					ThenDelete()
			},
			want: "MERGE INTO users AS target " +
				"USING (VALUES (?)) AS source (id) " +
				"ON source.id = target.id " +
				"WHEN MATCHED AND target.name = source.name AND target.age = source.age " +
				"THEN DELETE",
			wantErr: false,
		},
		{
			name: "Build With WhenNotMatched And WhenMatched",
			setup: func(b *MutationBuilder[string]) {
				b.MergeInto("users").
					UsingValues(Setter{Field: "id", Value: 1}).
					On(MergeCondition{
						SourceCol: "id",
						TargetCol: "id",
						Op:        domain.Eq,
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
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			builder := NewMutationBuilder[string]()
			tt.setup(builder)

			got, gotErr := builder.build()

			fixture.ExpectationsWereMet(t, tt.want, got, tt.wantErr, gotErr)
		})
	}
}
