// Copyright 2018-2026 Celer Network

package chain

import (
	"encoding/hex"
	"errors"
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"
)

// dataErrStub mimics the DataError interface that go-ethereum's RPC client
// returns from eth_estimateGas / eth_call when the contract reverts with a
// custom error. The hex string encoding ("0x" + selector + abi-encoded
// args) is what geth puts in the JSON-RPC `data` field; we mirror that.
type dataErrStub struct {
	msg  string
	data string
}

func (e *dataErrStub) Error() string         { return e.msg }
func (e *dataErrStub) ErrorCode() int        { return 3 }
func (e *dataErrStub) ErrorData() interface{} { return e.data }

func TestErrorSelector_MatchesEthereumKeccak(t *testing.T) {
	// Cross-validate the helper against go-ethereum's canonical Keccak256
	// implementation. Solidity's custom-error selector is the first 4
	// bytes of keccak256(signature) — pinning the hashing primitive
	// guarantees the selectors we compute match what the contract emits
	// on revert. (If the helper ever drifted from the EVM's keccak, this
	// test would catch it before the dispute e2e does.)
	sigs := []string{
		"ConditionNotFinalized()",
		"NotOperator()",
		"NotWalletOwner()",
		"ZeroAddress()",
		"BalanceLimitExceeded(uint256,uint256)",
		"SeqNumOutOfOrder(uint256,uint256)",
	}
	for _, sig := range sigs {
		want := "0x" + hex.EncodeToString(crypto.Keccak256([]byte(sig))[:4])
		got := ErrorSelectorHex(sig)
		if got != want {
			t.Errorf("ErrorSelectorHex(%q) = %q, want %q", sig, got, want)
		}
	}
}

func TestParseRevertSelector_FromDataError(t *testing.T) {
	want := ErrorSelector("ConditionNotFinalized()")
	stub := &dataErrStub{msg: "execution reverted", data: ErrorSelectorHex("ConditionNotFinalized()")}

	got, ok := ParseRevertSelector(stub)
	if !ok {
		t.Fatalf("ParseRevertSelector returned ok=false on a DataError")
	}
	if got != want {
		t.Errorf("got selector %x, want %x", got, want)
	}
}

func TestParseRevertSelector_FromDataErrorWithArgs(t *testing.T) {
	// Solidity custom errors with arguments encode as
	// "0x<4-byte-selector><abi-encoded-args>". The parser must take only
	// the leading 4 bytes regardless of any args that follow.
	sigSel := ErrorSelector("BalanceLimitExceeded(uint256,uint256)")
	dataHex := ErrorSelectorHex("BalanceLimitExceeded(uint256,uint256)") +
		"0000000000000000000000000000000000000000000000000000000000000064" +
		"0000000000000000000000000000000000000000000000000000000000000032"
	stub := &dataErrStub{msg: "execution reverted", data: dataHex}

	got, ok := ParseRevertSelector(stub)
	if !ok || got != sigSel {
		t.Errorf("got (%x, %v), want (%x, true)", got, ok, sigSel)
	}
}

func TestWrapWithRevertSelector_RoundTripsAcrossStringFlatten(t *testing.T) {
	// Simulates the full cross-boundary path:
	//   contract revert → rpc.DataError surfaces at OSP →
	//   chain.WrapWithRevertSelector embeds the hex selector in the
	//   error message → gRPC layer flattens to a string error →
	//   test side calls ParseRevertSelector on the flat string and
	//   recovers the same selector.
	want := ErrorSelector("ConditionNotFinalized()")
	stub := &dataErrStub{msg: "execution reverted", data: ErrorSelectorHex("ConditionNotFinalized()")}
	wrapped := WrapWithRevertSelector(stub)
	if wrapped == nil {
		t.Fatal("WrapWithRevertSelector returned nil on non-nil input")
	}

	// Flatten through string — what gRPC's status.Error(...) does to the
	// underlying typed error.
	flat := errors.New(wrapped.Error())
	got, ok := ParseRevertSelector(flat)
	if !ok {
		t.Fatalf("ParseRevertSelector failed to recover selector from flattened err: %v", flat)
	}
	if got != want {
		t.Errorf("got selector %x, want %x (flat err: %s)", got, want, flat)
	}
}

func TestWrapWithRevertSelector_NoOpOnPlainError(t *testing.T) {
	// Errors that don't carry a custom-error revert payload must round-trip
	// unchanged so callers can't false-positive on incidental hex bytes.
	plain := errors.New("dial tcp: connection refused")
	got := WrapWithRevertSelector(plain)
	if got != plain {
		t.Errorf("expected plain error to round-trip unchanged; got %v", got)
	}
	if _, ok := ParseRevertSelector(plain); ok {
		t.Errorf("ParseRevertSelector unexpectedly succeeded on a plain error")
	}
}

func TestWrapWithRevertSelector_NilSafe(t *testing.T) {
	if got := WrapWithRevertSelector(nil); got != nil {
		t.Errorf("WrapWithRevertSelector(nil) = %v, want nil", got)
	}
	if _, ok := ParseRevertSelector(nil); ok {
		t.Errorf("ParseRevertSelector(nil) returned ok=true")
	}
}

func TestParseRevertSelector_WrappedWithFmtErrorf(t *testing.T) {
	// goutils transactor wraps errors with `fmt.Errorf("tx dry-run err: %w, ...", err)`.
	// errors.As must still find the DataError through that chain.
	stub := &dataErrStub{msg: "execution reverted", data: ErrorSelectorHex("NotOperator()")}
	wrapped := fmt.Errorf("tx dry-run err: %w, calldata: 0xabcd", stub)
	got, ok := ParseRevertSelector(wrapped)
	if !ok {
		t.Fatalf("ParseRevertSelector failed through fmt.Errorf wrap")
	}
	if got != ErrorSelector("NotOperator()") {
		t.Errorf("got %x, want %x", got, ErrorSelector("NotOperator()"))
	}
}
