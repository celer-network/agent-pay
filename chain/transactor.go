// Copyright 2018-2026 Celer Network

package chain

import (
	"github.com/celer-network/goutils/eth"
	"github.com/ethereum/go-ethereum/core/types"
)

// SubmitWaitMined wraps `eth.TransactorPool.SubmitWaitMined` so any returned
// error carries the contract's custom-error selector via
// `WrapWithRevertSelector`. Use this from any site whose error may cross a
// transport boundary that flattens `rpc.DataError` to a string (e.g. gRPC
// to a remote client / SDK), so the receiver can still match against the
// canonical 4-byte selector.
//
// Direct `eth.TransactorPool.SubmitWaitMined` calls should stay only at
// internal call sites whose error stays in-process — they don't need the
// extra wrap. New cross-boundary call sites should default to this helper.
func SubmitWaitMined(
	pool *eth.TransactorPool,
	description string,
	method eth.TxMethod,
	opts ...eth.TxOption,
) (*types.Receipt, error) {
	receipt, err := pool.SubmitWaitMined(description, method, opts...)
	return receipt, WrapWithRevertSelector(err)
}

// TransactWaitMined wraps `eth.Transactor.TransactWaitMined` with the same
// selector preservation as `SubmitWaitMined`. Use it from cross-boundary
// call sites that operate on a single transactor (admin paths, OSP-only
// flows that still surface errors to clients).
func TransactWaitMined(
	transactor *eth.Transactor,
	description string,
	method eth.TxMethod,
	opts ...eth.TxOption,
) (*types.Receipt, error) {
	receipt, err := transactor.TransactWaitMined(description, method, opts...)
	return receipt, WrapWithRevertSelector(err)
}
