// Copyright 2018-2025 Celer Network

package msghdl

import (
	"fmt"

	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/config"
	"github.com/celer-network/agent-pay/utils"
)

func (h *CelerMsgHandler) HandleRoutingRequest(frame *common.MsgFrame) error {
	msg := frame.Message.GetRoutingRequest()
	var err error
	if h.routeController == nil {
		if config.EventListenerHttp == "" {
			return fmt.Errorf("both routeController and EventListenerHttp are empty")
		}
		err = utils.RecvRoutingInfo(config.EventListenerHttp, msg)
	} else {
		err = h.routeController.RecvBcastRoutingInfo(msg)
	}
	if err != nil {
		return fmt.Errorf("RecvBcastRoutingInfo err: %w", err)
	}
	return nil
}
