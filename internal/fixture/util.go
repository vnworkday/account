package fixture

import (
	"github.com/gookit/goutil/reflects"
	"github.com/pkg/errors"
)

func ExpectationsWereMet[T any](want, got T, wantErr bool, err error) error {
	if !wantErr && err != nil {
		return errors.Errorf("\nwant: %v\n got: %v", "no error", err)
	}

	if !wantErr && !reflects.IsEqual(want, got) {
		return errors.Errorf("\nwant: %v\n got: %v", want, got)
	}

	if wantErr && err == nil {
		return errors.Errorf("\nwant: %v\n got: %v", "error occurred", "no error")
	}

	return nil
}
