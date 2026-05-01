// Copyright 2018-2025 Celer Network

package webapi

import (
	"context"
	"testing"

	"github.com/celer-network/agent-pay/celersdk"
	"github.com/celer-network/agent-pay/webapi/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestGetAppSessionUnknownID(t *testing.T) {
	s := &ApiServer{
		appSessionMap: make(map[string]*celersdk.AppSession),
	}
	_, err := s.getAppSession("never-created")
	if status.Code(err) != codes.NotFound {
		t.Fatalf("getAppSession(unknown) error code = %v, want %v (err=%v)",
			status.Code(err), codes.NotFound, err)
	}
}

func TestGetAppSessionAfterDelete(t *testing.T) {
	s := &ApiServer{
		appSessionMap: make(map[string]*celersdk.AppSession),
	}
	// Inject a non-nil session-like sentinel so the "exists then deleted" path
	// uses the same map shape DeleteAppSession produces.
	s.appSessionMap["sid"] = &celersdk.AppSession{ID: "sid"}
	got, err := s.getAppSession("sid")
	if err != nil || got == nil {
		t.Fatalf("getAppSession(sid) before delete returned (%v, %v); want non-nil session", got, err)
	}

	// Mirror what DeleteAppSession does to the map (the apiClient.EndAppSession
	// half is bypassed here because the test stub doesn't have one).
	delete(s.appSessionMap, "sid")

	_, err = s.getAppSession("sid")
	if status.Code(err) != codes.NotFound {
		t.Fatalf("getAppSession(sid) after delete error code = %v, want %v (err=%v)",
			status.Code(err), codes.NotFound, err)
	}
}

// Public-handler regression: a query for an unknown session ID must return
// codes.NotFound at the gRPC surface (rather than nil-deref panicking, the
// pre-fix behavior). Exercises the actual surviving handlers, not just the
// helper they delegate to.
func TestGetDeployedAddressForAppSessionUnknownID(t *testing.T) {
	s := &ApiServer{
		appSessionMap: make(map[string]*celersdk.AppSession),
	}
	_, err := s.GetDeployedAddressForAppSession(
		context.Background(),
		&rpc.SessionID{SessionId: "never-created"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("GetDeployedAddressForAppSession(unknown) error code = %v, want %v (err=%v)",
			status.Code(err), codes.NotFound, err)
	}
}

func TestGetBooleanOutcomeForAppSessionUnknownID(t *testing.T) {
	s := &ApiServer{
		appSessionMap: make(map[string]*celersdk.AppSession),
	}
	_, err := s.GetBooleanOutcomeForAppSession(
		context.Background(),
		&rpc.GetBooleanOutcomeForAppSessionRequest{SessionId: "never-created"})
	if status.Code(err) != codes.NotFound {
		t.Fatalf("GetBooleanOutcomeForAppSession(unknown) error code = %v, want %v (err=%v)",
			status.Code(err), codes.NotFound, err)
	}
}
