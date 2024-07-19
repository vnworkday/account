package repo

import (
	"fmt"

	"github.com/vnworkday/account/internal/common/domain"

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

func stringifyOp(operator domain.Op) (string, error) {
	operators := map[domain.Op]string{
		domain.Eq:          "=",
		domain.Ne:          "<>",
		domain.Gt:          ">",
		domain.Lt:          "<",
		domain.Ge:          ">=",
		domain.Le:          "<=",
		domain.In:          "IN",
		domain.NotIn:       "NOT IN",
		domain.Contains:    "LIKE",
		domain.NotContains: "NOT LIKE",
		domain.StartsWith:  "LIKE",
		domain.EndsWith:    "LIKE",
		domain.Null:        "IS NULL",
		domain.NotNull:     "IS NOT NULL",
		domain.Between:     "BETWEEN",
	}

	op, exists := operators[operator]
	if !exists {
		return "", errors.Errorf("repository: unsupported filter operator: %d", operator)
	}

	return op, nil
}
