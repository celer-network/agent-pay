//go:build cgo

package storage

import (
	"os"
	"testing"
)

func TestNewKVStoreSQLAllowsSQLiteWithCgo(t *testing.T) {
	stFile := tempStoreFile()
	st, err := NewKVStoreSQL(stDriverLT, stFile)
	if err != nil {
		t.Fatalf("expected SQLite store with cgo enabled to succeed, got %v", err)
	}
	defer st.Close()
	defer os.Remove(stFile)
}