// Copyright 2018-2025 Celer Network

package cnode

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"

	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/dispute"
	"github.com/levigross/grequests"
	"google.golang.org/protobuf/proto"
)

type SgnResponseWithHeight struct {
	Height int64           `json:"height"`
	Result json.RawMessage `json:"result"`
}

type SgnSubscription struct {
	EthAddress string  `json:"eth_address"`
	Deposit    big.Int `json:"deposit"`
	Spend      big.Int `json:"spend"`
}

type SgnRequest struct {
	ChannelId               []byte `json:"channel_id"`
	SeqNum                  uint64 `json:"seq_num"`
	SimplexSender           string `json:"simplex_sender"`
	SimplexReceiver         string `json:"simplex_receiver"`
	SignedSimplexStateBytes []byte `json:"signed_simplex_state_bytes"`
	DisputeTimeout          uint64 `json:"dispute_timeout"`
	TriggerTxHash           string `json:"trigger_tx_hash"`
	TriggerTxBlkNum         uint64 `json:"trigger_tx_blk_num"`
	GuardTxHash             string `json:"guard_tx_hash"`
	GuardTxBlkNum           uint64 `json:"guard_tx_blk_num"`
	GuardSender             string `json:"guard_sender"`
}

func (c *CNode) RequestSgnGuardState(cid ctype.CidType) error {
	_, signedPeerSimplexState, found, err := c.dal.GetPeerSimplex(cid)
	if err != nil {
		return err
	}
	if !found {
		return common.ErrChannelNotFound
	}

	sigSortedPeerSimplexState, err := dispute.SigSortedSimplexState(signedPeerSimplexState)
	if err != nil {
		return err
	}

	sigSortedPeerSimplexStateBytes, err := proto.Marshal(sigSortedPeerSimplexState)
	if err != nil {
		return err
	}

	myPeerSimplexSig := c.SignState(sigSortedPeerSimplexStateBytes)

	resp, err := grequests.Post(
		c.sgnGw+"/guard/requestGuard",
		&grequests.RequestOptions{
			JSON: map[string]string{
				"signed_simplex_state_bytes": ctype.Bytes2Hex(sigSortedPeerSimplexStateBytes),
				"simplex_receiver_sig":       ctype.Bytes2Hex(myPeerSimplexSig),
			},
		},
	)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("requestGuard failed: %s", resp)
	}

	return nil
}

func (c *CNode) GetSgnSubscription() (*SgnSubscription, error) {
	resp, err := grequests.Get(
		c.sgnGw+"/guard/subscription/"+c.EthAddress.String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get subscription err %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get subscription status %d, err %s", resp.StatusCode, resp.String())
	}

	var responseWithHeight SgnResponseWithHeight
	err = json.Unmarshal(resp.Bytes(), &responseWithHeight)
	if err != nil {
		return nil, fmt.Errorf("Parse subscription response err: %w", err)
	}
	var subscription SgnSubscription
	err = json.Unmarshal(responseWithHeight.Result, &subscription)
	if err != nil {
		return nil, fmt.Errorf("Parse subscription err: %w", err)
	}

	return &subscription, nil
}

func (c *CNode) GetSgnGuardRequest(cid ctype.CidType) (*SgnRequest, error) {
	resp, err := grequests.Get(
		c.sgnGw+"/guard/request/"+ctype.Cid2Hex(cid)+"/"+c.EthAddress.String(),
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("get guard request err %w", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get subscription status %d, err %s", resp.StatusCode, resp.String())
	}

	var responseWithHeight SgnResponseWithHeight
	err = json.Unmarshal(resp.Bytes(), &responseWithHeight)
	if err != nil {
		return nil, fmt.Errorf("Parse guard request response err: %w", err)
	}
	var request SgnRequest
	err = json.Unmarshal(responseWithHeight.Result, &request)
	if err != nil {
		return nil, fmt.Errorf("Parse subscription err: %w", err)
	}

	return &request, nil
}
