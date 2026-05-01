// Copyright 2018-2025 Celer Network

// Helpers for the surviving SimpleSingleSessionApp test fixture.
// This file (and singlesessionapp.go) is kept for back-compat with
// agent-pay-x402, which registers SimpleSingleSessionApp via
// CreateAppSessionOnVirtualContract. See
// docs/progress/app-session-simplification.md §7 for the migration plan
// that retires this surface.

package testapp

import (
	"math/big"
	"strings"

	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/goutils/log"
	"github.com/ethereum/go-ethereum/accounts/abi"
)

var (
	AppCode = ctype.Hex2Bytes(SimpleSingleSessionAppBin)
	Nonce   = big.NewInt(666)
	Timeout = big.NewInt(2)
)

// GetSingleSessionConstructor generates an abi-conforming constructor blob for
// SimpleSingleSessionApp. Used by agent-pay e2e tests and by agent-pay-x402's
// testinfra to register the app via CreateAppSessionOnVirtualContract.
func GetSingleSessionConstructor(players []ctype.Addr) []byte {
	parsedABI, err := abi.JSON(strings.NewReader(SimpleSingleSessionAppABI))
	if err != nil {
		log.Error(err)
		return nil
	}
	input, err := parsedABI.Pack("", players, Nonce, Timeout)
	if err != nil {
		log.Error(err)
		return nil
	}
	return input
}
