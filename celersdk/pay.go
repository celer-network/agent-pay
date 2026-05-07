// Copyright 2018-2025 Celer Network

// payment related interface for celer sdk

package celersdk

import (
	"errors"
	"time"

	"github.com/celer-network/agent-pay/celersdkintf"
	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	"github.com/celer-network/agent-pay/utils"
	"github.com/celer-network/goutils/log"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"google.golang.org/protobuf/types/known/anypb"
)

const cPayTimeout = 600 // timeout in seconds for cpay ie. no app channel condition

// noteTypeUrl should be type url of any.Any.
// noteStr should be string representation of []byte in note (any.Any)
func (mc *Client) SendNative(receiver string, amtWei string, noteTypeUrl string, noteValueByte []byte) (string, error) {
	return mc.SendToken(nil, receiver, amtWei, noteTypeUrl, noteValueByte)
}

// SendToken sends an ERC20 or native token to receiver. Caller can optionally add a note in the pay.
func (mc *Client) SendToken(tk *Token, receiver string, amtWei string, noteTypeUrl string, noteValueByte []byte) (string, error) {

	xfer := createXfer(tk, receiver, amtWei)
	note := &anypb.Any{
		TypeUrl: noteTypeUrl,
		Value:   noteValueByte,
	}
	payID, err := mc.c.AddBooleanPay(
		xfer, []*entity.Condition{}, uint64(time.Now().Unix())+cPayTimeout, note, 0)
	if err != nil {
		log.Errorln("SendToken:", err)
		return ctype.ZeroPayIDHex, err
	}
	ret := ctype.PayID2Hex(payID)
	log.Debugln("Sent pay:", ret)
	return ret, nil
}

// ConfirmPay settles the condpay, ie. actually paid to pay dest
func (mc *Client) ConfirmPay(payID string) error {
	return mc.c.ConfirmBooleanPay(ctype.Hex2PayID(payID))
}

// RejectPay cancels the pay, ie. ask OSP and pay src to not pay
func (mc *Client) RejectPay(payID string) error {
	return mc.c.RejectBooleanPay(ctype.Hex2PayID(payID))
}

// RemoveExpiredPays clears pending pays that have expired, if tk is nil, means native token
func (mc *Client) RemoveExpiredPays(tk *Token) error {
	token := sdkToken2entityToken(tk)
	return mc.c.SettleExpiredPays(token)
}

// ResolvePayOnChain settles the payment onchain and receives the payment from OSP.
//
// VIRTUAL_CONTRACT conditions are deployed automatically before the resolve tx
// is submitted: this client walks the pay's conditions and, for any
// VIRTUAL_CONTRACT registered locally via NewAppChannelOnVirtualContract,
// triggers the same deploy path used by OnChainGetBooleanOutcome. Conditions
// registered on a different node are left to that node to deploy.
func (mc *Client) ResolvePayOnChain(payID string) error {
	err := mc.c.ResolveCondPayOnChain(ctype.Hex2PayID(payID))
	if err != nil {
		return err
	}
	return mc.c.SettleOnChainResolvedPay(ctype.Hex2PayID(payID))
}

// ConfirmOnChainResolvedPays confirms pays that have been onchain resolved, if tk is nil, means native token
func (mc *Client) ConfirmOnChainResolvedPays(tk *Token) error {
	token := sdkToken2entityToken(tk)
	return mc.c.ConfirmOnChainResolvedPays(token)
}

// Get incoming payment status code
func (mc *Client) GetIncomingPaymentStatus(payId string) int {
	return mc.c.GetIncomingPaymentStatus(ctype.Hex2PayID(payId))
}

// GetIncomingPaymentInfo returns the related payment info for an incoming payment ID.
// It returns ErrPayNotFound if the payment does not belong to this client as receiver.
func (mc *Client) GetIncomingPaymentInfo(paymentID string) (*celersdkintf.Payment, error) {
	payment, err := mc.c.GetPayment(ethcommon.HexToHash(paymentID))
	if err != nil {
		return nil, err
	}
	if ctype.Hex2Addr(payment.Receiver) != mc.c.GetMyEthAddr() {
		return nil, common.ErrPayNotFound
	}
	return payment, nil
}

// Get outgoing payment status code
func (mc *Client) GetOutgoingPaymentStatus(payId string) int {
	return mc.c.GetOutgoingPaymentStatus(ctype.Hex2PayID(payId))
}

func (mc *Client) GetOnChainPaymentInfo(paymentID string) (*OnChainPaymentInfo, error) {
	amount, resolveDeadline, err := mc.c.GetCondPayInfoFromRegistry(ethcommon.HexToHash(paymentID))
	if err != nil {
		return nil, err
	}
	return &OnChainPaymentInfo{Amount: amount.String(), ResolveDeadline: resolveDeadline}, nil
}

// ResolveIncomingPaymentOnChain submits PayResolver.resolvePaymentByConditions
// for the given payment. Locally-registered VIRTUAL_CONTRACT conditions are
// auto-deployed first (see ResolvePayOnChain doc).
func (mc *Client) ResolveIncomingPaymentOnChain(payId string) error {
	return mc.c.ResolveCondPayOnChain(ctype.Hex2PayID(payId))
}

func (mc *Client) SettleOnChainResolvedIncomingPayment(payId string) error {
	return mc.c.SettleOnChainResolvedPay(ctype.Hex2PayID(payId))
}

func (mc *Client) SendConditionalPayment(
	tokenInfo *TokenInfo,
	destination string,
	amount string,
	transferLogicType TransferLogicType,
	conditions []*Condition,
	timeout int64,
	note *anypb.Any) (string, error) {
	if transferLogicType != transferLogicTypeBooleanAnd {
		return "", errors.New("Unsupported transfer logic type")
	}
	token := &entity.TokenInfo{
		TokenType:    entity.TokenType(int32(tokenInfo.TokenType)),
		TokenAddress: ctype.Hex2Bytes(tokenInfo.TokenAddress),
	}
	transfer := &entity.TokenTransfer{
		Token: token,
		Receiver: &entity.AccountAmtPair{
			Account: ctype.Hex2Bytes(destination),
			Amt:     utils.Wei2BigInt(amount).Bytes(),
		},
	}
	entityConditions := make([]*entity.Condition, len(conditions))
	for i, condition := range conditions {
		entityConditions[i] = conditionToEntityCondition(condition)
	}
	payID, err := mc.c.AddBooleanPay(
		transfer,
		entityConditions,
		uint64(time.Now().Unix())+uint64(timeout),
		note, 0)
	if err != nil {
		log.Error(err)
		return ctype.ZeroPayIDHex, err
	}
	ret := ctype.PayID2Hex(payID)
	log.Debugln("Sent pay:", ret)
	return ret, nil
}

func (mc *Client) SettleExpiredPayments(tokenInfo *TokenInfo) error {
	return mc.c.SettleExpiredPays(&entity.TokenInfo{
		TokenType:    entity.TokenType(int32(tokenInfo.TokenType)),
		TokenAddress: ctype.Hex2Bytes(tokenInfo.TokenAddress),
	})
}

// GetPayment returns the related payment info of a specified payment ID
func (mc *Client) GetPayment(paymentID string) (*celersdkintf.Payment, error) {
	return mc.c.GetPayment(ethcommon.HexToHash(paymentID))
}

// GetAllPayments returns all payments info.
// **CAUTION**: This function costs heavy lookup on several tables and joins
// information from those tables, please take performance into consideration
// before using this.
// PaymentList.PayList is list of all payments. But due to gomobile limitation (no return list).
// mobile app needs to do following
// payList = GetAllPayments()
//
//	for i=0; i<payList.Length; i++ {
//	    pay = payList.Get(i)
//	}
func (mc *Client) GetAllPayments() (*celersdkintf.PaymentList, error) {
	allpays, err := mc.c.GetAllPayments()
	if err != nil {
		return nil, err
	}

	return &celersdkintf.PaymentList{
		Length:  len(allpays),
		PayList: allpays,
	}, nil
}
