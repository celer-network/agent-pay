// Copyright 2018-2026 Celer Network

package chain

import (
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	ethrpc "github.com/ethereum/go-ethereum/rpc"
	"golang.org/x/crypto/sha3"
)

// Solidity custom errors revert with a 4-byte selector + ABI-encoded args.
// Geth's eth_estimateGas / eth_call propagate the revert payload via the
// JSON-RPC error's `data` field, which go-ethereum exposes through the
// `rpc.DataError` interface. The legacy require-string path put the reason
// in the error message; custom errors do not — the message is just
// "execution reverted", and the precise reason lives in the data field.
//
// Helpers below extract the selector from an arbitrary error chain and
// embed it in error messages crossing layers (gRPC, OSP-to-client) so the
// custom-error name is preserved as a hex selector that the receiver can
// match against `ErrorSelector(sig)`.
//
// The hex format used in error messages is: `revert selector: 0xXXXXXXXX`.

// errorSelectorPrefix is the substring marker we use when embedding a
// selector into a wrapping error message. Receivers grep for this prefix
// when they want to recover the selector across a transport boundary that
// flattens error to string (e.g. gRPC).
const errorSelectorPrefix = "revert selector: 0x"

// ErrorSelector returns the 4-byte selector for a Solidity custom-error
// signature, e.g. "ConditionNotFinalized()" or "InvalidSeqNum(uint256)".
// The selector is the first 4 bytes of keccak256(sig).
func ErrorSelector(sig string) [4]byte {
	h := sha3.NewLegacyKeccak256()
	h.Write([]byte(sig))
	var sel [4]byte
	copy(sel[:], h.Sum(nil)[:4])
	return sel
}

// ErrorSelectorHex is ErrorSelector formatted as `0x` + 8 hex chars.
func ErrorSelectorHex(sig string) string {
	sel := ErrorSelector(sig)
	return "0x" + hex.EncodeToString(sel[:])
}

// ParseRevertSelector walks the error chain looking for an `rpc.DataError`,
// pulls its `ErrorData()` (the hex-encoded revert payload), and returns the
// leading 4-byte selector. Falls back to scanning the error message for an
// embedded `revert selector: 0x...` token left by `WrapWithRevertSelector`
// — used when the error has crossed a transport (e.g. gRPC) that flattened
// the original `rpc.DataError`.
func ParseRevertSelector(err error) (selector [4]byte, ok bool) {
	if err == nil {
		return selector, false
	}
	// Try the direct rpc.DataError path first (in-process call chain).
	var dataErr ethrpc.DataError
	if errors.As(err, &dataErr) {
		data := dataErr.ErrorData()
		if hexStr, isStr := data.(string); isStr {
			if sel, sok := decodeSelectorHex(hexStr); sok {
				return sel, true
			}
		}
	}
	// Fall back to the embedded-token path (post-transport).
	return parseEmbeddedSelector(err.Error())
}

// WrapWithRevertSelector wraps `err` with a "revert selector: 0xXXXXXXXX"
// suffix when the error chain carries a custom-error revert payload. Use
// this on the side that holds the original `rpc.DataError` (typically the
// OSP) before returning across a transport boundary that flattens the
// error to a string. If no selector is present in `err`, or if the
// message already carries the selector token (i.e. an upstream layer
// already wrapped), the original error is returned unchanged.
func WrapWithRevertSelector(err error) error {
	if err == nil {
		return nil
	}
	// Idempotent: if an upstream layer already embedded the token, don't
	// duplicate it. This guards against double-wrap noise when the same
	// error walks through several wrap-aware layers (e.g. helper at the
	// transactor + helper at the dispute path) before crossing transport.
	if strings.Contains(err.Error(), errorSelectorPrefix) {
		return err
	}
	sel, ok := ParseRevertSelector(err)
	if !ok {
		return err
	}
	return fmt.Errorf("%w (%s%s)", err, errorSelectorPrefix, hex.EncodeToString(sel[:]))
}

func decodeSelectorHex(hexStr string) ([4]byte, bool) {
	var sel [4]byte
	hexStr = strings.TrimPrefix(hexStr, "0x")
	if len(hexStr) < 8 {
		return sel, false
	}
	raw, err := hex.DecodeString(hexStr[:8])
	if err != nil {
		return sel, false
	}
	copy(sel[:], raw)
	return sel, true
}

func parseEmbeddedSelector(msg string) ([4]byte, bool) {
	idx := strings.Index(msg, errorSelectorPrefix)
	if idx < 0 {
		return [4]byte{}, false
	}
	return decodeSelectorHex(msg[idx+len(errorSelectorPrefix):])
}
