package main

import (
	"reflect"
	"testing"

	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	"github.com/celer-network/agent-pay/rpc"
	"google.golang.org/protobuf/types/known/anypb"
)

type recordedPaymentCallbackEvent struct {
	kind   string
	reason rpc.PaymentSettleReason
	errMsg string
	payID  ctype.PayIDType
	pay    *entity.ConditionalPay
	note   *anypb.Any
}

type testPaymentCallbackSink struct {
	events []recordedPaymentCallbackEvent
}

func (s *testPaymentCallbackSink) HandleReceivingStart(payID ctype.PayIDType, pay *entity.ConditionalPay, note *anypb.Any) {
	s.events = append(s.events, recordedPaymentCallbackEvent{kind: "recv-start", payID: payID, pay: pay, note: note})
}

func (s *testPaymentCallbackSink) HandleReceivingDone(
	payID ctype.PayIDType,
	pay *entity.ConditionalPay,
	note *anypb.Any,
	reason rpc.PaymentSettleReason) {
	s.events = append(s.events, recordedPaymentCallbackEvent{kind: "recv-done", reason: reason, payID: payID, pay: pay, note: note})
}

func (s *testPaymentCallbackSink) HandleSendComplete(
	payID ctype.PayIDType,
	pay *entity.ConditionalPay,
	note *anypb.Any,
	reason rpc.PaymentSettleReason) {
	s.events = append(s.events, recordedPaymentCallbackEvent{kind: "send-complete", reason: reason, payID: payID, pay: pay, note: note})
}

func (s *testPaymentCallbackSink) HandleDestinationUnreachable(payID ctype.PayIDType, pay *entity.ConditionalPay, note *anypb.Any) {
	s.events = append(s.events, recordedPaymentCallbackEvent{kind: "dest-unreachable", payID: payID, pay: pay, note: note})
}

func (s *testPaymentCallbackSink) HandleSendFail(payID ctype.PayIDType, pay *entity.ConditionalPay, note *anypb.Any, errMsg string) {
	s.events = append(s.events, recordedPaymentCallbackEvent{kind: "send-fail", errMsg: errMsg, payID: payID, pay: pay, note: note})
}

func TestPaymentCallbacksMuxFansOutAllEvents(t *testing.T) {
	primary := &testPaymentCallbackSink{}
	secondary := &testPaymentCallbackSink{}
	mux := newPaymentCallbacksMux(primary, secondary)

	payID := ctype.ZeroPayID
	pay := &entity.ConditionalPay{Src: []byte{1}, Dest: []byte{2}}
	note := &anypb.Any{TypeUrl: "type.googleapis.com/test.Note", Value: []byte("note")}

	mux.HandleReceivingStart(payID, pay, note)
	mux.HandleReceivingDone(payID, pay, note, rpc.PaymentSettleReason_PAY_REJECTED)
	mux.HandleSendComplete(payID, pay, note, rpc.PaymentSettleReason_PAY_PAID_MAX)
	mux.HandleDestinationUnreachable(payID, pay, note)
	mux.HandleSendFail(payID, pay, note, "send failed")

	if len(primary.events) != 5 {
		t.Fatalf("primary sink saw %d events, want 5", len(primary.events))
	}
	if len(secondary.events) != 5 {
		t.Fatalf("secondary sink saw %d events, want 5", len(secondary.events))
	}
	if !reflect.DeepEqual(primary.events, secondary.events) {
		t.Fatalf("primary and secondary sinks diverged: primary=%#v secondary=%#v", primary.events, secondary.events)
	}

	gotKinds := make([]string, 0, len(primary.events))
	for _, event := range primary.events {
		gotKinds = append(gotKinds, event.kind)
	}
	wantKinds := []string{"recv-start", "recv-done", "send-complete", "dest-unreachable", "send-fail"}
	if !reflect.DeepEqual(gotKinds, wantKinds) {
		t.Fatalf("event kinds mismatch: got %v want %v", gotKinds, wantKinds)
	}
	if primary.events[1].reason != rpc.PaymentSettleReason_PAY_REJECTED {
		t.Fatalf("recv-done reason mismatch: got %v", primary.events[1].reason)
	}
	if primary.events[2].reason != rpc.PaymentSettleReason_PAY_PAID_MAX {
		t.Fatalf("send-complete reason mismatch: got %v", primary.events[2].reason)
	}
	if primary.events[4].errMsg != "send failed" {
		t.Fatalf("send-fail err mismatch: got %q", primary.events[4].errMsg)
	}
}