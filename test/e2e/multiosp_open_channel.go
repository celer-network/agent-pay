// Copyright 2018-2025 Celer Network

package e2e

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/common/structs"
	"github.com/celer-network/agent-pay/config"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/rpc"
	tf "github.com/celer-network/agent-pay/testing"
	"github.com/celer-network/agent-pay/utils"
	"github.com/celer-network/goutils/log"
	"google.golang.org/grpc"
)

func multiOspOpenChannelTest(t *testing.T) {
	log.Info("============== start test multiOspOpenChannelTest ==============")
	defer log.Info("============== end test multiOspOpenChannelTest ==============")
	// Let osp2 initiate openning channel with osp1.
	err := ensureOpenChannel(o2AdminWeb, osp1EthAddr, initOspToOspBalance, initOspToOspBalance, tokenAddrNative)
	if err != nil {
		t.Error(err)
		return
	}
	if err = buildRoutingTablesForNative(o1AdminWeb, o2AdminWeb); err != nil {
		t.Error(err)
		return
	}
	log.Infoln("done open channel, waiting")
	sleep(6)
	log.Infoln("sending token")
	// requestSvrSendToken is defined in admin.go. It will ask osp1 to send 1 token to osp2EthAddr.
	payID, err := requestSendToken(o1AdminWeb, osp2EthAddr, "1", tokenAddrNative)
	if err != nil {
		t.Error(err)
		return
	}
	log.Infoln("done sending token, waiting")
	time.Sleep(1 * time.Second)
	dal1, dal2, _, _, _ := getMultiOspDALs()
	token := utils.GetTokenInfoFromAddress(ctype.Hex2Addr(tokenAddrNative))
	cid12, found, err := dal1.GetCidByPeerToken(ctype.Hex2Addr(osp2EthAddr), token)
	if err != nil {
		t.Error(err)
		return
	}
	if !found {
		t.Error("channel cid12 not found")
		return
	}
	err = checkOspPayState(dal1, payID, ctype.ZeroCid, structs.PayState_NULL, cid12, structs.PayState_COSIGNED_PAID, 5)
	if err != nil {
		t.Errorf("pay err at o1: %s", err)
		return
	}
	err = checkOspPayState(dal2, payID, cid12, structs.PayState_COSIGNED_PAID, ctype.ZeroCid, structs.PayState_NULL, 5)
	if err != nil {
		t.Errorf("pay err at o2: %s", err)
		return
	}

	log.Infoln("sending back")
	requestSendToken(o2AdminWeb, osp1EthAddr, "1", "")
}

func multiOspOpenChannelPolicyTest(t *testing.T) {
	log.Info("============== start test multiOspOpenChannelPolicyTest ==============")
	defer log.Info("============== end test multiOspOpenChannelPolicyTest ==============")
	tf.FundAccountsWithErc20(tokenAddrErc20, []string{osp2EthAddr}, accountBalance)
	// Let osp2 initiate openning channel with osp1 using bad deposit combination.
	err := requestOpenChannel(o2AdminWeb, osp1EthAddr, "20000000000000000000", "20000000000000000000", tokenAddrNative)
	if err == nil {
		t.Error("Expect to break policy")
		return
	}
	// ask osp1 to deposit 8. This should break ratio policy which is set to 1.0 in rt_config.json
	err = requestOpenChannel(o2AdminWeb, osp1EthAddr, "8", "1", tokenAddrNative)
	if err == nil {
		t.Error("Expect to break matching ratio policy")
		return
	}
}

func multiOspOpenChannelPolicyFallbackTest(t *testing.T) {
	log.Info("============== start test multiOspOpenChannelPolicyFallbackTest ==============")
	defer log.Info("============== end test multiOspOpenChannelPolicyFallbackTest ==============")
	tf.FundAccountsWithErc20(tokenAddrErc20, []string{osp2EthAddr}, accountBalance)
	// Let osp2 initiate openning channel with osp1 using erc20. This should fallback to client-osp policy
	// and fail because it doesn't meet the fallback policy.
	err := requestOpenChannel(o2AdminWeb, osp1EthAddr, "2" /*peerDeposit*/, "2" /*selfDeposit*/, tokenAddrErc20)
	if err == nil {
		t.Error("Expect to fail due to exceeding deposit in fallback")
		return
	}
	// Let osp2 initiate openning channel with osp1 using erc20. Deposit meets fallback aka client-osp policy
	err = ensureOpenChannel(o2AdminWeb, osp1EthAddr, "1" /*peerDeposit*/, "1" /*selfDeposit*/, tokenAddrErc20)
	if err != nil {
		t.Error("Unable to fallback", err)
		return
	}
}

func ensureOpenChannel(adminWebAddr, peerAddr, peerDeposit, selfDeposit, tokenAddr string) error {
	var lastErr error
	for attempt := 0; attempt < 60; attempt++ {
		err := requestOpenChannel(adminWebAddr, peerAddr, peerDeposit, selfDeposit, tokenAddr)
		if err == nil || strings.Contains(err.Error(), "channel already exist") {
			return nil
		}
		if strings.Contains(err.Error(), "no RPC connection") || strings.Contains(err.Error(), "Deadline out of range") {
			lastErr = err
			time.Sleep(500 * time.Millisecond)
			continue
		}
		return err
	}
	return lastErr
}

func registerStreamWithRetry(adminWebAddr string, peerAddr ctype.Addr, peerHostPort string) error {
	var lastErr error
	for attempt := 0; attempt < 40; attempt++ {
		err := utils.RequestRegisterStream(adminWebAddr, peerAddr, peerHostPort)
		if err == nil || strings.Contains(err.Error(), "celer stream already exists") {
			return nil
		}
		lastErr = err
		time.Sleep(250 * time.Millisecond)
	}
	return lastErr
}

func buildRoutingTablesForNative(adminWebAddrs ...string) error {
	for _, adminWebAddr := range adminWebAddrs {
		if err := utils.RequestBuildRoutingTable(adminWebAddr, ctype.NativeTokenAddr); err != nil {
			return err
		}
	}
	return nil
}

func requestOpenChannel(adminWebAddr, peerAddr, peerDeposit, selfDeposit, tokenAddr string) error {
	peerDepositInt, ok := new(big.Int).SetString(peerDeposit, 10)
	if !ok {
		return common.ErrInvalidArg
	}
	selfDepositInt, ok := new(big.Int).SetString(selfDeposit, 10)
	if !ok {
		return common.ErrInvalidArg
	}
	return utils.RequestOpenChannel(adminWebAddr, ctype.Hex2Addr(peerAddr), ctype.Hex2Addr(tokenAddr), peerDepositInt, selfDepositInt)
}

func getNativeBalance(ospHTTPTarget string, osp2Addr string) (string, error) {
	conn, err := grpc.Dial(ospHTTPTarget, utils.GetClientTlsOption(), grpc.WithBlock(),
		grpc.WithTimeout(8*time.Second), grpc.WithKeepaliveParams(config.KeepAliveClientParams))
	if err != nil {
		return "", fmt.Errorf("fail to get peer status: %s", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()
	response, err := rpc.NewRpcClient(conn).CelerGetPeerStatus(
		ctx,
		&rpc.PeerAddress{
			Address:   osp2Addr,
			TokenAddr: tokenAddrNative,
		},
	)
	if err != nil {
		return "", fmt.Errorf("fail to get peer status: %s", err)
	}
	return response.GetFreeBalance(), nil
}
