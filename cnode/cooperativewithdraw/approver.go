// Copyright 2018-2025 Celer Network

package cooperativewithdraw

import (
	"errors"

	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/rpc"
	"github.com/celer-network/agent-pay/utils"
	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
	"google.golang.org/protobuf/proto"
)

func (p *Processor) ProcessRequest(frame *common.MsgFrame) error {
	request := frame.Message.GetWithdrawRequest()
	cid := ctype.Bytes2Cid(request.WithdrawInfo.ChannelId)
	frame.LogEntry.FromCid = ctype.Cid2Hex(cid)
	return p.sendResponse(cid, request)
}

func (p *Processor) sendResponse(
	cid ctype.CidType, request *rpc.CooperativeWithdrawRequest) error {
	withdrawInfo := request.WithdrawInfo

	peer, found, err := p.dal.GetChanPeer(cid)
	if err != nil {
		return err
	}
	if !found {
		return common.ErrChannelNotFound
	}
	serializedInfo, err := proto.Marshal(withdrawInfo)
	if err != nil {
		return err
	}
	if !eth.IsSignatureValid(peer, serializedInfo, request.RequesterSig) {
		return errors.New("Invalid CooperativeWithdrawRequest signature")
	}

	approverSig, err := p.signer.SignEthMessage(serializedInfo)
	if err != nil {
		return err
	}

	err = p.dal.Transactional(p.checkWithdrawBalanceTx, cid, withdrawInfo)
	if err != nil {
		return err
	}

	response := &rpc.CooperativeWithdrawResponse{
		WithdrawInfo: withdrawInfo,
		RequesterSig: request.RequesterSig,
		ApproverSig:  approverSig,
	}
	msg := &rpc.CelerMsg{
		Message: &rpc.CelerMsg_WithdrawResponse{
			WithdrawResponse: response,
		},
	}
	log.Infof("Sending cooperative withdraw response to %s. %s",
		ctype.Addr2Hex(peer), utils.PrintCooperativeWithdrawInfo(withdrawInfo))
	err = p.streamWriter.WriteCelerMsg(peer, msg)
	if err != nil {
		return err
	}
	return nil
}
