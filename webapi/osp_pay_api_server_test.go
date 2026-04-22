package webapi

import (
	"context"
	"strings"
	"testing"

	"github.com/celer-network/agent-pay/webapi/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

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