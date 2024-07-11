package repository

import (
	"fmt"

	"github.com/vnworkday/account/internal/model"

	"github.com/pkg/errors"
)

func stringifyField(field string, sensitive bool, alias string) (string, error) {
	if field == "" {
		return "", errors.New("repository: sort field is required")
	}

	baseFormat := "%s"

	if sensitive {
		baseFormat = "LOWER(" + baseFormat + ")"
	}

	if alias != "" {
		field = fmt.Sprintf("%s.%s", alias, field)
	}

	return fmt.Sprintf(baseFormat, field), nil
}

func stringifyOp(operator model.Op) (string, error) {
	operators := map[model.Op]string{
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
