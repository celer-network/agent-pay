// Copyright 2018-2025 Celer Network

package webapi

import (
	"context"

	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/webapi/rpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OspPayBackend interface {
	SendToken(*rpc.SendTokenRequest) (ctype.PayIDType, error)
	SendConditionalPayment(*rpc.SendConditionalPaymentRequest) (ctype.PayIDType, error)
	CreateAppSessionOnVirtualContract(*rpc.CreateAppSessionOnVirtualContractRequest) (string, error)
	DeleteAppSession(string) error
	GetStatusForAppSession(string) (uint8, error)
	GetIncomingPaymentState(ctype.PayIDType) (int, error)
	GetIncomingPaymentRecord(ctype.PayIDType) (*PaymentRecord, error)
	GetOutgoingPaymentState(ctype.PayIDType) (int, error)
	ConfirmOutgoingPayment(ctype.PayIDType) error
	RejectIncomingPayment(ctype.PayIDType) error
}

type OspPayApiServer struct {
	rpc.UnimplementedWebApiServer
	backend   OspPayBackend
	eventFeed *PaymentEventFeed
}

func NewOspPayApiServer(backend OspPayBackend, eventFeed *PaymentEventFeed) *OspPayApiServer {
	if eventFeed == nil {
		eventFeed = NewPaymentEventFeed()
	}
	return &OspPayApiServer{backend: backend, eventFeed: eventFeed}
}

func ospAdminBalanceGuidance(method string, ambiguous bool) error {
	msg := method + " is not supported on OSP WebAPI in phase 1; use Admin gRPC CelerGetPeerStatus(peer, token) or osp-cli -dbview channel"
	if ambiguous {
		msg = method + " is not supported on OSP WebAPI in phase 1: peer/channel selection is ambiguous on an OSP; use Admin gRPC CelerGetPeerStatus(peer, token) or osp-cli -dbview channel"
	}
	return status.Error(codes.Unimplemented, msg)
}

func (s *OspPayApiServer) SendToken(
	context context.Context,
	request *rpc.SendTokenRequest) (*rpc.PaymentID, error) {
	payID, err := s.backend.SendToken(request)
	if err != nil {
		return nil, err
	}
	return &rpc.PaymentID{PaymentId: ctype.PayID2Hex(payID)}, nil
}

func (s *OspPayApiServer) SendConditionalPayment(
	context context.Context,
	request *rpc.SendConditionalPaymentRequest) (*rpc.PaymentID, error) {
	payID, err := s.backend.SendConditionalPayment(request)
	if err != nil {
		return nil, err
	}
	return &rpc.PaymentID{PaymentId: ctype.PayID2Hex(payID)}, nil
}

func (s *OspPayApiServer) CreateAppSessionOnVirtualContract(
	context context.Context,
	request *rpc.CreateAppSessionOnVirtualContractRequest) (*rpc.SessionID, error) {
	sessionID, err := s.backend.CreateAppSessionOnVirtualContract(request)
	if err != nil {
		return nil, err
	}
	return &rpc.SessionID{SessionId: sessionID}, nil
}

func (s *OspPayApiServer) DeleteAppSession(
	context context.Context,
	request *rpc.SessionID) (*emptypb.Empty, error) {
	err := s.backend.DeleteAppSession(request.GetSessionId())
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *OspPayApiServer) GetStatusForAppSession(
	context context.Context,
	request *rpc.SessionID) (*rpc.AppSessionStatus, error) {
	statusValue, err := s.backend.GetStatusForAppSession(request.GetSessionId())
	if err != nil {
		return nil, err
	}
	return &rpc.AppSessionStatus{Status: uint32(statusValue)}, nil
}

func (s *OspPayApiServer) GetIncomingPaymentStatus(
	context context.Context,
	request *rpc.PaymentID) (*rpc.PaymentStatus, error) {
	state, err := s.backend.GetIncomingPaymentState(ctype.Hex2PayID(request.PaymentId))
	if err != nil {
		return nil, err
	}
	return &rpc.PaymentStatus{Status: uint32(payStateToSdkStatus(state))}, nil
}

func (s *OspPayApiServer) GetIncomingPaymentInfo(
	context context.Context,
	request *rpc.PaymentID) (*rpc.PaymentInfo, error) {
	record, err := s.backend.GetIncomingPaymentRecord(ctype.Hex2PayID(request.PaymentId))
	if err != nil {
		return nil, err
	}
	return paymentInfoFromRecord(record), nil
}

func (s *OspPayApiServer) GetOutgoingPaymentStatus(
	context context.Context,
	request *rpc.PaymentID) (*rpc.PaymentStatus, error) {
	state, err := s.backend.GetOutgoingPaymentState(ctype.Hex2PayID(request.PaymentId))
	if err != nil {
		return nil, err
	}
	return &rpc.PaymentStatus{Status: uint32(payStateToSdkStatus(state))}, nil
}

func (s *OspPayApiServer) GetBalance(context.Context, *rpc.TokenInfo) (*rpc.GetBalanceResponse, error) {
	return nil, ospAdminBalanceGuidance("GetBalance", true)
}

func (s *OspPayApiServer) GetPeerFreeBalance(context.Context, *rpc.GetPeerFreeBalanceRequest) (*rpc.FreeBalance, error) {
	return nil, ospAdminBalanceGuidance("GetPeerFreeBalance", false)
}

func (s *OspPayApiServer) ConfirmOutgoingPayment(
	context context.Context,
	request *rpc.PaymentID) (*emptypb.Empty, error) {
	err := s.backend.ConfirmOutgoingPayment(ctype.Hex2PayID(request.PaymentId))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *OspPayApiServer) RejectIncomingPayment(
	context context.Context,
	request *rpc.PaymentID) (*emptypb.Empty, error) {
	err := s.backend.RejectIncomingPayment(ctype.Hex2PayID(request.PaymentId))
	if err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *OspPayApiServer) SubscribeIncomingPayments(
	_ *emptypb.Empty,
	stream rpc.WebApi_SubscribeIncomingPaymentsServer) error {
	payments, release, err := s.eventFeed.SubscribeIncoming()
	if err != nil {
		return err
	}
	defer release()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case payment, ok := <-payments:
			if !ok {
				return nil
			}
			if err := stream.Send(paymentInfoFromPayment(payment)); err != nil {
				return err
			}
		}
	}
}

func (s *OspPayApiServer) SubscribeOutgoingPayments(
	_ *emptypb.Empty,
	stream rpc.WebApi_SubscribeOutgoingPaymentsServer) error {
	events, release, err := s.eventFeed.SubscribeOutgoing()
	if err != nil {
		return err
	}
	defer release()

	for {
		select {
		case <-stream.Context().Done():
			return nil
		case event, ok := <-events:
			if !ok {
				return nil
			}
			if err := stream.Send(outgoingPaymentInfoFromPayment(event.payment, event.err)); err != nil {
				return err
			}
		}
	}
}