// Copyright 2018-2025 Celer Network

package e2e

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/celer-network/agent-pay/celersdkintf"
	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/common/structs"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/fsm"
	"github.com/celer-network/agent-pay/storage"
	tf "github.com/celer-network/agent-pay/testing"
	"github.com/celer-network/agent-pay/webapi"
	"github.com/celer-network/goutils/log"
	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	ethrpc "github.com/ethereum/go-ethereum/rpc"
)

// runtime variables, will be initialized by TestMain
var (
	// root dir with ending / for all files, outRootDirPrefix + epoch seconds
	// due to testframework etc in a different testing package, we have to define
	// same var in testframework.go and expose a set api
	outRootDir     string
	envDir         = "../../testing/env"
	noProxyProfile string // full file path to profile.json
	// erc20 token addr hex
	// map from app type to deployed addr, updated by SetupOnChain
	appAddrMap     = make(map[string]ctype.Addr)
	tokenAddrErc20 string // set by setupOnchain deploy erc20 contract

	tokenAddrNet1 string
	tokenAddrNet2 string
)

// toBuild map package subpath to binary file name eg. server -> server means build agent-pay/server and output server
var toBuild = map[string]string{
	"server":             "server",
	"testing/testclient": "testclient",
	"tools/osp-cli":      "ospcli",
}

// todo: remove addr arg
func getEthClient(addr string) (*ethclient.Client, error) {
	ws, err := ethrpc.Dial(ethGateway)
	if err != nil {
		return nil, err
	}
	conn := ethclient.NewClient(ws)
	return conn, nil
}

func sleep(second time.Duration) {
	time.Sleep(second * time.Second)
}

func payStatusFinalized(status int) bool {
	if status == celersdkintf.PAY_STATUS_PAID || status == celersdkintf.PAY_STATUS_PAID_RESOLVED_ONCHAIN ||
		status == celersdkintf.PAY_STATUS_UNPAID || status == celersdkintf.PAY_STATUS_UNPAID_EXPIRED ||
		status == celersdkintf.PAY_STATUS_UNPAID_REJECTED || status == celersdkintf.PAY_STATUS_UNPAID_DEST_UNREACHABLE {
		return true
	}
	return false
}

func waitForPaymentCompletion(payID string, sender, receiver *tf.ClientController) error {
	// Small grace period to allow message propagation
	time.Sleep(200 * time.Millisecond)
	// Extend retries to tolerate slower CI/macOS timing and background startup jitter (~30s total)
	const retryLimit = 20
	var status int
	var err error
	if sender != nil {
		var done bool
		for retry := 0; retry < retryLimit; retry++ {
			status, err = sender.GetOutgoingPaymentStatus(payID)
			if err != nil {
				return err
			}
			if payStatusFinalized(status) {
				done = true
				break
			}
			time.Sleep(400 * time.Millisecond)
		}
		if !done {
			return fmt.Errorf("payment not sent successfully, payID %s %s", payID, webapi.PayStatusName(status))
		}
	}
	if receiver != nil {
		for retry := 0; retry < retryLimit; retry++ {
			status, err = receiver.GetIncomingPaymentStatus(payID)
			if err != nil {
				return err
			}
			if payStatusFinalized(status) {
				return nil
			}
			time.Sleep(400 * time.Millisecond)
		}
		return fmt.Errorf("payment not received successfully, payID %s %s", payID, webapi.PayStatusName(status))
	}
	return nil
}

func waitForPaymentPending(payID string, sender, receiver *tf.ClientController) error {
	time.Sleep(200 * time.Millisecond)
	const retryLimit = 20
	var status int
	var err error
	if sender != nil {
		var done bool
		for retry := 0; retry < retryLimit; retry++ {
			status, err = sender.GetOutgoingPaymentStatus(payID)
			if err != nil {
				return err
			}
			if status == celersdkintf.PAY_STATUS_PENDING {
				done = true
				break
			}
			time.Sleep(400 * time.Millisecond)
		}
		if !done {
			return fmt.Errorf("payment not sent successfully, payID %s %s", payID, webapi.PayStatusName(status))
		}
	}
	if receiver != nil {
		for retry := 0; retry < retryLimit; retry++ {
			status, err = receiver.GetIncomingPaymentStatus(payID)
			if err != nil {
				return err
			}
			if status == celersdkintf.PAY_STATUS_PENDING {
				return nil
			}
			time.Sleep(400 * time.Millisecond)
		}
		return fmt.Errorf("payment not received successfully, payID %s %s", payID, webapi.PayStatusName(status))
	}
	return nil
}

func waitForPayStatus(payID string, sender, receiver *tf.ClientController, expStatus, timeoutSec int) error {
	time.Sleep(200 * time.Millisecond)
	deadline := time.Now().Add(time.Duration(timeoutSec) * time.Second)
	var status int
	var err error
	if sender != nil {
		var done bool
		for time.Now().Before(deadline) {
			status, err = sender.GetOutgoingPaymentStatus(payID)
			if err != nil {
				return err
			}
			if status == expStatus {
				done = true
				break
			}
			time.Sleep(400 * time.Millisecond)
		}
		if !done {
			return fmt.Errorf("payment not sent successfully, payID %s %s", payID, webapi.PayStatusName(status))
		}
	}
	if receiver != nil {
		for time.Now().Before(deadline) {
			status, err = receiver.GetIncomingPaymentStatus(payID)
			if err != nil {
				return err
			}
			if status == expStatus {
				return nil
			}
			time.Sleep(400 * time.Millisecond)
		}
		return fmt.Errorf("payment not received successfully, payID %s %s", payID, webapi.PayStatusName(status))
	}
	return nil
}

func checkOspPayState(dal *storage.DAL, payID string,
	expInCid ctype.CidType, expInState int, expOutCid ctype.CidType, expOutState, timeoutSec int) error {
	deadline := time.Now().Add(time.Duration(timeoutSec) * time.Second)
	var inCid, outCid ctype.CidType
	var inState, outState int
	var err error
	for time.Now().Before(deadline) {
		_, _, inCid, inState, outCid, outState, _, _, err = dal.GetPaymentInfo(ctype.Hex2PayID(payID))
		if err != nil {
			return err
		}
		if inCid == expInCid && inState == expInState && outCid == expOutCid && outState == expOutState {
			return nil
		}
		// value finalized but no match
		if (inCid != ctype.ZeroCid && inCid != expInCid) || (outCid != ctype.ZeroCid && outCid != expOutCid) ||
			(payStateFinalized(inState) && inState != expInState) || (payStateFinalized(outState) && outState != expOutState) {
			break
		}
		time.Sleep(400 * time.Millisecond)
	}
	return fmt.Errorf("payID %s, in: %x %s out: %x %s, expect in: %x %s out: %x %s",
		payID, inCid, fsm.PayStateName(inState), outCid, fsm.PayStateName(outState),
		expInCid, fsm.PayStateName(expInState), expOutCid, fsm.PayStateName(expOutState))
}

func payStateFinalized(state int) bool {
	if state == structs.PayState_COSIGNED_PAID || state == structs.PayState_COSIGNED_CANCELED {
		return true
	}
	return false
}

// save json as file path
func saveProfile(p *common.ProfileJSON, fpath string) {
	b, _ := json.Marshal(p)
	ioutil.WriteFile(fpath, b, 0644)
}

func SaveProfile(p *common.ProfileJSON, fpath string) {
	saveProfile(p, fpath)
}

// misc test configs from running, save for later reuse
type misc struct {
	GethPid int
	Erc20   string
	AppMap  map[string]ctype.Addr
}

func saveMisc(fpath string, pid int, erc string, m map[string]ctype.Addr) {
	s := &misc{
		GethPid: pid,
		Erc20:   erc,
		AppMap:  m,
	}
	b, _ := json.Marshal(s)
	ioutil.WriteFile(fpath, b, 0644)
}

func loadMisc(fpath string) *misc {
	ret := new(misc)
	b, _ := ioutil.ReadFile(fpath)
	json.Unmarshal(b, ret)
	return ret
}

// start process to handle eth rpc, and fund etherbase and server account
func StartChain() (*os.Process, error) {
	log.Infoln("outRootDir", outRootDir, "envDir", envDir)
	chainDataDir := outRootDir + "chaindata"
	logFname := outRootDir + "chain.log"
	if err := os.MkdirAll(chainDataDir, os.ModePerm); err != nil {
		return nil, err
	}

	cmdCopy := exec.Command("cp", "-a", "keystore", chainDataDir)
	cmdCopy.Dir, _ = filepath.Abs(envDir)
	log.Infoln("copy", filepath.Join(envDir, "keystore"), filepath.Join(chainDataDir, "keystore"))
	if err := cmdCopy.Run(); err != nil {
		return nil, err
	}

	// Allow overriding geth verbosity via env; default to 5 (debug)
	verbosity := os.Getenv("E2E_GETH_VERBOSITY")
	if verbosity == "" {
		verbosity = "5"
	}
	// Build geth args
	args := []string{
		"--verbosity", verbosity,
		"--dev",
		"--dev.period", "1",
		"--cache", "256",
		"--datadir", chainDataDir,
		"--keystore", filepath.Join(chainDataDir, "keystore"),
		"--miner.etherbase", "0x" + etherBaseAddr,
		// HTTP RPC
		"--http",
		"--http.addr", "127.0.0.1",
		"--http.port", "8545",
		"--http.corsdomain", "*",
		"--http.api", "admin,debug,eth,miner,net,txpool,web3",
		// WebSocket RPC
		"--ws",
		"--ws.addr", "127.0.0.1",
		"--ws.port", "8546",
		"--ws.api", "admin,debug,eth,miner,net,txpool,web3",
		"--ws.origins", "*",
		// Global RPC toggles
		"--rpc.allow-unprotected-txs",
	}
	// Run geth in --dev mode with fast empty-block sealing
	cmd := exec.Command("geth", args...)
	// set working dir so relative files resolve
	cmd.Dir, _ = filepath.Abs(envDir)

	logF, _ := os.Create(logFname)
	// By default, write geth logs only to file; enable stdout with E2E_GETH_LOG_STDOUT=1
	if os.Getenv("E2E_GETH_LOG_STDOUT") == "1" || os.Getenv("E2E_GETH_LOG_STDOUT") == "true" {
		mw := io.MultiWriter(logF, os.Stdout)
		cmd.Stderr = mw
		cmd.Stdout = mw
	} else {
		cmd.Stderr = logF
		cmd.Stdout = logF
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	fmt.Println("geth pid:", cmd.Process.Pid)
	// Wait for chain to start producing blocks and fund etherbase via mining in dev mode
	// so subsequent transactions signed by etherBasePriv have balance.
	go func() {
		// best-effort readiness probe; don't fail tests here
		// Wait until HTTP is reachable
		for i := 0; i < 50; i++ { // ~10s
			if _, err := os.Stat(logFname); err == nil {
				// small delay for RPC to come up
				time.Sleep(100 * time.Millisecond)
			}
			// try dialing
			if _, err := ethrpc.Dial(ethGateway); err == nil {
				break
			}
			time.Sleep(200 * time.Millisecond)
		}
		// check height and balance
		if conn, err := getEthClient(ethGateway); err == nil {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			// wait for at least block 1
			for i := 0; i < 50; i++ { // ~10s
				hdr, e := conn.HeaderByNumber(ctx, nil)
				if e == nil && hdr != nil && hdr.Number != nil && hdr.Number.Uint64() >= 1 {
					break
				}
				time.Sleep(200 * time.Millisecond)
			}
			// wait until etherbase has non-zero balance
			eb := ethcommon.HexToAddress("0x" + etherBaseAddr)
			for i := 0; i < 50; i++ { // ~10s
				bal, e := conn.BalanceAt(ctx, eb, nil)
				if e == nil && bal != nil && bal.Cmp(big.NewInt(0)) > 0 {
					break
				}
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()
	// in case geth exits with non-zero, exit test early
	// if geth is killed by ethProc.Signal, it exits w/ 0
	go func() {
		if err := cmd.Wait(); err != nil {
			fmt.Println("geth process failed:", err)
			os.Exit(1)
		}
	}()
	return cmd.Process, nil
}

const agentPayRepo = "github.com/celer-network/agent-pay/"

func buildBins(rootDir string) error {
	for pkg, bin := range toBuild {
		err := buildPkgBin(rootDir, pkg, bin)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildPkgBin(rootDir, pkg, bin string) error {
	fmt.Println("Building", pkg, "->", bin)
	cmd := exec.Command("go", "build", "-o", rootDir+bin, agentPayRepo+pkg)
	cmd.Stderr, _ = os.OpenFile(rootDir+"build.err", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func chkErr(e error, msg string) {
	if e != nil {
		fmt.Println("Err:", msg, e)
		os.Exit(1)
	}
}

func CheckError(e error, msg string) {
	chkErr(e, msg)
}

func SetEnvDir(dir string) {
	envDir = dir
}

func SetOutRootDir(dir string) {
	outRootDir = dir
}
