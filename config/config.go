// Copyright 2018-2025 Celer Network

package config

import (
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/rtconfig"
	"github.com/celer-network/goutils/eth"
	"google.golang.org/grpc/keepalive"
)

// envUint returns the env var parsed as a uint64, or def if unset / unparsable.
// Used for the "safe margin" knobs below so e2e tests can shrink production-safe
// 60-second slacks down to a few seconds without recompiling the server binary.
func envUint(key string, def uint64) uint64 {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.ParseUint(v, 10, 64); err == nil {
			return n
		}
	}
	return def
}

// NOTE: not protected by lock, only set once at initialization.
// All "*Timeout" / "*Deadline" / "*SafeMargin" values below are in seconds —
// the contracts measure deadlines and challenge windows in `block.timestamp`
// (unix seconds), so off-chain code follows the same unit and comparisons
// against on-chain state stay aligned.
var (
	ChainId               *big.Int
	ChannelDisputeTimeout = uint64(86400) // 1 day, seconds
	BlockDelay            = uint64(5)
	BlockIntervalSec      = uint64(10)
	EventListenerHttp     = ""
	RouterBcastInterval   = 293 * time.Second
	RouterBuildInterval   = 367 * time.Second
	RouterAliveTimeout    = 900 * time.Second
	OspClearPaysInterval  = 613 * time.Second
	OspReportInverval     = 887 * time.Second

	// Safe-margin knobs are env-var tunable so e2e tests can shrink them. Production
	// defaults (60s) absorb chain-confirmation slack past a deadline; tests typically
	// set AGENTPAY_*_SAFE_MARGIN_S=5 to keep the timeout-and-sweep flow snappy.
	WithdrawTimeoutSafeMargin = envUint("AGENTPAY_WITHDRAW_SAFE_MARGIN_S", 60) // seconds
	PaySendTimeoutSafeMargin  = envUint("AGENTPAY_SEND_SAFE_MARGIN_S", 60)     // seconds
	PayRecvTimeoutSafeMargin  = envUint("AGENTPAY_RECV_SAFE_MARGIN_S", 60)     // seconds
)

const (
	ClientCacheSize            = 1000
	ServerCacheSize            = 16
	OpenChannelTimeout         = uint64(600)    // seconds
	CooperativeWithdrawTimeout = uint64(60)     // seconds
	PayResolveTimeout          = uint64(60)     // seconds (on-chain partial-resolve challenge window)
	AdminSendTokenTimeout      = uint64(600)    // seconds
	QuickCatchBlockDelay       = uint64(2)      // blocks (unrelated to deadlines — fast-path reorg confirmation)
	TcbTimeoutSeconds          = uint64(604800) // 7 days, seconds

	// Protocol Version in AuthReq, >=1 support sync
	AuthProtocolVersion = uint64(1)
	// AuthAckTimeout is duration client will wait for AuthAck msg
	AuthAckTimeout = 5 * time.Second

	// grpc dial timeout seconds, block until 15s
	GrpcDialTimeout = 15

	EventListenerLeaseName          = "eventlistener"
	EventListenerLeaseRenewInterval = 60 * time.Second
	EventListenerLeaseTimeout       = 90 * time.Second

	// used by clients to control onchain query frequency
	QueryName_OnChainBalance      = "onchainBalance"
	QueryName_OnChainResolvedPays = "onchainResolvedPays"
)

// KeepAliveClientParams is grpc client side keeyalive parameters
// Make sure these parameters are set in coordination with the keepalive policy
// on the server, as incompatible settings can result in closing of connection
var KeepAliveClientParams = keepalive.ClientParameters{
	Time:                15 * time.Second, // send pings every 15 seconds if there is no activity
	Timeout:             3 * time.Second,  // wait 3 seconds for ping ack before considering the connection dead
	PermitWithoutStream: true,             // send pings even without active streams
}

// KeepAliveServerParams is grpc server side keeyalive parameters
var KeepAliveServerParams = keepalive.ServerParameters{
	Time:    20 * time.Second, // send pings every 20 seconds if there is no activity
	Timeout: 3 * time.Second,  // wait 3 seconds for ping ack before considering the connection dead
}

// KeepAliveEnforcePolicy is grpc server side policy
var KeepAliveEnforcePolicy = keepalive.EnforcementPolicy{
	MinTime:             12 * time.Second, // must be smaller than clientParam.Time
	PermitWithoutStream: true,
}

func SetGlobalConfigFromProfile(profile *common.CProfile) {
	ChainId = big.NewInt(profile.ChainId)
	BlockDelay = profile.BlockDelayNum
	if profile.PollingInterval != 0 {
		BlockIntervalSec = profile.PollingInterval
	}
	if profile.DisputeTimeout != 0 {
		ChannelDisputeTimeout = profile.DisputeTimeout
	}
}

func WaitMinedOptions() []eth.TxOption {
	return []eth.TxOption{
		eth.WithBlockDelay(BlockDelay),
		eth.WithPollingInterval(time.Duration(BlockIntervalSec) * time.Second),
		eth.WithTimeout(time.Duration(rtconfig.GetWaitMinedTxTimeout()) * time.Second),
		eth.WithQueryTimeout(time.Duration(rtconfig.GetWaitMinedTxQueryTimeout()) * time.Second),
		eth.WithQueryRetryInterval(time.Duration(rtconfig.GetWaitMinedTxQueryRetryInterval()) * time.Second),
	}
}

func TransactOptions(opts ...eth.TxOption) []eth.TxOption {
	options := []eth.TxOption{
		eth.WithMinGasGwei(float64(rtconfig.GetMinGasGwei())),
		eth.WithMaxGasGwei(float64(rtconfig.GetMaxGasGwei())),
		eth.WithAddGasGwei(float64(rtconfig.GetAddGasGwei())),
	}
	options = append(options, WaitMinedOptions()...)
	return append(options, opts...)
}

func QuickTransactOptions(opts ...eth.TxOption) []eth.TxOption {
	options := TransactOptions(opts...)
	if QuickCatchBlockDelay < BlockDelay {
		// this will overwrite the previous WithBlockDelay option
		options = append(options, eth.WithBlockDelay(QuickCatchBlockDelay))
	}
	return options
}
