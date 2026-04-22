//go:build !cgo

package storage

import (
	"errors"
	"testing"
)

func TestNewKVStoreSQLRejectsSQLiteWithoutCgo(t *testing.T) {
	_, err := NewKVStoreSQL(stDriverLT, tempStoreFile())
	if !errors.Is(err, ErrSQLiteRequiresCgo) {
		t.Fatalf("expected ErrSQLiteRequiresCgo, got %v", err)
	}
}