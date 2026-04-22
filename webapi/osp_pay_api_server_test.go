package webapi

import (
	"context"
	"strings"
	"testing"

	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/webapi/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type stubOspPayBackend struct {
	createRequest   *rpc.CreateAppSessionOnVirtualContractRequest
	createSessionID string
	createErr       error
	deleteSessionID string
	deleteErr       error
	statusSessionID string
	statusValue     uint8
	statusErr       error
}

func (*stubOspPayBackend) SendToken(*rpc.SendTokenRequest) (ctype.PayIDType, error) {
	return ctype.ZeroPayID, nil
}

func (*stubOspPayBackend) SendConditionalPayment(*rpc.SendConditionalPaymentRequest) (ctype.PayIDType, error) {
	return ctype.ZeroPayID, nil
}

func (b *stubOspPayBackend) CreateAppSessionOnVirtualContract(request *rpc.CreateAppSessionOnVirtualContractRequest) (string, error) {
	b.createRequest = request
	return b.createSessionID, b.createErr
}

func (b *stubOspPayBackend) DeleteAppSession(sessionID string) error {
	b.deleteSessionID = sessionID
	return b.deleteErr
}

func (b *stubOspPayBackend) GetStatusForAppSession(sessionID string) (uint8, error) {
	b.statusSessionID = sessionID
	return b.statusValue, b.statusErr
}

func (*stubOspPayBackend) GetIncomingPaymentState(ctype.PayIDType) (int, error) {
	return 0, nil
}

func (*stubOspPayBackend) GetIncomingPaymentRecord(ctype.PayIDType) (*PaymentRecord, error) {
	return nil, nil
}

func (*stubOspPayBackend) GetOutgoingPaymentState(ctype.PayIDType) (int, error) {
	return 0, nil
}

func (*stubOspPayBackend) ConfirmOutgoingPayment(ctype.PayIDType) error {
	return nil
}

func (*stubOspPayBackend) RejectIncomingPayment(ctype.PayIDType) error {
	return nil
}

func TestOspPayApiServerCreateAppSessionOnVirtualContract(t *testing.T) {
	backend := &stubOspPayBackend{createSessionID: "session-123"}
	server := NewOspPayApiServer(backend, nil)
	request := &rpc.CreateAppSessionOnVirtualContractRequest{ContractBin: "beef", ContractConstructor: "cafe", Nonce: 7, OnChainTimeout: 8}

	response, err := server.CreateAppSessionOnVirtualContract(context.Background(), request)
	if err != nil {
		t.Fatal(err)
	}
	if response.GetSessionId() != backend.createSessionID {
		t.Fatalf("CreateAppSessionOnVirtualContract session_id = %q, want %q", response.GetSessionId(), backend.createSessionID)
	}
	if backend.createRequest != request {
		t.Fatal("CreateAppSessionOnVirtualContract request not forwarded to backend")
	}
}

func TestOspPayApiServerGetStatusForAppSession(t *testing.T) {
	backend := &stubOspPayBackend{statusValue: 3}
	server := NewOspPayApiServer(backend, nil)

	response, err := server.GetStatusForAppSession(context.Background(), &rpc.SessionID{SessionId: "session-123"})
	if err != nil {
		t.Fatal(err)
	}
	if response.GetStatus() != uint32(backend.statusValue) {
		t.Fatalf("GetStatusForAppSession status = %d, want %d", response.GetStatus(), backend.statusValue)
	}
	if backend.statusSessionID != "session-123" {
		t.Fatalf("GetStatusForAppSession session_id = %q, want %q", backend.statusSessionID, "session-123")
	}
}

func TestOspPayApiServerDeleteAppSession(t *testing.T) {
	backend := &stubOspPayBackend{}
	server := NewOspPayApiServer(backend, nil)

	_, err := server.DeleteAppSession(context.Background(), &rpc.SessionID{SessionId: "session-123"})
	if err != nil {
		t.Fatal(err)
	}
	if backend.deleteSessionID != "session-123" {
		t.Fatalf("DeleteAppSession session_id = %q, want %q", backend.deleteSessionID, "session-123")
	}
}

func TestOspPayApiServerGetBalanceGuidance(t *testing.T) {
	server := NewOspPayApiServer(nil, nil)

	_, err := server.GetBalance(context.Background(), &rpc.TokenInfo{})
	st := status.Convert(err)
	if st.Code() != codes.Unimplemented {
		t.Fatalf("GetBalance code = %v, want %v", st.Code(), codes.Unimplemented)
	}
	if !strings.Contains(st.Message(), "CelerGetPeerStatus(peer, token)") {
		t.Fatalf("GetBalance message missing admin guidance: %q", st.Message())
	}
	if !strings.Contains(st.Message(), "ambiguous") {
		t.Fatalf("GetBalance message missing ambiguity guidance: %q", st.Message())
	}
}

func TestOspPayApiServerGetPeerFreeBalanceGuidance(t *testing.T) {
	server := NewOspPayApiServer(nil, nil)

	_, err := server.GetPeerFreeBalance(context.Background(), &rpc.GetPeerFreeBalanceRequest{})
	st := status.Convert(err)
	if st.Code() != codes.Unimplemented {
		t.Fatalf("GetPeerFreeBalance code = %v, want %v", st.Code(), codes.Unimplemented)
	}
	if !strings.Contains(st.Message(), "CelerGetPeerStatus(peer, token)") {
		t.Fatalf("GetPeerFreeBalance message missing admin guidance: %q", st.Message())
	}
}