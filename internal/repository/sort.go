package repository

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/vnworkday/account/internal/model"
)

func StringifySort(sort model.Sort, optAlias ...string) (string, error) {
	var alias string

	if len(optAlias) > 0 {
		alias = optAlias[0]
	}

	field, fieldErr := stringifyField(sort.Field, sort.IsCaseSensitive, alias)
	if fieldErr != nil {
		return "", errors.Wrap(fieldErr, "repository: failed to stringify sort field")
	}

	var order string

	switch sort.Order {
	case model.Asc:
		order = "ASC"
	case model.Desc:
		order = "DESC"
	default:
		return "", errors.New("repository: invalid sort order" + string(sort.Order))
	}

	return fmt.Sprintf("%s %s", field, order), nil
}
