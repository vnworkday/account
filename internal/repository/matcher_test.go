package repository

import (
	"strings"
	"testing"

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

			if tt.matched && matcher.matchClause.String() != tt.want {
				t.Errorf("MatchClause = %v, want %v", matcher.matchClause.String(), tt.want)
			}

			if !tt.matched && matcher.notMatchClause.String() != tt.want {
				t.Errorf("NotMatchClause = %v, want %v", matcher.notMatchClause.String(), tt.want)
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

			if !tt.wantError {
				if tt.matched && matcher.mb.matchClause.String() != tt.want {
					t.Errorf("MatchClause = %v, want %v", matcher.mb.matchClause.String(), tt.want)
				}

				if !tt.matched && matcher.mb.notMatchClause.String() != tt.want {
					t.Errorf("NotMatchClause = %v, want %v", matcher.mb.notMatchClause.String(), tt.want)
				}
			} else if matcher.err == nil {
				t.Error("Expected error, got nil")
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
			name:       "ThenUpdate Matched Single Setter",
			matched:    true,
			setters:    []Setter{{Field: "name", Value: "source.name"}},
			want:       " THEN UPDATE SET name = source.name",
			wantError:  false,
			thenAction: "update",
		},
		{
			name:       "ThenUpdate Matched Multiple Setters",
			matched:    true,
			setters:    []Setter{{Field: "name", Value: "source.name"}, {Field: "age", Value: "source.age"}},
			want:       " THEN UPDATE SET name = source.name, age = source.age",
			wantError:  false,
			thenAction: "update",
		},
		{
			name:       "ThenUpdate NotMatched",
			matched:    false,
			setters:    []Setter{{Field: "name", Value: "source.name"}},
			want:       "",
			wantError:  true,
			thenAction: "update",
		},
		{
			name:       "ThenInsert NotMatched Single Setter",
			matched:    false,
			setters:    []Setter{{Field: "name", Value: "source.name"}},
			want:       " THEN INSERT (name) VALUES (source.name)",
			wantError:  false,
			thenAction: "insert",
		},
		{
			name:       "ThenInsert NotMatched Multiple Setters",
			matched:    false,
			setters:    []Setter{{Field: "name", Value: "source.name"}, {Field: "age", Value: "source.age"}},
			want:       " THEN INSERT (name, age) VALUES (source.name, source.age)",
			wantError:  false,
			thenAction: "insert",
		},
		{
			name:       "ThenInsert Matched Error",
			matched:    true,
			setters:    []Setter{{Field: "name", Value: "source.name"}},
			want:       "",
			wantError:  true,
			thenAction: "insert",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mb := &MutationBuilder[string]{matchClause: strings.Builder{}, notMatchClause: strings.Builder{}}

			if tt.thenAction == "update" {
				mb = NewMatcher(mb, tt.matched).ThenUpdate(tt.setters...)
			} else {
				mb = NewMatcher(mb, tt.matched).ThenInsert(tt.setters...)
			}

			if !tt.wantError {
				var got string
				if tt.thenAction == "update" {
					got = mb.matchClause.String()
				} else {
					got = mb.notMatchClause.String()
				}

				if got != tt.want {
					t.Errorf("got:\n%v\nwant:\n%v", got, tt.want)
				}
			}

			if tt.wantError && mb.err == nil {
				t.Error("Expected error, got nil")
			}
		})
	}
}
