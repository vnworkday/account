package repository

import (
	"strings"
	"testing"

	"github.com/vnworkday/account/internal/fixture"

	"github.com/pkg/errors"
)

func TestMatcherBuilder_ThenDoNothing(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		matched bool
		want    string
	}{
		{
			name:    "ThenDoNothing When Matched",
			matched: true,
			want:    " THEN DO NOTHING",
		},
		{
			name:    "ThenDoNothing When Not Matched",
			matched: false,
			want:    " THEN DO NOTHING",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mb := &MutationBuilder[string]{matchClause: strings.Builder{}, notMatchClause: strings.Builder{}}
			matcher := NewMatcher(mb, tt.matched).ThenDoNothing()

			var got string
			if tt.matched {
				got = matcher.matchClause.String()
			} else {
				got = matcher.notMatchClause.String()
			}

			if err := fixture.ExpectationsWereMet(tt.want, got, false, matcher.err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMatcherBuilder_ThenDelete(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		matched   bool
		want      string
		wantError bool
	}{
		{
			name:      "ThenDelete When Matched",
			matched:   true,
			want:      " THEN DELETE",
			wantError: false,
		},
		{
			name:      "ThenDelete When Not Matched",
			matched:   false,
			want:      " THEN DELETE",
			wantError: false,
		},
		{
			name:      "ThenDelete With Error",
			matched:   true,
			want:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mb := &MutationBuilder[string]{matchClause: strings.Builder{}, notMatchClause: strings.Builder{}}
			matcher := NewMatcher(mb, tt.matched)

			if tt.wantError {
				matcher.err = errors.New("test error")
			}

			matcher.ThenDelete()

			var got string
			if tt.matched {
				got = matcher.mb.matchClause.String()
			} else {
				got = matcher.mb.notMatchClause.String()
			}

			if err := fixture.ExpectationsWereMet(tt.want, got, tt.wantError, matcher.err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMatcherBuilder_ThenUpdateOrInsert_Columns(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		matched    bool
		columns    []string
		want       string
		wantError  bool
		thenAction string
	}{
		{
			name:       "ThenUpdate With Matched And Single Column",
			matched:    true,
			columns:    []string{"name"},
			want:       " THEN UPDATE SET name = source.name",
			wantError:  false,
			thenAction: matchActionUpdate,
		},
		{
			name:       "ThenUpdate With Matched And Multiple Columns",
			matched:    true,
			columns:    []string{"name", "age"},
			want:       " THEN UPDATE SET name = source.name, age = source.age",
			wantError:  false,
			thenAction: matchActionUpdate,
		},
		{
			name:       "ThenUpdate With Not Matched",
			matched:    false,
			columns:    []string{"name"},
			want:       "",
			wantError:  true,
			thenAction: matchActionUpdate,
		},
		{
			name:       "ThenUpdate With Matched And No Columns",
			matched:    true,
			columns:    []string{},
			want:       " THEN DO NOTHING",
			wantError:  false,
			thenAction: matchActionUpdate,
		},
		{
			name:       "ThenInsert With Matched And Single Column",
			matched:    false,
			columns:    []string{"name"},
			want:       " THEN INSERT (name) VALUES (source.name)",
			wantError:  false,
			thenAction: matchActionInsert,
		},
		{
			name:       "ThenInsert With Matched And Multiple Columns",
			matched:    false,
			columns:    []string{"name", "age"},
			want:       " THEN INSERT (name, age) VALUES (source.name, source.age)",
			wantError:  false,
			thenAction: matchActionInsert,
		},
		{
			name:       "ThenInsert With Matched",
			matched:    true,
			columns:    []string{"name"},
			want:       "",
			wantError:  true,
			thenAction: matchActionInsert,
		},
		{
			name:       "ThenInsert With Matched And No Columns",
			matched:    false,
			columns:    []string{},
			want:       " THEN DO NOTHING",
			wantError:  false,
			thenAction: matchActionInsert,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mb := &MutationBuilder[string]{matchClause: strings.Builder{}, notMatchClause: strings.Builder{}}

			if tt.thenAction == matchActionUpdate {
				mb = NewMatcher(mb, tt.matched).ThenUpdate(tt.columns...)
			} else {
				mb = NewMatcher(mb, tt.matched).ThenInsert(tt.columns...)
			}

			var got string
			if tt.thenAction == matchActionUpdate {
				got = mb.matchClause.String()
			} else {
				got = mb.notMatchClause.String()
			}

			if err := fixture.ExpectationsWereMet(tt.want, got, tt.wantError, mb.err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestMatcherBuilder_ThenUpdateOrInsert(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		matched    bool
		setters    []Setter
		want       string
		wantError  bool
		thenAction string
	}{
		{
			name:       "thenUpdate Matched Single Setter",
			matched:    true,
			setters:    []Setter{{Field: "name", Value: "source.name"}},
			want:       " THEN UPDATE SET name = source.name",
			wantError:  false,
			thenAction: matchActionUpdate,
		},
		{
			name:       "thenUpdate Matched Multiple Setters",
			matched:    true,
			setters:    []Setter{{Field: "name", Value: "source.name"}, {Field: "age", Value: "source.age"}},
			want:       " THEN UPDATE SET name = source.name, age = source.age",
			wantError:  false,
			thenAction: matchActionUpdate,
		},
		{
			name:       "thenUpdate NotMatched",
			matched:    false,
			setters:    []Setter{{Field: "name", Value: "source.name"}},
			want:       "",
			wantError:  true,
			thenAction: matchActionUpdate,
		},
		{
			name:       "thenInsert NotMatched Single Setter",
			matched:    false,
			setters:    []Setter{{Field: "name", Value: "source.name"}},
			want:       " THEN INSERT (name) VALUES (source.name)",
			wantError:  false,
			thenAction: matchActionInsert,
		},
		{
			name:       "thenInsert NotMatched Multiple Setters",
			matched:    false,
			setters:    []Setter{{Field: "name", Value: "source.name"}, {Field: "age", Value: "source.age"}},
			want:       " THEN INSERT (name, age) VALUES (source.name, source.age)",
			wantError:  false,
			thenAction: matchActionInsert,
		},
		{
			name:       "thenInsert Matched Error",
			matched:    true,
			setters:    []Setter{{Field: "name", Value: "source.name"}},
			want:       "",
			wantError:  true,
			thenAction: matchActionInsert,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mb := &MutationBuilder[string]{matchClause: strings.Builder{}, notMatchClause: strings.Builder{}}

			if tt.thenAction == matchActionUpdate {
				mb = NewMatcher(mb, tt.matched).thenUpdate(tt.setters...)
			} else {
				mb = NewMatcher(mb, tt.matched).thenInsert(tt.setters...)
			}

			var got string
			if tt.thenAction == matchActionUpdate {
				got = mb.matchClause.String()
			} else {
				got = mb.notMatchClause.String()
			}

			if err := fixture.ExpectationsWereMet(tt.want, got, tt.wantError, mb.err); err != nil {
				t.Error(err)
			}
		})
	}
}
