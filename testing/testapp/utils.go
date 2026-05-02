// Copyright 2018-2025 Celer Network

// Helpers for the SimpleSingleSessionApp test fixture. New code should
// prefer BooleanCondMock (in this same package) — a stateless IBooleanCond
// implementation that doesn't carry the legacy session-state-machine surface.

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
	timeout = big.NewInt(2)
)

// GetSingleSessionConstructor generates an abi-conforming constructor blob
// for SimpleSingleSessionApp.
func GetSingleSessionConstructor(players []ctype.Addr) []byte {
	parsedABI, err := abi.JSON(strings.NewReader(SimpleSingleSessionAppABI))
	if err != nil {
		log.Error(err)
		return nil
	}
	input, err := parsedABI.Pack("", players, Nonce, timeout)
	if err != nil {
		log.Error(err)
		return nil
	}
	return input
}
