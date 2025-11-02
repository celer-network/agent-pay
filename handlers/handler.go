// Copyright 2018-2025 Celer Network

package handlers

import (
	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/ctype"
)

type ForwardToServerCallback func(dest ctype.Addr, retry bool, msg interface{}) (bool, error)

type CelerMsgHandler interface {
	GetMsgName() string
	CelerMsgRunnable
}

type CelerMsgRunnable interface {
	Run(msg *common.MsgFrame) error
}
