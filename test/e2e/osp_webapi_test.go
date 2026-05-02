// Copyright 2018-2025 Celer Network

package e2e

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/celer-network/agent-pay/celersdkintf"
	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	"github.com/celer-network/agent-pay/storage"
	tf "github.com/celer-network/agent-pay/testing"
	"github.com/celer-network/agent-pay/testing/testapp"
	"github.com/celer-network/agent-pay/tools/toolsetup"
	webrpc "github.com/celer-network/agent-pay/webapi/rpc"
	"github.com/celer-network/goutils/log"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

func TestOSPWebApi(t *testing.T) {
	toKill := setUpWithServerArgs("-webapigrpc", sWebApiRPC)
	defer tearDownMultiSvr(toKill)

	t.Run("ospWebApiAppSessionSubset", ospWebApiAppSessionSubset)
	t.Run("ospWebApiPaySubset", ospWebApiPaySubset)
}

func TestOSPWebApiRoutingBehavior(t *testing.T) {
	toKill, dal, err := setUpOspWebApiRoutingOsps()
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownMultiSvr(toKill)

	ks, addrs, err := tf.CreateAccountsWithBalance(2, accountBalance)
	if err != nil {
		t.Fatal(err)
	}
	c1KeyStore := ks[0]
	c1EthAddr := addrs[0]
	c2KeyStore := ks[1]
	c2EthAddr := addrs[1]

	c1, err := tf.StartC1WithoutProxy(c1KeyStore)
	if err != nil {
		t.Fatal(err)
	}
	defer c1.Kill()

	c2, err := startClientForOsp2(c2KeyStore)
	if err != nil {
		t.Fatal(err)
	}
	defer c2.Kill()

	if _, err = c1.OpenChannel(c1EthAddr, entity.TokenType_ETH, tokenAddrEth, initialBalance, initialBalance); err != nil {
		t.Fatal(err)
	}
	if _, err = c2.OpenChannel(c2EthAddr, entity.TokenType_ETH, tokenAddrEth, initialBalance, initialBalance); err != nil {
		t.Fatal(err)
	}
	if err = buildRoutingTablesForEth(o1AdminWeb, o2AdminWeb); err != nil {
		t.Fatal(err)
	}
	time.Sleep(2 * time.Second)

	conn, ospClient, err := tf.DialWebApiClient(sWebApiRPC)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	directResp, err := ospClient.SendToken(context.Background(), &webrpc.SendTokenRequest{
		TokenInfo:   &webrpc.TokenInfo{TokenType: entity.TokenType_ETH, TokenAddress: tokenAddrEth},
		Destination: c1EthAddr,
		Amount:      sendAmt,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err = assertOspWebApiOutgoingPayRouteShape(dal, directResp.GetPaymentId(), true); err != nil {
		t.Fatal(err)
	}
	if err = waitForOspOutgoingPaymentCompletion(directResp.GetPaymentId(), ospClient, c1); err != nil {
		t.Fatal(err)
	}

	routedResp, err := ospClient.SendToken(context.Background(), &webrpc.SendTokenRequest{
		TokenInfo:   &webrpc.TokenInfo{TokenType: entity.TokenType_ETH, TokenAddress: tokenAddrEth},
		Destination: c2EthAddr,
		Amount:      sendAmt,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err = assertOspWebApiOutgoingPayRouteShape(dal, routedResp.GetPaymentId(), false); err != nil {
		t.Fatal(err)
	}
	if err = waitForOspOutgoingPaymentCompletion(routedResp.GetPaymentId(), ospClient, c2); err != nil {
		t.Fatal(err)
	}
}

func ospWebApiPaySubset(t *testing.T) {
	log.Info("============== start test ospWebApiPaySubset ==============")
	defer log.Info("============== end test ospWebApiPaySubset ==============")

	ks, addrs, err := tf.CreateAccountsWithBalance(1, accountBalance)
	if err != nil {
		t.Fatal(err)
	}
	c1KeyStore := ks[0]
	c1EthAddr := addrs[0]

	c1, err := tf.StartC1WithoutProxy(c1KeyStore)
	if err != nil {
		t.Fatal(err)
	}
	defer c1.Kill()

	_, err = c1.OpenChannel(c1EthAddr, entity.TokenType_ETH, tokenAddrEth, initialBalance, initialBalance)
	if err != nil {
		t.Fatal(err)
	}

	conn, ospClient, err := tf.DialWebApiClient(sWebApiRPC)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	incomingSub, err := ospClient.SubscribeIncomingPayments(context.Background(), &empty.Empty{})
	if err != nil {
		t.Fatal(err)
	}
	incomingEvents := make(chan string, 4)
	go func() {
		for {
			payment, err := incomingSub.Recv()
			if err != nil {
				return
			}
			incomingEvents <- payment.GetPaymentId()
		}
	}()

	outgoingResp, err := ospClient.SendToken(context.Background(), &webrpc.SendTokenRequest{
		TokenInfo:   &webrpc.TokenInfo{TokenType: entity.TokenType_ETH, TokenAddress: tokenAddrEth},
		Destination: c1EthAddr,
		Amount:      sendAmt,
	})
	if err != nil {
		t.Fatal(err)
	}

	err = waitForOspOutgoingPaymentCompletion(outgoingResp.GetPaymentId(), ospClient, c1)
	if err != nil {
		t.Fatal(err)
	}

	if _, err = ospClient.GetIncomingPaymentInfo(context.Background(), &webrpc.PaymentID{PaymentId: outgoingResp.GetPaymentId()}); err == nil {
		t.Fatal("OSP unexpectedly fetched outgoing pay via GetIncomingPaymentInfo")
	}

	appChanID, err := c1.NewAppChannelOnVirtualContract(
		ctype.Hex2Bytes(testapp.BooleanCondMockBin),
		[]byte{},
		1006)
	if err != nil {
		t.Fatal(err)
	}

	cond := &entity.Condition{
		ConditionType:          entity.ConditionType_VIRTUAL_CONTRACT,
		VirtualContractAddress: ctype.Hex2Bytes(appChanID),
		ArgsQueryFinalization:  []byte{},
		ArgsQueryOutcome:       []byte{2},
	}
	incomingPayID, err := c1.SendPaymentWithBooleanConditions(
		ospEthAddr,
		sendAmt,
		entity.TokenType_ETH,
		tokenAddrEth,
		[]*entity.Condition{cond},
		100)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForOspIncomingPaymentPending(incomingPayID, c1, ospClient)
	if err != nil {
		t.Fatal(err)
	}

	err = waitForWebApiPaymentEvent(incomingEvents, incomingPayID)
	if err != nil {
		t.Fatal(err)
	}

	incomingInfo, err := ospClient.GetIncomingPaymentInfo(context.Background(), &webrpc.PaymentID{PaymentId: incomingPayID})
	if err != nil {
		t.Fatal(err)
	}
	if incomingInfo.GetPaymentId() != incomingPayID ||
		ctype.Hex2Addr(incomingInfo.GetSender()) != ctype.Hex2Addr(c1EthAddr) ||
		ctype.Hex2Addr(incomingInfo.GetReceiver()) != ctype.Hex2Addr(ospEthAddr) ||
		incomingInfo.GetAmount() != sendAmt ||
		int(incomingInfo.GetStatus()) != celersdkintf.PAY_STATUS_PENDING {
		t.Fatalf("wrong OSP incoming payment info: %+v", incomingInfo)
	}

	_, err = ospClient.RejectIncomingPayment(context.Background(), &webrpc.PaymentID{PaymentId: incomingPayID})
	if err != nil {
		t.Fatal(err)
	}

	err = waitForOspIncomingPaymentFinalized(incomingPayID, c1, ospClient)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ospClient.GetBalance(context.Background(), &webrpc.TokenInfo{TokenType: entity.TokenType_ETH, TokenAddress: tokenAddrEth})
	if status.Code(err) != codes.Unimplemented {
		t.Fatalf("GetBalance error code = %v, want %v (err=%v)", status.Code(err), codes.Unimplemented, err)
	}
}

func ospWebApiAppSessionSubset(t *testing.T) {
	log.Info("============== start test ospWebApiAppSessionSubset ==============")
	defer log.Info("============== end test ospWebApiAppSessionSubset ==============")

	ks, addrs, err := tf.CreateAccountsWithBalance(1, accountBalance)
	if err != nil {
		t.Fatal(err)
	}
	c1KeyStore := ks[0]
	c1EthAddr := addrs[0]

	c1, err := tf.StartC1WithoutProxy(c1KeyStore)
	if err != nil {
		t.Fatal(err)
	}
	defer c1.Kill()

	_, err = c1.OpenChannel(c1EthAddr, entity.TokenType_ETH, tokenAddrEth, initialBalance, initialBalance)
	if err != nil {
		t.Fatal(err)
	}

	conn, ospClient, err := tf.DialWebApiClient(sWebApiRPC)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	// This test exercises only the create→pay→reject→delete cycle on the OSP
	// WebAPI; it never disputes or queries the registered contract, so the
	// underlying bytecode is incidental.
	sessionResp, err := ospClient.CreateAppSessionOnVirtualContract(context.Background(), &webrpc.CreateAppSessionOnVirtualContractRequest{
		ContractBin:         testapp.BooleanCondMockBin,
		ContractConstructor: "",
		Nonce:               1007,
	})
	if err != nil {
		t.Fatal(err)
	}
	sessionID := sessionResp.GetSessionId()
	if sessionID == "" {
		t.Fatal("CreateAppSessionOnVirtualContract returned empty session id")
	}

	payResp, err := ospClient.SendConditionalPayment(context.Background(), &webrpc.SendConditionalPaymentRequest{
		TokenInfo:         &webrpc.TokenInfo{TokenType: entity.TokenType_ETH, TokenAddress: tokenAddrEth},
		Destination:       c1EthAddr,
		Amount:            sendAmt,
		TransferLogicType: entity.TransferFunctionType_BOOLEAN_AND,
		Conditions:        []*webrpc.Condition{{OnChainDeployed: false, ContractAddress: sessionID, IsFinalizedArgs: []byte{}, GetOutcomeArgs: []byte{2}}},
		Timeout:           100,
	})
	if err != nil {
		t.Fatal(err)
	}
	if payResp.GetPaymentId() == "" {
		t.Fatal("SendConditionalPayment returned empty payment id")
	}

	err = waitForOspOutgoingPaymentPending(payResp.GetPaymentId(), ospClient, c1)
	if err != nil {
		t.Fatal(err)
	}

	err = c1.RejectBooleanPay(payResp.GetPaymentId())
	if err != nil {
		t.Fatal(err)
	}

	err = waitForOspOutgoingPaymentCompletion(payResp.GetPaymentId(), ospClient, c1)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ospClient.DeleteAppSession(context.Background(), &webrpc.SessionID{SessionId: sessionID})
	if err != nil {
		t.Fatal(err)
	}
}

func waitForOspOutgoingPaymentCompletion(payID string, ospClient webrpc.WebApiClient, receiver *tf.ClientController) error {
	time.Sleep(200 * time.Millisecond)
	const retryLimit = 20
	var statusCode int
	for retry := 0; retry < retryLimit; retry++ {
		statusResp, err := ospClient.GetOutgoingPaymentStatus(context.Background(), &webrpc.PaymentID{PaymentId: payID})
		if err != nil {
			return err
		}
		statusCode = int(statusResp.GetStatus())
		if payStatusFinalized(statusCode) {
			break
		}
		time.Sleep(400 * time.Millisecond)
		if retry == retryLimit-1 {
			return fmt.Errorf("OSP outgoing payment not finalized, payID %s %d", payID, statusCode)
		}
	}
	if receiver == nil {
		return nil
	}
	for retry := 0; retry < retryLimit; retry++ {
		recvStatus, err := receiver.GetIncomingPaymentStatus(payID)
		if err != nil {
			return err
		}
		if payStatusFinalized(recvStatus) {
			return nil
		}
		time.Sleep(400 * time.Millisecond)
	}
	return fmt.Errorf("receiver payment not finalized, payID %s", payID)
}

func waitForOspOutgoingPaymentPending(payID string, ospClient webrpc.WebApiClient, receiver *tf.ClientController) error {
	time.Sleep(200 * time.Millisecond)
	const retryLimit = 20
	for retry := 0; retry < retryLimit; retry++ {
		statusResp, err := ospClient.GetOutgoingPaymentStatus(context.Background(), &webrpc.PaymentID{PaymentId: payID})
		if err != nil {
			return err
		}
		recvStatus, err := receiver.GetIncomingPaymentStatus(payID)
		if err != nil {
			return err
		}
		if int(statusResp.GetStatus()) == celersdkintf.PAY_STATUS_PENDING && recvStatus == celersdkintf.PAY_STATUS_PENDING {
			return nil
		}
		time.Sleep(400 * time.Millisecond)
	}
	return fmt.Errorf("OSP outgoing payment did not reach pending state, payID %s", payID)
}

func waitForOspIncomingPaymentPending(payID string, sender *tf.ClientController, ospClient webrpc.WebApiClient) error {
	time.Sleep(200 * time.Millisecond)
	const retryLimit = 20
	for retry := 0; retry < retryLimit; retry++ {
		sendStatus, err := sender.GetOutgoingPaymentStatus(payID)
		if err != nil {
			return err
		}
		recvStatusResp, err := ospClient.GetIncomingPaymentStatus(context.Background(), &webrpc.PaymentID{PaymentId: payID})
		if err != nil {
			return err
		}
		if sendStatus == celersdkintf.PAY_STATUS_PENDING && int(recvStatusResp.GetStatus()) == celersdkintf.PAY_STATUS_PENDING {
			return nil
		}
		time.Sleep(400 * time.Millisecond)
	}
	return fmt.Errorf("OSP incoming payment did not reach pending state, payID %s", payID)
}

func waitForOspIncomingPaymentFinalized(payID string, sender *tf.ClientController, ospClient webrpc.WebApiClient) error {
	time.Sleep(200 * time.Millisecond)
	const retryLimit = 20
	for retry := 0; retry < retryLimit; retry++ {
		sendStatus, err := sender.GetOutgoingPaymentStatus(payID)
		if err != nil {
			return err
		}
		recvStatusResp, err := ospClient.GetIncomingPaymentStatus(context.Background(), &webrpc.PaymentID{PaymentId: payID})
		if err != nil {
			return err
		}
		if payStatusFinalized(sendStatus) && payStatusFinalized(int(recvStatusResp.GetStatus())) {
			return nil
		}
		time.Sleep(400 * time.Millisecond)
	}
	return fmt.Errorf("OSP incoming payment did not finalize, payID %s", payID)
}

func waitForWebApiPaymentEvent(events <-chan string, payID string) error {
	timeout := time.After(10 * time.Second)
	for {
		select {
		case got := <-events:
			if got == payID {
				return nil
			}
		case <-timeout:
			return fmt.Errorf("timeout waiting for payment event %s", payID)
		}
	}
}

func setUpOspWebApiRoutingOsps() ([]Killable, *storage.DAL, error) {
	if err := tf.RegisterRouters([]string{osp2Keystore}); err != nil {
		return nil, nil, err
	}
	os.RemoveAll(sStoreDir)

	killables := make([]Killable, 0, 2)
	o1 := tf.StartServerController(outRootDir+toBuild["server"],
		"-profile", noProxyProfile,
		"-port", o1Port,
		"-storedir", sStoreDir,
		"-ks", ospKeystore,
		"-nopassword",
		"-rtc", rtConfigMultiOSP,
		"-svrname", "o1",
		"-logcolor",
		"-logprefix", "o1_"+ospEthAddr[:4],
		"-webapigrpc", sWebApiRPC)
	killables = append(killables, o1)

	o2 := tf.StartServerController(outRootDir+toBuild["server"],
		"-profile", noProxyProfile,
		"-port", o2Port,
		"-storedir", sStoreDir,
		"-adminrpc", o2AdminRPC,
		"-adminweb", o2AdminWeb,
		"-ks", osp2Keystore,
		"-nopassword",
		"-rtc", rtConfigMultiOSP,
		"-svrname", "o2",
		"-logcolor",
		"-logprefix", "o2_"+osp2EthAddr[:4])
	killables = append(killables, o2)

	cleanupErr := func(err error) ([]Killable, *storage.DAL, error) {
		tearDownMultiSvr(killables)
		return nil, nil, err
	}

	time.Sleep(3 * time.Second)
	if err := registerStreamWithRetry(o2AdminWeb, ctype.Hex2Addr(ospEthAddr), localhost+o1Port); err != nil {
		return cleanupErr(err)
	}
	if err := ensureOpenChannel(o2AdminWeb, osp1EthAddr, initOspToOspBalance, initOspToOspBalance, tokenAddrEth); err != nil {
		return cleanupErr(err)
	}
	if err := buildRoutingTablesForEth(o1AdminWeb, o2AdminWeb); err != nil {
		return cleanupErr(err)
	}
	sleep(6)

	return killables, openOspStoreDAL(ospEthAddr), nil
}

func openOspStoreDAL(ospAddr string) *storage.DAL {
	return toolsetup.NewDAL(&common.CProfile{StoreDir: sStoreDir + "/" + ospAddr})
}

func assertOspWebApiOutgoingPayRouteShape(dal *storage.DAL, payID string, expectDirect bool) error {
	pay, _, found, err := dal.GetPayment(ctype.Hex2PayID(payID))
	if err != nil {
		return err
	}
	if !found || pay == nil {
		return fmt.Errorf("payment %s not found in OSP store", payID)
	}

	conditions := pay.GetConditions()
	if expectDirect {
		if len(conditions) != 0 {
			return fmt.Errorf("payment %s expected direct-pay fast path with no prepended conditions, got %d", payID, len(conditions))
		}
		return nil
	}

	if len(conditions) == 0 {
		return fmt.Errorf("payment %s expected prepended hash-lock condition for routed pay, got none", payID)
	}
	hashLock := conditions[0].GetHashLock()
	if len(hashLock) == 0 {
		return fmt.Errorf("payment %s expected first condition to be hash-lock, got %#v", payID, conditions[0])
	}
	_, found, err = dal.GetSecret(ctype.Bytes2Hex(hashLock))
	if err != nil {
		return err
	}
	if !found {
		return fmt.Errorf("payment %s expected stored secret for routed hash-lock %x", payID, hashLock)
	}
	return nil
}
