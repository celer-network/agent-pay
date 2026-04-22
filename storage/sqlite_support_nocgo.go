//go:build !cgo

package storage

func sqliteSupported() bool {
	return false
}