package main

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/celer-network/agent-pay/cnode/cooperativewithdraw"
	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAdminServiceCooperativeWithdrawUsesExplicitAmount(t *testing.T) {
	var gotCid ctype.CidType
	var gotAmount *big.Int
	var removed string

	svc := &adminService{
		cooperativeWithdrawWithCallback: func(cid ctype.CidType, amount *big.Int, cb cooperativewithdraw.Callback) (string, error) {
			gotCid = cid
			gotAmount = new(big.Int).Set(amount)
			cb.OnWithdraw("withdraw-hash", "0xtx")
			return "withdraw-hash", nil
		},
		getChannelBalance: func(cid ctype.CidType) (*common.ChannelBalance, error) {
			t.Fatalf("getChannelBalance should not be called for explicit amount")
			return nil, nil
		},
		removeCooperativeWithdrawJob: func(withdrawHash string) {
			removed = withdrawHash
		},
	}

	resp, err := svc.CooperativeWithdraw(context.Background(), &rpc.ChannelOpRequest{
		Cid: "0x1",
		Wei: "123",
	})
	if err != nil {
		t.Fatalf("CooperativeWithdraw() error = %v", err)
	}
	if resp.GetStatus() != 0 {
		t.Fatalf("CooperativeWithdraw() status = %d, want 0", resp.GetStatus())
	}
	if gotCid != ctype.Hex2Cid("1") {
		t.Fatalf("CooperativeWithdraw() cid = %s, want %s", ctype.Cid2Hex(gotCid), ctype.Cid2Hex(ctype.Hex2Cid("1")))
	}
	if gotAmount == nil || gotAmount.Cmp(big.NewInt(123)) != 0 {
		t.Fatalf("CooperativeWithdraw() amount = %v, want 123", gotAmount)
	}
	if removed != "withdraw-hash" {
		t.Fatalf("RemoveCooperativeWithdrawJob() hash = %q, want withdraw-hash", removed)
	}
}

func TestAdminServiceCooperativeWithdrawUsesFreeBalanceWhenWeiMissing(t *testing.T) {
	var gotAmount *big.Int

	svc := &adminService{
		cooperativeWithdrawWithCallback: func(cid ctype.CidType, amount *big.Int, cb cooperativewithdraw.Callback) (string, error) {
			gotAmount = new(big.Int).Set(amount)
			cb.OnWithdraw("withdraw-hash", "0xtx")
			return "withdraw-hash", nil
		},
		getChannelBalance: func(cid ctype.CidType) (*common.ChannelBalance, error) {
			return &common.ChannelBalance{MyFree: big.NewInt(456)}, nil
		},
		removeCooperativeWithdrawJob: func(string) {},
	}

	resp, err := svc.CooperativeWithdraw(context.Background(), &rpc.ChannelOpRequest{Cid: "2"})
	if err != nil {
		t.Fatalf("CooperativeWithdraw() error = %v", err)
	}
	if resp.GetStatus() != 0 {
		t.Fatalf("CooperativeWithdraw() status = %d, want 0", resp.GetStatus())
	}
	if gotAmount == nil || gotAmount.Cmp(big.NewInt(456)) != 0 {
		t.Fatalf("CooperativeWithdraw() amount = %v, want 456", gotAmount)
	}
}

func TestAdminServiceCooperativeWithdrawReturnsCallbackError(t *testing.T) {
	var removed string

	svc := &adminService{
		cooperativeWithdrawWithCallback: func(cid ctype.CidType, amount *big.Int, cb cooperativewithdraw.Callback) (string, error) {
			cb.OnError("withdraw-hash", "withdraw failed")
			return "withdraw-hash", nil
		},
		getChannelBalance:            func(cid ctype.CidType) (*common.ChannelBalance, error) { return nil, nil },
		removeCooperativeWithdrawJob: func(withdrawHash string) { removed = withdrawHash },
	}

	resp, err := svc.CooperativeWithdraw(context.Background(), &rpc.ChannelOpRequest{Cid: "3", Wei: "1"})
	if err == nil {
		t.Fatal("CooperativeWithdraw() error = nil, want grpc error")
	}
	if status.Code(err) != codes.Unavailable {
		t.Fatalf("CooperativeWithdraw() grpc code = %s, want %s", status.Code(err), codes.Unavailable)
	}
	if resp.GetStatus() != 1 || resp.GetError() != "withdraw failed" {
		t.Fatalf("CooperativeWithdraw() response = %+v, want status=1 error=withdraw failed", resp)
	}
	if removed != "withdraw-hash" {
		t.Fatalf("RemoveCooperativeWithdrawJob() hash = %q, want withdraw-hash", removed)
	}
}

func TestAdminServiceCooperativeWithdrawRejectsInvalidCid(t *testing.T) {
	svc := &adminService{
		cooperativeWithdrawWithCallback: func(cid ctype.CidType, amount *big.Int, cb cooperativewithdraw.Callback) (string, error) {
			t.Fatal("cooperativeWithdrawWithCallback should not be called for invalid cid")
			return "", nil
		},
		getChannelBalance:            func(cid ctype.CidType) (*common.ChannelBalance, error) { return nil, nil },
		removeCooperativeWithdrawJob: func(string) {},
	}

	resp, err := svc.CooperativeWithdraw(context.Background(), &rpc.ChannelOpRequest{Cid: "zz", Wei: "1"})
	if err == nil {
		t.Fatal("CooperativeWithdraw() error = nil, want grpc error")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("CooperativeWithdraw() grpc code = %s, want %s", status.Code(err), codes.InvalidArgument)
	}
	if resp.GetStatus() != 1 || resp.GetError() == "" {
		t.Fatalf("CooperativeWithdraw() response = %+v, want status=1 with error", resp)
	}
}

func TestAdminServiceCooperativeWithdrawReturnsStartError(t *testing.T) {
	svc := &adminService{
		cooperativeWithdrawWithCallback: func(cid ctype.CidType, amount *big.Int, cb cooperativewithdraw.Callback) (string, error) {
			return "", errors.New("start failed")
		},
		getChannelBalance:            func(cid ctype.CidType) (*common.ChannelBalance, error) { return nil, nil },
		removeCooperativeWithdrawJob: func(string) {},
	}

	resp, err := svc.CooperativeWithdraw(context.Background(), &rpc.ChannelOpRequest{Cid: "4", Wei: "1"})
	if err == nil {
		t.Fatal("CooperativeWithdraw() error = nil, want grpc error")
	}
	if status.Code(err) != codes.Unavailable {
		t.Fatalf("CooperativeWithdraw() grpc code = %s, want %s", status.Code(err), codes.Unavailable)
	}
	if resp.GetStatus() != 1 || resp.GetError() != "start failed" {
		t.Fatalf("CooperativeWithdraw() response = %+v, want status=1 error=start failed", resp)
	}
}
