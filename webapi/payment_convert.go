// Copyright 2018-2025 Celer Network

package webapi

import (
	"math/big"
	"strings"
	"time"

	"github.com/celer-network/agent-pay/celersdkintf"
	enums "github.com/celer-network/agent-pay/common/structs"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	msgrpc "github.com/celer-network/agent-pay/rpc"
	"github.com/celer-network/agent-pay/utils"
	"github.com/celer-network/agent-pay/webapi/rpc"
	"google.golang.org/protobuf/types/known/anypb"
)

type PaymentRecord struct {
	PayID    ctype.PayIDType
	Pay      *entity.ConditionalPay
	Note     *anypb.Any
	PayState int
}

func paymentFromCondPay(
	payID ctype.PayIDType,
	pay *entity.ConditionalPay,
	payNote *anypb.Any,
	status int) *celersdkintf.Payment {
	if pay == nil {
		return &celersdkintf.Payment{}
	}

	payJSON, _ := utils.PbToJSONString(pay)
	payNoteJSON, _ := utils.PbToJSONString(payNote)

	payNoteType := ""
	if payNote != nil {
		typeURL := payNote.GetTypeUrl()
		if i := strings.LastIndex(typeURL, "/"); i >= 0 && i+1 < len(typeURL) {
			payNoteType = typeURL[i+1:]
		} else {
			payNoteType = typeURL
		}
	}

	maxTransfer := pay.TransferFunc.MaxTransfer
	payTimestamp := int64(pay.PayTimestamp / uint64(time.Millisecond))

	payment := &celersdkintf.Payment{
		Sender:       ctype.Bytes2Hex(pay.GetSrc()),
		Receiver:     ctype.Bytes2Hex(pay.GetDest()),
		AmtWei:       new(big.Int).SetBytes(maxTransfer.Receiver.Amt).String(),
		TokenAddr:    ctype.Bytes2Hex(maxTransfer.Token.TokenAddress),
		UID:          ctype.PayID2Hex(payID),
		PayJSON:      payJSON,
		Status:       status,
		PayNoteType:  payNoteType,
		PayNoteJSON:  payNoteJSON,
		PayTimestamp: payTimestamp,
	}
	if maxTransfer.Token.TokenType == entity.TokenType_ETH {
		payment.TokenAddr = ""
	}
	return payment
}

func paymentInfoFromPayment(payment *celersdkintf.Payment) *rpc.PaymentInfo {
	tokenAddr := ctype.Hex2Addr(payment.TokenAddr)
	var tokenType entity.TokenType
	if tokenAddr == ctype.Hex2Addr(ctype.EthTokenAddrStr) {
		tokenType = entity.TokenType_ETH
	} else {
		tokenType = entity.TokenType_ERC20
	}
	return &rpc.PaymentInfo{
		PaymentId: payment.UID,
		Sender:    payment.Sender,
		Receiver:  payment.Receiver,
		TokenInfo: &rpc.TokenInfo{
			TokenType:    tokenType,
			TokenAddress: payment.TokenAddr,
		},
		Amount:      payment.AmtWei,
		PaymentJson: payment.PayJSON,
		Status:      uint32(payment.Status),
	}
}

func paymentInfoFromRecord(record *PaymentRecord) *rpc.PaymentInfo {
	if record == nil {
		return nil
	}
	status := payStateToSdkStatus(record.PayState)
	return paymentInfoFromPayment(paymentFromCondPay(record.PayID, record.Pay, record.Note, status))
}

func outgoingPaymentInfoFromPayment(payment *celersdkintf.Payment, errInfo *celersdkintf.E) *rpc.OutgoingPaymentInfo {
	var errReason string
	var errCode int64
	if errInfo != nil {
		errReason = errInfo.Reason
		errCode = int64(errInfo.Code)
	}
	return &rpc.OutgoingPaymentInfo{
		Payment:     paymentInfoFromPayment(payment),
		ErrorReason: errReason,
		ErrorCode:   errCode,
	}
}

func payStateToSdkStatus(state int) int {
	switch state {
	case enums.PayState_ONESIG_PENDING, enums.PayState_COSIGNED_PENDING:
		return celersdkintf.PAY_STATUS_INITIALIZING
	case enums.PayState_SECRET_REVEALED, enums.PayState_ONESIG_PAID,
		enums.PayState_ONESIG_CANCELED, enums.PayState_INGRESS_REJECTED:
		return celersdkintf.PAY_STATUS_PENDING
	case enums.PayState_COSIGNED_PAID:
		return celersdkintf.PAY_STATUS_PAID
	case enums.PayState_COSIGNED_CANCELED, enums.PayState_NACKED:
		return celersdkintf.PAY_STATUS_UNPAID
	default:
		return celersdkintf.PAY_STATUS_INVALID
	}
}

func settleReasonToSdkStatus(reason msgrpc.PaymentSettleReason) int {
	switch reason {
	case msgrpc.PaymentSettleReason_PAY_EXPIRED:
		return celersdkintf.PAY_STATUS_UNPAID_EXPIRED
	case msgrpc.PaymentSettleReason_PAY_REJECTED:
		return celersdkintf.PAY_STATUS_UNPAID_REJECTED
	case msgrpc.PaymentSettleReason_PAY_RESOLVED_ONCHAIN:
		return celersdkintf.PAY_STATUS_PAID_RESOLVED_ONCHAIN
	case msgrpc.PaymentSettleReason_PAY_PAID_MAX:
		return celersdkintf.PAY_STATUS_PAID
	case msgrpc.PaymentSettleReason_PAY_DEST_UNREACHABLE:
		return celersdkintf.PAY_STATUS_UNPAID_DEST_UNREACHABLE
	default:
		return celersdkintf.PAY_STATUS_INVALID
	}
}