// Copyright 2018-2025 Celer Network

package webapi

import (
	"sync"

	"github.com/celer-network/agent-pay/celersdkintf"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	"github.com/celer-network/agent-pay/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/anypb"
)

const paymentEventFeedBufferSize = 16

type outgoingPaymentEvent struct {
	payment *celersdkintf.Payment
	err     *celersdkintf.E
}

// PaymentEventFeed keeps at most one active subscriber per direction.
// Publish is best-effort and non-blocking, so slow subscribers may drop events.
type PaymentEventFeed struct {
	incomingMu sync.Mutex
	incoming   chan *celersdkintf.Payment
	outgoingMu sync.Mutex
	outgoing   chan *outgoingPaymentEvent
}

func NewPaymentEventFeed() *PaymentEventFeed {
	return &PaymentEventFeed{}
}

func (f *PaymentEventFeed) SubscribeIncoming() (<-chan *celersdkintf.Payment, func(), error) {
	f.incomingMu.Lock()
	defer f.incomingMu.Unlock()
	if f.incoming != nil {
		return nil, nil, status.Error(codes.FailedPrecondition, "incoming payment subscription already active")
	}
	ch := make(chan *celersdkintf.Payment, paymentEventFeedBufferSize)
	f.incoming = ch
	return ch, func() {
		f.incomingMu.Lock()
		defer f.incomingMu.Unlock()
		if f.incoming == ch {
			close(ch)
			f.incoming = nil
		}
	}, nil
}

func (f *PaymentEventFeed) SubscribeOutgoing() (<-chan *outgoingPaymentEvent, func(), error) {
	f.outgoingMu.Lock()
	defer f.outgoingMu.Unlock()
	if f.outgoing != nil {
		return nil, nil, status.Error(codes.FailedPrecondition, "outgoing payment subscription already active")
	}
	ch := make(chan *outgoingPaymentEvent, paymentEventFeedBufferSize)
	f.outgoing = ch
	return ch, func() {
		f.outgoingMu.Lock()
		defer f.outgoingMu.Unlock()
		if f.outgoing == ch {
			close(ch)
			f.outgoing = nil
		}
	}, nil
}

func (f *PaymentEventFeed) HandleReceivingStart(payID ctype.PayIDType, pay *entity.ConditionalPay, note *anypb.Any) {
	f.publishIncoming(paymentFromCondPay(payID, pay, note, celersdkintf.PAY_STATUS_PENDING))
}

func (f *PaymentEventFeed) HandleReceivingDone(
	payID ctype.PayIDType,
	pay *entity.ConditionalPay,
	note *anypb.Any,
	reason rpc.PaymentSettleReason) {
	f.publishIncoming(paymentFromCondPay(payID, pay, note, settleReasonToSdkStatus(reason)))
}

func (f *PaymentEventFeed) HandleSendComplete(
	payID ctype.PayIDType,
	pay *entity.ConditionalPay,
	note *anypb.Any,
	reason rpc.PaymentSettleReason) {
	f.publishOutgoing(&outgoingPaymentEvent{
		payment: paymentFromCondPay(payID, pay, note, settleReasonToSdkStatus(reason)),
	})
}

func (f *PaymentEventFeed) HandleDestinationUnreachable(payID ctype.PayIDType, pay *entity.ConditionalPay, note *anypb.Any) {
	f.publishOutgoing(&outgoingPaymentEvent{
		payment: paymentFromCondPay(payID, pay, note, celersdkintf.PAY_STATUS_UNPAID_DEST_UNREACHABLE),
		err: &celersdkintf.E{
			Reason: "Unreachable to " + ctype.Bytes2Hex(pay.GetDest()),
			Code:   -1,
		},
	})
}

func (f *PaymentEventFeed) HandleSendFail(payID ctype.PayIDType, pay *entity.ConditionalPay, note *anypb.Any, errMsg string) {
	f.publishOutgoing(&outgoingPaymentEvent{
		payment: paymentFromCondPay(payID, pay, note, celersdkintf.PAY_STATUS_UNPAID),
		err: &celersdkintf.E{
			Reason: errMsg + ": Send token failed to " + ctype.Bytes2Hex(pay.GetDest()),
			Code:   -1,
		},
	})
}

func (f *PaymentEventFeed) publishIncoming(payment *celersdkintf.Payment) {
	f.incomingMu.Lock()
	defer f.incomingMu.Unlock()
	if f.incoming == nil {
		return
	}
	select {
	case f.incoming <- payment:
	default:
	}
}

func (f *PaymentEventFeed) publishOutgoing(event *outgoingPaymentEvent) {
	f.outgoingMu.Lock()
	defer f.outgoingMu.Unlock()
	if f.outgoing == nil {
		return
	}
	select {
	case f.outgoing <- event:
	default:
	}
}