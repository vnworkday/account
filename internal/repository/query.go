package repository

import (
	"fmt"

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
