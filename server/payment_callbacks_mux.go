// Copyright 2018-2025 Celer Network

package main

import (
	"github.com/celer-network/agent-pay/common/event"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	"github.com/celer-network/agent-pay/rpc"
	"google.golang.org/protobuf/types/known/anypb"
)

type paymentCallbackSink interface {
	event.OnReceivingTokenCallback
	event.OnSendingTokenCallback
}

type paymentCallbacksMux struct {
	primary   paymentCallbackSink
	secondary paymentCallbackSink
}

type noopPaymentCallbackSink struct{}

func newPaymentCallbacksMux(primary, secondary paymentCallbackSink) *paymentCallbacksMux {
	if primary == nil {
		primary = noopPaymentCallbackSink{}
	}
	if secondary == nil {
		secondary = noopPaymentCallbackSink{}
	}
	return &paymentCallbacksMux{
		primary:   primary,
		secondary: secondary,
	}
}

func (noopPaymentCallbackSink) HandleReceivingStart(ctype.PayIDType, *entity.ConditionalPay, *anypb.Any) {}

func (noopPaymentCallbackSink) HandleReceivingDone(
	ctype.PayIDType,
	*entity.ConditionalPay,
	*anypb.Any,
	rpc.PaymentSettleReason) {
}

func (noopPaymentCallbackSink) HandleSendComplete(
	ctype.PayIDType,
	*entity.ConditionalPay,
	*anypb.Any,
	rpc.PaymentSettleReason) {
}

func (noopPaymentCallbackSink) HandleDestinationUnreachable(
	ctype.PayIDType,
	*entity.ConditionalPay,
	*anypb.Any) {
}

func (noopPaymentCallbackSink) HandleSendFail(
	ctype.PayIDType,
	*entity.ConditionalPay,
	*anypb.Any,
	string) {
}

func (m *paymentCallbacksMux) HandleReceivingStart(payID ctype.PayIDType, pay *entity.ConditionalPay, note *anypb.Any) {
	m.primary.HandleReceivingStart(payID, pay, note)
	m.secondary.HandleReceivingStart(payID, pay, note)
}

func (m *paymentCallbacksMux) HandleReceivingDone(
	payID ctype.PayIDType,
	pay *entity.ConditionalPay,
	note *anypb.Any,
	reason rpc.PaymentSettleReason) {
	m.primary.HandleReceivingDone(payID, pay, note, reason)
	m.secondary.HandleReceivingDone(payID, pay, note, reason)
}

func (m *paymentCallbacksMux) HandleSendComplete(
	payID ctype.PayIDType,
	pay *entity.ConditionalPay,
	note *anypb.Any,
	reason rpc.PaymentSettleReason) {
	m.primary.HandleSendComplete(payID, pay, note, reason)
	m.secondary.HandleSendComplete(payID, pay, note, reason)
}

func (m *paymentCallbacksMux) HandleDestinationUnreachable(
	payID ctype.PayIDType,
	pay *entity.ConditionalPay,
	note *anypb.Any) {
	m.primary.HandleDestinationUnreachable(payID, pay, note)
	m.secondary.HandleDestinationUnreachable(payID, pay, note)
}

func (m *paymentCallbacksMux) HandleSendFail(
	payID ctype.PayIDType,
	pay *entity.ConditionalPay,
	note *anypb.Any,
	errMsg string) {
	m.primary.HandleSendFail(payID, pay, note, errMsg)
	m.secondary.HandleSendFail(payID, pay, note, errMsg)
}