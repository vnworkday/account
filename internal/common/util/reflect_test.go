package util

import (
	"testing"

	"github.com/vnworkday/account/internal/common/fixture"
)

func TestType(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Name string
	}

	tests := []struct {
		name  string
		input any
		want  string
	}{
		{"WithNilValue", nil, "nil"},
		{"WithPrimitiveType", 42, "int"},
		{"WithAnonymousStruct", struct{ Name string }{}, "struct"},
		{"WithNamedStruct", testStruct{}, "testStruct"},
		{"WithPointerToPrimitive", new(int), "*int"},
		{"WithPointerToAnonymousStruct", &struct{ Name string }{}, "*struct"},
		{"WithPointerToNamedStruct", &testStruct{}, "*testStruct"},
		{"WithSlice", []int{1, 2, 3}, "slice"},
		{"WithMap", map[string]int{"key": 42}, "map"},
		{"WithInterface", "test", "string"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			result := Type(tt.input)
			fixture.ExpectationsWereMet(t, tt.want, result, false, nil)
		})
	}
}
