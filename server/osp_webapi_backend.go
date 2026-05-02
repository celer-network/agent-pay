// Copyright 2018-2025 Celer Network

package main

import (
	"bytes"
	"errors"
	"time"

	"github.com/celer-network/agent-pay/cnode"
	"github.com/celer-network/agent-pay/common"
	enums "github.com/celer-network/agent-pay/common/structs"
	"github.com/celer-network/agent-pay/config"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	"github.com/celer-network/agent-pay/utils"
	"github.com/celer-network/agent-pay/webapi"
	webrpc "github.com/celer-network/agent-pay/webapi/rpc"
	"google.golang.org/protobuf/types/known/anypb"
)

// ospWebapiDefaultPayTimeoutSec is the resolve-deadline window used by the
// `SendToken` WebAPI shortcut, which doesn't take a caller-supplied timeout.
// The deadline is `now + this value` (seconds), and must be long enough to
// cover the source→destination forward, the reveal-secret round-trip, and any
// settle-request retry under load. 600s (10 min) gives multi-hop crossnet
// payments comfortable headroom while still being short enough that abandoned
// pays expire promptly.
const ospWebapiDefaultPayTimeoutSec = uint64(600)

var _ webapi.OspPayBackend = (*ospWebapiBackend)(nil)

type ospWebapiBackend struct {
	cNode  *cnode.CNode
	myAddr ctype.Addr
}

func newOspWebapiBackend(cNode *cnode.CNode) *ospWebapiBackend {
	return &ospWebapiBackend{
		cNode:  cNode,
		myAddr: cNode.EthAddress,
	}
}

func (b *ospWebapiBackend) SendToken(request *webrpc.SendTokenRequest) (ctype.PayIDType, error) {
	return b.sendBooleanPayment(request.GetTokenInfo(), request.GetDestination(), request.GetAmount(), nil, ospWebapiDefaultPayTimeoutSec, request.GetNote())
}

func (b *ospWebapiBackend) SendConditionalPayment(request *webrpc.SendConditionalPaymentRequest) (ctype.PayIDType, error) {
	if request.GetTransferLogicType() != entity.TransferFunctionType_BOOLEAN_AND {
		return ctype.ZeroPayID, errors.New("Unsupported transfer logic type")
	}
	conditions := make([]*entity.Condition, len(request.GetConditions()))
	for i, condition := range request.GetConditions() {
		conditions[i] = &entity.Condition{
			ConditionType:          getConditionType(condition.GetOnChainDeployed()),
			DeployedContractAddress: getDeployedContractAddress(condition),
			VirtualContractAddress:  getVirtualContractAddress(condition),
			ArgsQueryFinalization:   condition.GetIsFinalizedArgs(),
			ArgsQueryOutcome:        condition.GetGetOutcomeArgs(),
		}
	}
	return b.sendBooleanPayment(
		request.GetTokenInfo(),
		request.GetDestination(),
		request.GetAmount(),
		conditions,
		request.GetTimeout(),
		request.GetNote())
}

func (b *ospWebapiBackend) CreateAppSessionOnVirtualContract(request *webrpc.CreateAppSessionOnVirtualContractRequest) (string, error) {
	return b.cNode.AppClient.NewAppChannelOnVirtualContract(
		ctype.Hex2Bytes(request.GetContractBin()),
		ctype.Hex2Bytes(request.GetContractConstructor()),
		request.GetNonce(),
	)
}

func (b *ospWebapiBackend) DeleteAppSession(sessionID string) error {
	b.cNode.AppClient.DeleteAppChannel(sessionID)
	return nil
}

func (b *ospWebapiBackend) GetIncomingPaymentState(payID ctype.PayIDType) (int, error) {
	inState, _, _, err := b.cNode.GetDAL().GetPayStates(payID)
	return inState, err
}

func (b *ospWebapiBackend) GetIncomingPaymentRecord(payID ctype.PayIDType) (*webapi.PaymentRecord, error) {
	pay, note, _, inState, _, _, _, found, err := b.cNode.GetDAL().GetPaymentInfo(payID)
	if err != nil {
		return nil, err
	}
	if !found || pay == nil || inState == enums.PayState_NULL || !bytes.Equal(pay.GetDest(), b.myAddr.Bytes()) {
		return nil, common.ErrPayNotFound
	}
	return &webapi.PaymentRecord{PayID: payID, Pay: pay, Note: note, PayState: inState}, nil
}

func (b *ospWebapiBackend) GetOutgoingPaymentState(payID ctype.PayIDType) (int, error) {
	_, outState, _, err := b.cNode.GetDAL().GetPayStates(payID)
	return outState, err
}

func (b *ospWebapiBackend) ConfirmOutgoingPayment(payID ctype.PayIDType) error {
	return b.cNode.ConfirmBooleanPay(payID)
}

func (b *ospWebapiBackend) RejectIncomingPayment(payID ctype.PayIDType) error {
	return b.cNode.RejectBooleanPay(payID)
}

func (b *ospWebapiBackend) sendBooleanPayment(
	tokenInfo *webrpc.TokenInfo,
	destination string,
	amount string,
	conditions []*entity.Condition,
	timeout uint64,
	note *anypb.Any) (ctype.PayIDType, error) {
	transfer, err := buildTokenTransfer(tokenInfo, destination, amount)
	if err != nil {
		return ctype.ZeroPayID, err
	}

	nowTs := uint64(time.Now().Unix())
	resolveDeadline := nowTs + timeout
	if resolveDeadline <= nowTs {
		return ctype.ZeroPayID, common.ErrDeadlinePassed
	}

	pay := &entity.ConditionalPay{
		Src:        b.myAddr.Bytes(),
		Dest:       transfer.Receiver.Account,
		Conditions: conditions,
		TransferFunc: &entity.TransferFunction{
			LogicType:   entity.TransferFunctionType_BOOLEAN_AND,
			MaxTransfer: transfer,
		},
		ResolveDeadline: resolveDeadline,
		ResolveTimeout:  config.PayResolveTimeout,
	}

	var payID ctype.PayIDType
	for i := 0; i < 10; i++ {
		payID, err = b.cNode.AddBooleanPay(pay, note, 0)
		if err != common.ErrPendingSimplex {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	return payID, err
}

func buildTokenTransfer(tokenInfo *webrpc.TokenInfo, destination string, amount string) (*entity.TokenTransfer, error) {
	if tokenInfo == nil {
		return nil, errors.New("Missing token info")
	}
	destAddr, err := utils.ValidateAndFormatAddress(destination)
	if err != nil {
		return nil, err
	}
	amt := utils.Wei2BigInt(amount)
	if amt == nil {
		return nil, common.ErrInvalidAmount
	}
	entityToken, err := webTokenToEntityToken(tokenInfo)
	if err != nil {
		return nil, err
	}
	return &entity.TokenTransfer{
		Token: entityToken,
		Receiver: &entity.AccountAmtPair{
			Account: destAddr.Bytes(),
			Amt:     amt.Bytes(),
		},
	}, nil
}

func webTokenToEntityToken(tokenInfo *webrpc.TokenInfo) (*entity.TokenInfo, error) {
	switch tokenInfo.GetTokenType() {
	case entity.TokenType_ETH:
		return &entity.TokenInfo{TokenType: entity.TokenType_ETH}, nil
	case entity.TokenType_ERC20:
		tokenAddr, err := utils.ValidateAndFormatAddress(tokenInfo.GetTokenAddress())
		if err != nil {
			return nil, err
		}
		return &entity.TokenInfo{TokenType: entity.TokenType_ERC20, TokenAddress: tokenAddr.Bytes()}, nil
	default:
		return nil, errors.New("Unknown token type")
	}
}

func getConditionType(onChain bool) entity.ConditionType {
	if onChain {
		return entity.ConditionType_DEPLOYED_CONTRACT
	}
	return entity.ConditionType_VIRTUAL_CONTRACT
}

func getDeployedContractAddress(condition *webrpc.Condition) []byte {
	if condition.GetOnChainDeployed() {
		return ctype.Hex2Bytes(condition.GetContractAddress())
	}
	return nil
}

func getVirtualContractAddress(condition *webrpc.Condition) []byte {
	if !condition.GetOnChainDeployed() {
		return ctype.Hex2Bytes(condition.GetContractAddress())
	}
	return nil
}