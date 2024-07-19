package repo

import (
	"fmt"

	"github.com/vnworkday/account/internal/common/domain"

	"github.com/pkg/errors"
)

func StringifySort(sort domain.Sort, optAlias ...string) (string, error) {
	var alias string

	if len(optAlias) > 0 {
		alias = optAlias[0]
	}

	field, fieldErr := stringifyField(sort.Field, sort.CaseSensitive, alias)
	if fieldErr != nil {
		return "", errors.Wrap(fieldErr, "repository: failed to stringify sort field")
	}

	var order string

	switch sort.Order {
	case domain.Asc:
		order = "ASC"
	case domain.Desc:
		order = "DESC"
	default:
		return "", errors.New("repository: invalid sort order" + string(sort.Order))
	}

	return fmt.Sprintf("%s %s", field, order), nil
}
