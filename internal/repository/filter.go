package repository

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/goutil/strutil"
	"github.com/pkg/errors"
	"github.com/vnworkday/account/internal/model"
)

const (
	keyValueSplitLen = 2
)

func AppendWhereClause(query *strings.Builder, filters ...string) {
	if len(filters) == 0 {
		return
	}

	query.WriteString(" WHERE ")
	query.WriteString(strings.Join(filters, " AND "))
}

func StringifyFilter(filter model.Filter, optAlias ...string) (string, error) {
	var alias string

	if len(optAlias) > 0 {
		alias = optAlias[0]
	}

	field, fieldErr := stringifyField(filter.Field, filter.IsCaseSensitive, alias)
	if fieldErr != nil {
		return "", errors.Wrap(fieldErr, "repository: failed to stringify filter field")
	}

	op, opErr := stringifyFilterOperator(filter.Operator)
	if opErr != nil {
		return "", errors.Wrap(opErr, "repository: failed to stringify filter operator")
	}

	wildcards, wildcardsErr := buildFilterWildcards(filter.Operator, filter.IsCaseSensitive)
	if wildcardsErr != nil {
		return "", errors.Wrap(wildcardsErr, "repository: failed to build filter wildcards")
	}

	ret := fmt.Sprintf("%s %s %s", field, op, wildcards)

	return strings.TrimSpace(ret), nil
}

func stringifyFilterOperator(operator model.FilterOperator) (string, error) {
	operators := map[model.FilterOperator]string{
		model.Eq:          "=",
		model.Ne:          "<>",
		model.Gt:          ">",
		model.Lt:          "<",
		model.Ge:          ">=",
		model.Le:          "<=",
		model.In:          "IN",
		model.NotIn:       "NOT IN",
		model.Contains:    "LIKE",
		model.NotContains: "NOT LIKE",
		model.StartsWith:  "LIKE",
		model.EndsWith:    "LIKE",
		model.Null:        "IS NULL",
		model.NotNull:     "IS NOT NULL",
		model.Between:     "BETWEEN",
	}

	op, exists := operators[operator]
	if !exists {
		return "", errors.Errorf("repository: unsupported filter operator: %d", operator)
	}

	return op, nil
}

func buildFilterWildcards(op model.FilterOperator, sensitive bool) (string, error) {
	wildcards := map[model.FilterOperator]string{
		model.Eq:          "?",
		model.Ne:          "?",
		model.Gt:          "?",
		model.Lt:          "?",
		model.Ge:          "?",
		model.Le:          "?",
		model.In:          "(?)",
		model.NotIn:       "(?)",
		model.Contains:    "'%' || ? || '%'",
		model.NotContains: "'%' || ? || '%'",
		model.StartsWith:  "? || '%'",
		model.EndsWith:    "'%' || ?",
		model.Null:        "",
		model.NotNull:     "",
		model.Between:     "? AND ?",
	}

	wildcard, exists := wildcards[op]
	if !exists {
		return "", errors.Errorf("repository: unsupported filter operator: %d", op)
	}

	if sensitive {
		switch op {
		case model.Eq, model.Ne:
			return "LOWER(?)", nil
		case model.Contains, model.NotContains, model.StartsWith, model.EndsWith:
			return strutil.Replaces(wildcard, map[string]string{"?": "LOWER(?)"}), nil
		default:
			return wildcard, nil
		}
	}

	return wildcard, nil
}

func castFilterValue(
	value string,
	valueType model.FilterValueType,
	op model.FilterOperator,
	caseSensitive bool,
) (any, error) {
	switch valueType {
	case model.String:
		return castStringValue(value, op, caseSensitive)
	case model.Integer:
		return castIntegerValue(value, op)
	case model.Float:
		return castFloatValue(value, op)
	case model.Boolean:
		return castBooleanValue(value, op)
	case model.Date, model.Time, model.DateTime:
		return castTimeValue(value, valueType, op)
	default:
		return nil, errors.Errorf("repository: unsupported filter value type: %d", valueType)
	}
}

func castStringValue(value string, op model.FilterOperator, caseSensitive bool) (any, error) {
	if caseSensitive {
		value = strings.ToLower(value)
	}

	if op == model.In || op == model.NotIn || op == model.Between {
		return strings.Split(value, ","), nil
	}

	return value, nil
}

func castIntegerValue(value string, op model.FilterOperator) (any, error) {
	if op == model.In || op == model.NotIn || op == model.Between {
		values := strings.Split(value, ",")

		return convertSliceToInt(values)
	}

	return strconv.Atoi(value)
}

func castFloatValue(value string, op model.FilterOperator) (any, error) {
	if op == model.In || op == model.NotIn || op == model.Between {
		values := strings.Split(value, ",")

		return convertSliceToFloat(values)
	}

	return strconv.ParseFloat(value, 64)
}

func castBooleanValue(value string, op model.FilterOperator) (any, error) {
	if op == model.Between {
		return nil, errors.Errorf("repository: unsupported filter operator for boolean: %d", op)
	} else if op == model.In || op == model.NotIn {
		values := strings.Split(value, ",")

		return convertSliceToBool(values)
	}

	return strconv.ParseBool(value)
}

func castTimeValue(value string, valueType model.FilterValueType, op model.FilterOperator) (any, error) {
	if op != model.Between {
		return strutil.ToTime(value, layoutForType(valueType))
	}

	values := strings.SplitN(value, ",", keyValueSplitLen)

	return convertSliceToDate(values, layoutForType(valueType))
}

func layoutForType(valueType model.FilterValueType) string {
	switch valueType {
	case model.Date:
		return "2006-01-02"
	case model.Time:
		return "15:04:05"
	case model.DateTime:
		return "2006-01-02 15:04:05"
	default:
		return ""
	}
}

func convertSliceToInt(values []string) ([]int, error) {
	intValues := make([]int, len(values))

	for idx, val := range values {
		if val == "" {
			continue
		}

		intValue, err := strconv.Atoi(val)
		if err != nil {
			return nil, errors.Wrapf(err, "repository: failed to cast filter value to integer: %s", val)
		}
		intValues[idx] = intValue
	}

	return intValues, nil
}

func convertSliceToFloat(values []string) ([]float64, error) {
	floatValues := make([]float64, len(values))

	for idx, val := range values {
		if val == "" {
			continue
		}

		floatValue, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil, errors.Wrapf(err, "repository: failed to cast filter value to float: %s", val)
		}
		floatValues[idx] = floatValue
	}

	return floatValues, nil
}

func convertSliceToBool(values []string) ([]bool, error) {
	boolValues := make([]bool, len(values))

	for idx, val := range values {
		if val == "" {
			continue
		}

		boolValue, err := strconv.ParseBool(val)
		if err != nil {
			return nil, errors.Wrapf(err, "repository: failed to cast filter value to boolean: %s", val)
		}
		boolValues[idx] = boolValue
	}

	return boolValues, nil
}

func convertSliceToDate(values []string, layout string) ([]time.Time, error) {
	dateValues := make([]time.Time, len(values))

	for idx, val := range values {
		if val == "" {
			continue
		}

		dateValue, err := strutil.ToTime(val, layout)
		if err != nil {
			return nil, errors.Wrapf(err, "repository: failed to cast filter value to date: %s", val)
		}
		dateValues[idx] = dateValue
	}

	return dateValues, nil
}
