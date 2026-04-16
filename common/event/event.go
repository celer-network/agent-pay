// Copyright 2018-2025 Celer Network

package event

import (
	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	"github.com/celer-network/agent-pay/rpc"
	"google.golang.org/protobuf/types/known/anypb"
)

const (
	CooperativeWithdraw = "CooperativeWithdraw"
	Deploy              = "Deploy"
	Deposit             = "Deposit"
	IntendSettle        = "IntendSettle"
	OpenChannel         = "OpenChannel"
	ConfirmSettle       = "ConfirmSettle"
	IntendWithdraw      = "IntendWithdraw"
	ConfirmWithdraw     = "ConfirmWithdraw"
	VetoWithdraw        = "VetoWithdraw"
	RouterUpdated       = "RouterUpdated"
	MigrateChannelTo    = "MigrateChannelTo"
)

type OpenChannelCallback interface {
	HandleOpenChannelFinish(cid ctype.CidType)
	HandleOpenChannelErr(e *common.E)
}
type OnNewStreamCallback interface {
	HandleNewCelerStream(addr ctype.Addr)
}
type OnReceivingTokenCallback interface {
	HandleReceivingStart(payID ctype.PayIDType, pay *entity.ConditionalPay, note *anypb.Any)
	HandleReceivingDone(
		payID ctype.PayIDType,
		pay *entity.ConditionalPay,
		note *anypb.Any,
		reason rpc.PaymentSettleReason)
}
type OnSendingTokenCallback interface {
	HandleSendComplete(
		payID ctype.PayIDType,
		pay *entity.ConditionalPay,
		note *anypb.Any,
		reason rpc.PaymentSettleReason)
	HandleDestinationUnreachable(payID ctype.PayIDType, pay *entity.ConditionalPay, note *anypb.Any)
	HandleSendFail(payID ctype.PayIDType, pay *entity.ConditionalPay, note *anypb.Any, errMsg string)
}
