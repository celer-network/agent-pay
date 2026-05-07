package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"time"

	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/entity"
	"github.com/celer-network/agent-pay/rpc"
	"github.com/celer-network/goutils/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
)

var ErrHttpReponse = errors.New("http response error")

func RequestSendToken(
	adminHostPort string, receiver, tokenAddr ctype.Addr, amount *big.Int, dstNetId uint64) (ctype.PayIDType, error) {
	return RequestSendTokenWithNote(adminHostPort, receiver, tokenAddr, amount, dstNetId, "", nil)
}

func RequestSendTokenWithNote(
	adminHostPort string,
	receiver, tokenAddr ctype.Addr, amount *big.Int, dstNetId uint64,
	noteTypeUrl string, noteValueByte []byte) (ctype.PayIDType, error) {

	request := &rpc.SendTokenRequest{
		DstAddr:   ctype.Addr2Hex(receiver),
		AmtWei:    amount.String(),
		TokenAddr: ctype.Addr2Hex(tokenAddr),
		DstNetId:  dstNetId,
	}
	if noteTypeUrl != "" {
		request.Note = &anypb.Any{
			TypeUrl: noteTypeUrl,
			Value:   noteValueByte,
		}
	}

	url := fmt.Sprintf("http://%s/admin/sendtoken", adminHostPort)
	resBody, err := HttpPost(url, request)
	if err != nil {
		if errors.Is(err, ErrHttpReponse) {
			err = fmt.Errorf("%w, err msg: %s", err, getGrpcHttpErrMsg(resBody))
		}
		return ctype.ZeroPayID, err
	}

	res := &rpc.SendTokenResponse{}
	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(resBody, res)
	if err != nil {
		return ctype.ZeroPayID, err
	}
	if res.Status != 0 {
		return ctype.ZeroPayID, fmt.Errorf("%s", res.Error)
	}
	return ctype.Hex2PayID(res.PayId), nil
}

func RequestRegisterStream(adminHostPort string, peerAddr ctype.Addr, peerHostPort string) error {
	request := &rpc.RegisterStreamRequest{
		PeerRpcAddress: peerHostPort,
		PeerEthAddress: peerAddr.Bytes(),
	}
	url := fmt.Sprintf("http://%s/admin/peer/registerstream", adminHostPort)
	resBody, err := HttpPost(url, request)
	if errors.Is(err, ErrHttpReponse) {
		return fmt.Errorf("%w, err msg: %s", err, getGrpcHttpErrMsg(resBody))
	}
	return err
}

func RequestOpenChannel(adminHostPort string, peerAddr, tokenAddr ctype.Addr, peerDeposit, selfDeposit *big.Int) error {
	tokenType := entity.TokenType_ERC20
	if tokenAddr == ctype.NativeTokenAddr {
		tokenType = entity.TokenType_NATIVE
	}
	request := &rpc.OspOpenChannelRequest{
		PeerEthAddress:    peerAddr.Bytes(),
		TokenType:         tokenType,
		TokenAddress:      tokenAddr.Bytes(),
		PeerDepositAmtWei: peerDeposit.String(),
		SelfDepositAmtWei: selfDeposit.String(),
	}
	url := fmt.Sprintf("http://%s/admin/peer/openchannel", adminHostPort)
	resBody, err := HttpPost(url, request)
	if errors.Is(err, ErrHttpReponse) {
		return fmt.Errorf("%w, err msg: %s", err, getGrpcHttpErrMsg(resBody))
	}
	return err
}

func RequestDeposit(
	adminHostPort string, peerAddr, tokenAddr ctype.Addr, amount *big.Int, toPeer bool, maxWaitSec uint64) (string, error) {
	request := &rpc.DepositRequest{
		PeerAddr:  ctype.Addr2Hex(peerAddr),
		TokenAddr: ctype.Addr2Hex(tokenAddr),
		ToPeer:    toPeer,
		AmtWei:    amount.String(),
		MaxWaitS:  maxWaitSec,
	}
	url := fmt.Sprintf("http://%s/admin/deposit", adminHostPort)
	resBody, err := HttpPost(url, request)
	if err != nil {
		if errors.Is(err, ErrHttpReponse) {
			err = fmt.Errorf("%w, err msg: %s", err, getGrpcHttpErrMsg(resBody))
		}
		return "", err
	}
	res := &rpc.DepositResponse{}
	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(resBody, res)
	if err != nil {
		return "", err
	}
	if res.Status != 0 {
		return "", fmt.Errorf("%s", res.Error)
	}
	return res.DepositId, nil
}

func QueryDeposit(adminHostPort string, depositID string) (*rpc.QueryDepositResponse, error) {
	request := &rpc.QueryDepositRequest{DepositId: depositID}
	url := fmt.Sprintf("http://%s/admin/query_deposit", adminHostPort)
	resBody, err := HttpPost(url, request)
	if err != nil {
		return nil, err
	}
	res := &rpc.QueryDepositResponse{}
	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(resBody, res)
	if err != nil {
		if errors.Is(err, ErrHttpReponse) {
			err = fmt.Errorf("%w, err msg: %s", err, getGrpcHttpErrMsg(resBody))
		}
		return nil, err
	}
	return res, nil
}

func QueryPeerOsps(adminHostPort string) (*rpc.PeerOspsResponse, error) {
	url := fmt.Sprintf("http://%s/admin/peer/peer_osps", adminHostPort)
	resBody, err := HttpPost(url, nil)
	if err != nil {
		return nil, err
	}
	res := &rpc.PeerOspsResponse{}
	err = protojson.UnmarshalOptions{DiscardUnknown: true}.Unmarshal(resBody, res)
	if err != nil {
		if errors.Is(err, ErrHttpReponse) {
			err = fmt.Errorf("%w, err msg: %s", err, getGrpcHttpErrMsg(resBody))
		}
		return nil, err
	}
	return res, nil
}

func RequestBuildRoutingTable(adminHostPort string, tokenAddr ctype.Addr) error {
	request := &rpc.BuildRoutingTableRequest{TokenAddress: tokenAddr.Bytes()}
	url := fmt.Sprintf("http://%s/admin/route/build", adminHostPort)
	_, err := HttpPost(url, request)
	return err
}

// HTTP request to send the routing info to the listener/routing server.
func RecvRoutingInfo(adminHostPort string, info *rpc.RoutingRequest) error {
	url := fmt.Sprintf("http://%s/admin/route/recv_bcast", adminHostPort)
	_, err := HttpPost(url, info)
	return err
}

func HttpPost(url string, input interface{}) ([]byte, error) {
	log.Debugln("URL:>", url)
	payload, err := json.Marshal(input)
	if err != nil {
		return nil, fmt.Errorf("json.Marshal err: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("NewRequestWithContext err: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("DefaultClient.Do(req) err: %w", err)
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ReadAll(resp.Body) err: %w, status is %s", err, resp.Status)
	}

	if resp.StatusCode != http.StatusOK {
		return buf, fmt.Errorf("%w: %s", ErrHttpReponse, resp.Status)
	}

	return buf, nil
}

func getGrpcHttpErrMsg(resBody []byte) string {
	res := &errorBody{}
	if err := json.Unmarshal(resBody, res); err != nil {
		return ""
	}
	return res.Message
}

type errorBody struct {
	Error   string `json:"error"`
	Code    int32  `json:"code"`
	Message string `json:"message"`
}
