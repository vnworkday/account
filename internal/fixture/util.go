package fixture

import (
	"github.com/gookit/goutil/testutil/assert"
	"testing"
)

func ExpectationsWereMet[T any](t *testing.T, want, got T, wantErr bool, err error) {
	t.Helper()

	assert.DisableColor()
	assert.HideFullPath()

	if wantErr {
		assert.Error(t, err)
	} else {
		assert.NoError(t, err)
		assert.Equal(t, want, got)
	}
}
