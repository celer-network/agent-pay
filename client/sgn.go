// Copyright 2018-2025 Celer Network

package client

import (
	"errors"

	"github.com/celer-network/agent-pay/cnode"
	"github.com/celer-network/agent-pay/entity"
	"github.com/celer-network/agent-pay/utils"
)

func (c *CelerClient) RequestSgnGuardState(token *entity.TokenInfo) error {
	cid, exist := c.getCidFromTokenInfo(token)
	if !exist {
		return errors.New("PSC_NOT_OPEN_" + utils.GetTokenAddrStr(token))
	}

	return c.cNode.RequestSgnGuardState(cid)
}

func (c *CelerClient) GetSgnSubscription() (*cnode.SgnSubscription, error) {
	return c.cNode.GetSgnSubscription()
}

func (c *CelerClient) GetSgnGuardRequest(token *entity.TokenInfo) (*cnode.SgnRequest, error) {
	cid, exist := c.getCidFromTokenInfo(token)
	if !exist {
		return nil, errors.New("PSC_NOT_OPEN_" + utils.GetTokenAddrStr(token))
	}

	return c.cNode.GetSgnGuardRequest(cid)
}
