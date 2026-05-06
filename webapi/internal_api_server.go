// Copyright 2018-2025 Celer Network

package webapi

import (
	"context"
	"errors"

	"github.com/celer-network/agent-pay/celersdk"
	"github.com/celer-network/agent-pay/entity"
	"github.com/celer-network/agent-pay/webapi/rpc"
	"google.golang.org/grpc"
)

type InternalApiServer struct {
	*ApiServer
	register func(*grpc.Server)
}

func NewInternalApiServer(
	webPort int,
	grpcPort int,
	allowedOrigins string,
	keystore string,
	password string,
	dataPath string,
	config string,
	extSigner bool) *InternalApiServer {
	apiServer := NewApiServer(webPort, grpcPort, allowedOrigins, keystore, password, dataPath, config, extSigner)
	return &InternalApiServer{ApiServer: apiServer}
}

func NewInternalApiServerWithExternalSigner(
	webPort int,
	grpcPort int,
	allowedOrigins string,
	addr string,
	dataPath string,
	config string,
	cb celersdk.ExternalSignerCallback,
	register func(*grpc.Server)) *InternalApiServer {
	apiServer := NewApiServerWithExternalSigner(webPort, grpcPort, allowedOrigins, addr, dataPath, config, cb)
	return &InternalApiServer{ApiServer: apiServer, register: register}
}

func (s *InternalApiServer) Start() {
	gs := grpc.NewServer()
	if s.register != nil {
		s.register(gs)
	}
	rpc.RegisterWebApiServer(gs, s.ApiServer)
	rpc.RegisterInternalWebApiServer(gs, s)
	s.ApiServer.serve(gs)
}

func (s *InternalApiServer) OpenTrustedPaymentChannel(
	context context.Context, request *rpc.OpenPaymentChannelRequest) (*rpc.ChannelID, error) {
	callbackImpl := s.callbackImpl
	tokenInfo := request.TokenInfo
	switch entity.TokenType(tokenInfo.TokenType) {
	case entity.TokenType_ETH:
		go s.apiClient.TcbOpenETHChannel(
			request.PeerAmount,
			s.callbackImpl)
	case entity.TokenType_ERC20:
		go s.apiClient.TcbOpenTokenChannel(
			&celersdk.Token{Erctype: "ERC20", Addr: tokenInfo.TokenAddress},
			request.PeerAmount,
			s.callbackImpl)
	default:
		return nil, errors.New("Unknown token type")
	}
	select {
	case cid := <-callbackImpl.channelOpened:
		return &rpc.ChannelID{ChannelId: cid}, nil
	case errMsg := <-callbackImpl.openChannelError:
		return nil, errors.New(errMsg)
	}
}

func (s *InternalApiServer) InstantiateTrustedPaymentChannel(
	context context.Context, request *rpc.TokenInfo) (*rpc.ChannelID, error) {
	var ercType string
	if request.TokenType == entity.TokenType_ETH {
		ercType = ""
	} else {
		ercType = "ERC20"
	}
	callbackImpl := s.callbackImpl
	go s.apiClient.InstantiateChannelForToken(
		&celersdk.Token{
			Erctype: ercType,
			Addr:    request.TokenAddress,
		},
		callbackImpl)
	select {
	case cid := <-callbackImpl.channelOpened:
		return &rpc.ChannelID{ChannelId: cid}, nil
	case errMsg := <-callbackImpl.openChannelError:
		return nil, errors.New(errMsg)
	}
}

func (s *InternalApiServer) DepositNonBlocking(
	context context.Context, request *rpc.DepositOrWithdrawRequest) (*rpc.DepositOrWithdrawJob, error) {
	return s.ApiServer.DepositNonBlocking(context, request)
}

func (s *InternalApiServer) CooperativeWithdrawNonBlocking(
	context context.Context,
	request *rpc.DepositOrWithdrawRequest) (*rpc.DepositOrWithdrawJob, error) {
	return s.ApiServer.CooperativeWithdrawNonBlocking(context, request)
}
