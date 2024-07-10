package repository

import (
	"reflect"
	"testing"
	"time"

	"github.com/vnworkday/account/internal/model"
)

func TestStringifyFilter(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		filter   model.Filter
		optAlias []string
		want     string
		wantErr  bool
	}{
		{
			name: "SimpleEqualityFilterWithoutAlias",
			filter: model.Filter{
				Field:           "username",
				Operator:        model.Eq,
				Value:           "john_doe",
				IsCaseSensitive: false,
			},
			want:    "username = ?",
			wantErr: false,
		},
		{
			name: "SensitiveContainsFilterWithAlias",
			filter: model.Filter{
				Field:           "email",
				Operator:        model.Contains,
				Value:           "@example.com",
				IsCaseSensitive: true,
			},
			optAlias: []string{"users"},
			want:     "LOWER(users.email) LIKE '%' || LOWER(?) || '%'",
			wantErr:  false,
		},
		{
			name: "InvalidOperator",
			filter: model.Filter{
				Field:           "status",
				Operator:        model.FilterOperator(999),
				Value:           "active",
				IsCaseSensitive: false,
			},
			wantErr: true,
		},
		{
			name: "EmptyFieldName",
			filter: model.Filter{
				Field:           "",
				Operator:        model.Eq,
				Value:           "value",
				IsCaseSensitive: false,
			},
			wantErr: true,
		},
		{
			name: "NullValueWithAlias",
			filter: model.Filter{
				Field:           "deleted_at",
				Operator:        model.Null,
				IsCaseSensitive: false,
			},
			optAlias: []string{"users"},
			want:     "users.deleted_at IS NULL",
			wantErr:  false,
		},
		{
			name: "InsensitiveStartsWithFilterWithoutAlias",
			filter: model.Filter{
				Field:           "name",
				Operator:        model.StartsWith,
				Value:           "John",
				IsCaseSensitive: false,
			},
			want:    "name LIKE ? || '%'",
			wantErr: false,
		},
		{
			name: "SensitiveEndsWithFilterWithAlias",
			filter: model.Filter{
				Field:           "email",
				Operator:        model.EndsWith,
				Value:           "@example.com",
				IsCaseSensitive: true,
			},
			optAlias: []string{"contacts"},
			want:     "LOWER(contacts.email) LIKE '%' || LOWER(?)",
			wantErr:  false,
		},
		{
			name: "InOperatorWithMultipleValues",
			filter: model.Filter{
				Field:           "status",
				Operator:        model.In,
				Value:           "active,inactive",
				IsCaseSensitive: false,
			},
			want:    "status IN (?)",
			wantErr: false,
		},
		{
			name: "NotInOperatorSensitiveWithAlias",
			filter: model.Filter{
				Field:           "type",
				Operator:        model.NotIn,
				Value:           "admin,user",
				IsCaseSensitive: true,
			},
			optAlias: []string{"users"},
			want:     "LOWER(users.type) NOT IN (?)",
			wantErr:  false,
		},
		{
			name: "UnsupportedValueType",
			filter: model.Filter{
				Field:           "created_at",
				Operator:        model.Eq,
				Value:           "2023-01-01",
				IsCaseSensitive: false,
			},
			want:    "created_at = ?",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := StringifyFilter(tt.filter, tt.optAlias...)
			if (err != nil) != tt.wantErr {
				t.Errorf("StringifyFilter() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("StringifyFilter() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringifyFilterOperator(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		operator model.FilterOperator
		want     string
		wantErr  bool
	}{
		{
			name:     "EqualityOperator",
			operator: model.Eq,
			want:     "=",
			wantErr:  false,
		},
		{
			name:     "NotEqualOperator",
			operator: model.Ne,
			want:     "<>",
			wantErr:  false,
		},
		{
			name:     "GreaterThanOperator",
			operator: model.Gt,
			want:     ">",
			wantErr:  false,
		},
		{
			name:     "LessThanOperator",
			operator: model.Lt,
			want:     "<",
			wantErr:  false,
		},
		{
			name:     "GreaterThanOrEqualOperator",
			operator: model.Ge,
			want:     ">=",
			wantErr:  false,
		},
		{
			name:     "LessThanOrEqualOperator",
			operator: model.Le,
			want:     "<=",
			wantErr:  false,
		},
		{
			name:     "InOperator",
			operator: model.In,
			want:     "IN",
			wantErr:  false,
		},
		{
			name:     "NotInOperator",
			operator: model.NotIn,
			want:     "NOT IN",
			wantErr:  false,
		},
		{
			name:     "ContainsOperator",
			operator: model.Contains,
			want:     "LIKE",
			wantErr:  false,
		},
		{
			name:     "NotContainsOperator",
			operator: model.NotContains,
			want:     "NOT LIKE",
			wantErr:  false,
		},
		{
			name:     "StartsWithOperator",
			operator: model.StartsWith,
			want:     "LIKE",
			wantErr:  false,
		},
		{
			name:     "EndsWithOperator",
			operator: model.EndsWith,
			want:     "LIKE",
			wantErr:  false,
		},
		{
			name:     "NullOperator",
			operator: model.Null,
			want:     "IS NULL",
			wantErr:  false,
		},
		{
			name:     "NotNullOperator",
			operator: model.NotNull,
			want:     "IS NOT NULL",
			wantErr:  false,
		},
		{
			name:     "BetweenOperator",
			operator: model.Between,
			want:     "BETWEEN",
			wantErr:  false,
		},
		{
			name:     "UnsupportedOperator",
			operator: model.FilterOperator(999),
			want:     "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := stringifyFilterOperator(tt.operator)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringifyFilterOperator() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("stringifyFilterOperator() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBuildFilterWildcards(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		op        model.FilterOperator
		sensitive bool
		want      string
		wantErr   bool
	}{
		{
			name:      "EqualityOperatorInsensitive",
			op:        model.Eq,
			sensitive: false,
			want:      "?",
			wantErr:   false,
		},
		{
			name:      "EqualityOperatorSensitive",
			op:        model.Eq,
			sensitive: true,
			want:      "LOWER(?)",
			wantErr:   false,
		},
		{
			name:      "ContainsOperatorInsensitive",
			op:        model.Contains,
			sensitive: false,
			want:      "'%' || ? || '%'",
			wantErr:   false,
		},
		{
			name:      "ContainsOperatorSensitive",
			op:        model.Contains,
			sensitive: true,
			want:      "'%' || LOWER(?) || '%'",
			wantErr:   false,
		},
		{
			name:      "StartsWithOperatorInsensitive",
			op:        model.StartsWith,
			sensitive: false,
			want:      "? || '%'",
			wantErr:   false,
		},
		{
			name:      "StartsWithOperatorSensitive",
			op:        model.StartsWith,
			sensitive: true,
			want:      "LOWER(?) || '%'",
			wantErr:   false,
		},
		{
			name:      "UnsupportedOperator",
			op:        model.FilterOperator(999),
			sensitive: false,
			want:      "",
			wantErr:   true,
		},
		{
			name:      "NullOperatorInsensitive",
			op:        model.Null,
			sensitive: false,
			want:      "",
			wantErr:   false,
		},
		{
			name:      "BetweenOperatorInsensitive",
			op:        model.Between,
			sensitive: false,
			want:      "? AND ?",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := buildFilterWildcards(tt.op, tt.sensitive)
			if (err != nil) != tt.wantErr {
				t.Errorf("buildFilterWildcards() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && got != tt.want {
				t.Errorf("buildFilterWildcards() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCastFilterValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		value         string
		valueType     model.FilterValueType
		op            model.FilterOperator
		want          any
		wantErr       bool
		caseSensitive bool
	}{
		{
			name:          "StringValueEquality",
			value:         "TeST",
			valueType:     model.String,
			op:            model.Eq,
			want:          "test",
			wantErr:       false,
			caseSensitive: true,
		},
		{
			name:      "IntegerValueInOperator",
			value:     "1,2,3",
			valueType: model.Integer,
			op:        model.In,
			want:      []int{1, 2, 3},
			wantErr:   false,
		},
		{
			name:      "FloatValueBetweenOperator",
			value:     "1.1,2.2",
			valueType: model.Float,
			op:        model.Between,
			want:      []float64{1.1, 2.2},
			wantErr:   false,
		},
		{
			name:      "BooleanValueNotInOperator",
			value:     "true,false",
			valueType: model.Boolean,
			op:        model.NotIn,
			want:      []bool{true, false},
			wantErr:   false,
		},
		{
			name:      "DateValueEquality",
			value:     "2023-01-01",
			valueType: model.Date,
			op:        model.Eq,
			want:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
			wantErr:   false,
		},
		{
			name:      "UnsupportedValueType",
			value:     "unsupported",
			valueType: model.FilterValueType(999),
			op:        model.Eq,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "UnsupportedOperatorForBoolean",
			value:     "true",
			valueType: model.Boolean,
			op:        model.Between,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "InvalidIntegerFormat",
			value:     "notAnInteger",
			valueType: model.Integer,
			op:        model.Eq,
			want:      0,
			wantErr:   true,
		},
		{
			name:      "InvalidFloatFormat",
			value:     "notAFloat",
			valueType: model.Float,
			op:        model.Eq,
			want:      0.0,
			wantErr:   true,
		},
		{
			name:      "InvalidBooleanFormat",
			value:     "notABoolean",
			valueType: model.Boolean,
			op:        model.Eq,
			want:      false,
			wantErr:   true,
		},
		{
			name:      "InvalidDateFormat",
			value:     "notADate",
			valueType: model.Date,
			op:        model.Eq,
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := castFilterValue(tt.value, tt.valueType, tt.op, tt.caseSensitive)
			if (err != nil) != tt.wantErr {
				t.Errorf("castFilterValue() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.want) {
				t.Errorf("castFilterValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCastStringValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		value         string
		op            model.FilterOperator
		want          any
		wantErr       bool
		caseSensitive bool
	}{
		{
			name:    "SimpleStringForEquality",
			value:   "test",
			op:      model.Eq,
			want:    "test",
			wantErr: false,
		},
		{
			name:          "StringWithCommaForInOperator",
			value:         "Test1,tESt2",
			op:            model.In,
			want:          []string{"test1", "test2"},
			wantErr:       false,
			caseSensitive: true,
		},
		{
			name:    "StringWithCommaForNotInOperator",
			value:   "test3,test4",
			op:      model.NotIn,
			want:    []string{"test3", "test4"},
			wantErr: false,
		},
		{
			name:    "StringWithCommaForBetweenOperator",
			value:   "start,end",
			op:      model.Between,
			want:    []string{"start", "end"},
			wantErr: false,
		},
		{
			name:    "EmptyString",
			value:   "",
			op:      model.Eq,
			want:    "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := castStringValue(tt.value, tt.op, tt.caseSensitive)
			if (err != nil) != tt.wantErr {
				t.Errorf("castStringValue() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("castTimeValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCastNumericValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     string
		op        model.FilterOperator
		want      any
		wantErr   bool
		valueType reflect.Kind
	}{
		{
			name:      "SingleValidFloat",
			value:     "123.456",
			op:        model.Eq,
			want:      123.456,
			wantErr:   false,
			valueType: reflect.Float64,
		},
		{
			name:      "InvalidFloatFormat",
			value:     "abc.def",
			op:        model.Eq,
			want:      0.0,
			wantErr:   true,
			valueType: reflect.Float64,
		},
		{
			name:      "InOperatorWithValidFloats",
			value:     "1.1,2.2,3.3",
			op:        model.In,
			want:      []float64{1.1, 2.2, 3.3},
			wantErr:   false,
			valueType: reflect.Float64,
		},
		{
			name:      "InOperatorWithOneInvalidFloat",
			value:     "4.4,notAFloat,5.5",
			op:        model.In,
			want:      nil,
			wantErr:   true,
			valueType: reflect.Float64,
		},
		{
			name:      "NotInOperatorWithValidFloats",
			value:     "6.6,7.7",
			op:        model.NotIn,
			want:      []float64{6.6, 7.7},
			wantErr:   false,
			valueType: reflect.Float64,
		},
		{
			name:      "BetweenOperatorWithValidFloats",
			value:     "8.8,9.9",
			op:        model.Between,
			want:      []float64{8.8, 9.9},
			wantErr:   false,
			valueType: reflect.Float64,
		},
		{
			name:      "BetweenOperatorWithInvalidFloats",
			value:     "10.10,invalid",
			op:        model.Between,
			want:      nil,
			wantErr:   true,
			valueType: reflect.Float64,
		},
		{
			name:      "SingleValidInteger",
			value:     "42",
			op:        model.Eq,
			want:      42,
			wantErr:   false,
			valueType: reflect.Int,
		},
		{
			name:      "InvalidIntegerFormat",
			value:     "notAnInteger",
			op:        model.Eq,
			want:      0,
			wantErr:   true,
			valueType: reflect.Int,
		},
		{
			name:      "InOperatorWithValidIntegers",
			value:     "1,2,3",
			op:        model.In,
			want:      []int{1, 2, 3},
			wantErr:   false,
			valueType: reflect.Int,
		},
		{
			name:      "InOperatorWithInvalidInteger",
			value:     "4,notAnInteger,6",
			op:        model.In,
			want:      nil,
			wantErr:   true,
			valueType: reflect.Int,
		},
		{
			name:      "NotInOperatorWithValidIntegers",
			value:     "7,8",
			op:        model.NotIn,
			want:      []int{7, 8},
			wantErr:   false,
			valueType: reflect.Int,
		},
		{
			name:      "BetweenOperatorWithValidIntegers",
			value:     "9,10",
			op:        model.Between,
			want:      []int{9, 10},
			wantErr:   false,
			valueType: reflect.Int,
		},
		{
			name:      "BetweenOperatorWithInvalidInteger",
			value:     "11,notAnInteger",
			op:        model.Between,
			want:      nil,
			wantErr:   true,
			valueType: reflect.Int,
		},
		{
			name:      "SingleTrueValue",
			value:     "true",
			op:        model.Eq,
			want:      true,
			wantErr:   false,
			valueType: reflect.Bool,
		},
		{
			name:      "SingleFalseValue",
			value:     "false",
			op:        model.Eq,
			want:      false,
			wantErr:   false,
			valueType: reflect.Bool,
		},
		{
			name:      "InvalidBooleanValue",
			value:     "not a boolean",
			op:        model.Eq,
			want:      false,
			wantErr:   true,
			valueType: reflect.Bool,
		},
		{
			name:      "InOperatorWithValidValues",
			value:     "true,false",
			op:        model.In,
			want:      []bool{true, false},
			wantErr:   false,
			valueType: reflect.Bool,
		},
		{
			name:      "InOperatorWithInvalidValue",
			value:     "true,not a boolean",
			op:        model.In,
			want:      nil,
			wantErr:   true,
			valueType: reflect.Bool,
		},
		{
			name:      "NotInOperatorWithValidValues",
			value:     "false,true",
			op:        model.NotIn,
			want:      []bool{false, true},
			wantErr:   false,
			valueType: reflect.Bool,
		},
		{
			name:      "BetweenOperatorUnsupported",
			value:     "true,false",
			op:        model.Between,
			want:      nil,
			wantErr:   true,
			valueType: reflect.Bool,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got any
			var err error
			var funcN string

			switch tt.valueType {
			case reflect.Float64:
				got, err = castFloatValue(tt.value, tt.op)
				funcN = "castFloatValue"
			case reflect.Int:
				got, err = castIntegerValue(tt.value, tt.op)
				funcN = "castIntegerValue"
			case reflect.Bool:
				got, err = castBooleanValue(tt.value, tt.op)
				funcN = "castBooleanValue"
			default:
				t.Errorf("unsupported value type: %v", tt.valueType)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("%s() error = %v, wantErr %v", funcN, err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("castTimeValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCastTimeValue(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		value     string
		valueType model.FilterValueType
		operator  model.FilterOperator
		want      any
		wantErr   bool
	}{
		{
			name:      "ValidDate",
			value:     "2023-01-01",
			valueType: model.Date,
			operator:  model.Eq,
			want:      time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
			wantErr:   false,
		},
		{
			name:      "ValidTime",
			value:     "15:04:05",
			valueType: model.Time,
			operator:  model.Eq,
			want:      time.Date(0, 1, 1, 15, 4, 5, 0, time.Local),
			wantErr:   false,
		},
		{
			name:      "ValidDateTime",
			value:     "2023-01-01 15:04:05",
			valueType: model.DateTime,
			operator:  model.Eq,
			want:      time.Date(2023, 1, 1, 15, 4, 5, 0, time.Local),
			wantErr:   false,
		},
		{
			name:      "InvalidDate",
			value:     "not a date",
			valueType: model.Date,
			operator:  model.Eq,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "BetweenDates",
			value:     "2023-01-01,2023-01-02",
			valueType: model.Date,
			operator:  model.Between,
			want: []time.Time{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
				time.Date(2023, 1, 2, 0, 0, 0, 0, time.Local),
			},
			wantErr: false,
		},
		{
			name:      "InvalidBetweenDates",
			value:     "2023-01-01,not a date",
			valueType: model.Date,
			operator:  model.Between,
			want:      nil,
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := castTimeValue(tt.value, tt.valueType, tt.operator)
			if (err != nil) != tt.wantErr {
				t.Errorf("castTimeValue() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("castTimeValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertSliceToNumeric(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		values    []string
		want      any
		wantErr   bool
		valueType reflect.Kind
	}{
		{
			name:      "AllValidIntegers",
			values:    []string{"1", "2", "3"},
			want:      []int{1, 2, 3},
			wantErr:   false,
			valueType: reflect.Int,
		},
		{
			name:      "ContainsInvalidInteger",
			values:    []string{"4", "not an int", "5"},
			want:      nil,
			wantErr:   true,
			valueType: reflect.Int,
		},
		{
			name:      "EmptySlice",
			values:    []string{},
			want:      []int{},
			wantErr:   false,
			valueType: reflect.Int,
		},
		{
			name:      "SingleValidInteger",
			values:    []string{"6"},
			want:      []int{6},
			wantErr:   false,
			valueType: reflect.Int,
		},
		{
			name:      "SingleInvalidInteger",
			values:    []string{"invalid"},
			want:      nil,
			wantErr:   true,
			valueType: reflect.Int,
		},
		{
			name:      "AllValidFloats",
			values:    []string{"1.1", "2.2", "3.3"},
			want:      []float64{1.1, 2.2, 3.3},
			wantErr:   false,
			valueType: reflect.Float64,
		},
		{
			name:      "ContainsInvalidFloat",
			values:    []string{"4.4", "not a float", "5.5"},
			want:      nil,
			wantErr:   true,
			valueType: reflect.Float64,
		},
		{
			name:      "EmptySlice",
			values:    []string{},
			want:      []float64{},
			wantErr:   false,
			valueType: reflect.Float64,
		},
		{
			name:      "SingleValidFloat",
			values:    []string{"6.6"},
			want:      []float64{6.6},
			wantErr:   false,
			valueType: reflect.Float64,
		},
		{
			name:      "SingleInvalidFloat",
			values:    []string{"invalid"},
			want:      nil,
			wantErr:   true,
			valueType: reflect.Float64,
		},
		{
			name:      "AllTrueVariants",
			values:    []string{"true", "True", "TRUE"},
			want:      []bool{true, true, true},
			wantErr:   false,
			valueType: reflect.Bool,
		},
		{
			name:      "AllFalseVariants",
			values:    []string{"false", "False", "FALSE"},
			want:      []bool{false, false, false},
			wantErr:   false,
			valueType: reflect.Bool,
		},
		{
			name:      "MixedTrueAndFalse",
			values:    []string{"true", "false", "True", "False", "TRUE", "FALSE"},
			want:      []bool{true, false, true, false, true, false},
			wantErr:   false,
			valueType: reflect.Bool,
		},
		{
			name:      "ContainsInvalidValue",
			values:    []string{"true", "not a bool", "false"},
			want:      nil,
			wantErr:   true,
			valueType: reflect.Bool,
		},
		{
			name:      "EmptySlice",
			values:    []string{},
			want:      []bool{},
			wantErr:   false,
			valueType: reflect.Bool,
		},
		{
			name:      "SingleInvalidValue",
			values:    []string{"invalid"},
			want:      nil,
			wantErr:   true,
			valueType: reflect.Bool,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var got any
			var err error
			var funcN string

			switch tt.valueType {
			case reflect.Int:
				got, err = convertSliceToInt(tt.values)
				funcN = "convertSliceToInt"
			case reflect.Float64:
				got, err = convertSliceToFloat(tt.values)
				funcN = "convertSliceToFloat"
			case reflect.Bool:
				got, err = convertSliceToBool(tt.values)
				funcN = "convertSliceToBool"
			default:
				t.Errorf("unsupported type: " + tt.valueType.String())
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("%s() error = %v, wantErr %v", funcN, err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("%s() got = %v, want %v", funcN, got, tt.want)
			}
		})
	}
}

func TestConvertSliceToDate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		values  []string
		layout  string
		want    []time.Time
		wantErr bool
	}{
		{
			name:   "SingleDate",
			values: []string{"2023-01-01"},
			layout: "2006-01-02",
			want:   []time.Time{time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local)},
		},
		{
			name:   "MultipleDates",
			values: []string{"2023-01-01", "2023-02-01"},
			layout: "2006-01-02",
			want: []time.Time{
				time.Date(2023, 1, 1, 0, 0, 0, 0, time.Local),
				time.Date(2023, 2, 1, 0, 0, 0, 0, time.Local),
			},
		},
		{
			name:    "InvalidDate",
			values:  []string{"2023-02-30"},
			layout:  "2006-01-02",
			wantErr: true,
		},
		{
			name:   "EmptySlice",
			values: []string{},
			layout: "2006-01-02",
			want:   []time.Time{},
		},
		{
			name:    "InvalidLayout",
			values:  []string{"01-203-02"},
			layout:  "02-206-01",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := convertSliceToDate(tt.values, tt.layout)
			if (err != nil) != tt.wantErr {
				t.Errorf("convertSliceToDate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.wantErr {
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("convertSliceToDate() got = %v, want %v", got, tt.want)
			}
		})
	}
}
