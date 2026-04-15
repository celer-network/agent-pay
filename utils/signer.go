// Copyright 2018-2025 Celer Network

package utils

import (
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/goutils/eth"
)

type onChainCompatibleSigner struct {
	eth.Signer
}

func NewOnChainCompatibleSigner(signer eth.Signer) eth.Signer {
	if signer == nil {
		return nil
	}
	return &onChainCompatibleSigner{Signer: signer}
}

func (s *onChainCompatibleSigner) SignEthMessage(data []byte) ([]byte, error) {
	sig, err := s.Signer.SignEthMessage(data)
	if err != nil {
		return nil, err
	}
	return ctype.ToOnChainSig(sig), nil
}

func (s *onChainCompatibleSigner) SignEthTransaction(rawTx []byte) ([]byte, error) {
	return s.Signer.SignEthTransaction(rawTx)
}