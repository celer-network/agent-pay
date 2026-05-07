package cooperativewithdraw

import (
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/celer-network/agent-pay/chain"
	ledgerbinding "github.com/celer-network/agent-pay/chain/channel-eth-go/ledger"
	"github.com/celer-network/agent-pay/common/structs"
	"github.com/celer-network/agent-pay/config"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/rpc"
	"github.com/celer-network/agent-pay/storage"
	"github.com/celer-network/agent-pay/utils"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type testContract struct{ addr ctype.Addr }

func (c testContract) GetAddr() ctype.Addr           { return c.addr }
func (testContract) GetABI() string                  { return "" }
func (testContract) GetETHClient() *ethclient.Client { return nil }
func (testContract) SendTransaction(*bind.TransactOpts, string, ...interface{}) (*types.Transaction, error) {
	return nil, nil
}
func (testContract) CallFunc(interface{}, string, ...interface{}) error { return nil }
func (testContract) WatchEvent(string, *bind.WatchOpts, <-chan bool) (types.Log, error) {
	return types.Log{}, nil
}
func (testContract) FilterEvent(string, *bind.FilterOpts, interface{}) (*chain.EventIterator, error) {
	return nil, nil
}
func (testContract) ParseEvent(string, types.Log, interface{}) error { return nil }

type testNodeConfig struct {
	self   ctype.Addr
	ledger chain.Contract
}

func (c testNodeConfig) GetOnChainAddr() ctype.Addr                    { return c.self }
func (testNodeConfig) GetNativeWrapAddr() ctype.Addr                      { return ctype.ZeroAddr }
func (testNodeConfig) GetEthConn() *ethclient.Client                   { return nil }
func (testNodeConfig) GetRPCAddr() string                              { return "" }
func (testNodeConfig) GetSvrName() string                              { return "" }
func (testNodeConfig) GetWalletContract() chain.Contract               { return nil }
func (c testNodeConfig) GetLedgerContract() chain.Contract             { return c.ledger }
func (c testNodeConfig) GetLedgerContractOn(ctype.Addr) chain.Contract { return c.ledger }
func (c testNodeConfig) GetAllLedgerContracts() map[ctype.Addr]chain.Contract {
	return map[ctype.Addr]chain.Contract{c.ledger.GetAddr(): c.ledger}
}
func (c testNodeConfig) GetLedgerContractOf(ctype.CidType) chain.Contract { return c.ledger }
func (testNodeConfig) GetVirtResolverContract() chain.Contract            { return nil }
func (testNodeConfig) GetPayResolverContract() chain.Contract             { return nil }
func (testNodeConfig) GetPayRegistryContract() chain.Contract             { return nil }
func (testNodeConfig) GetRouterRegistryContract() chain.Contract          { return nil }
func (testNodeConfig) GetCheckInterval(string) uint64                     { return 0 }

type testWithdrawCallback struct {
	txHash chan string
	err    chan string
}

func (cb *testWithdrawCallback) OnWithdraw(_ string, txHash string) {
	cb.txHash <- txHash
}

func (cb *testWithdrawCallback) OnError(_ string, err string) {
	cb.err <- err
}

func newTestDAL(t *testing.T) *storage.DAL {
	t.Helper()
	storeFile, err := os.CreateTemp("", "coop-withdraw-*.db")
	if err != nil {
		t.Fatalf("CreateTemp() err = %v", err)
	}
	storePath := storeFile.Name()
	if err := storeFile.Close(); err != nil {
		t.Fatalf("Close() err = %v", err)
	}
	os.Remove(storePath)
	t.Cleanup(func() { os.Remove(storePath) })

	st, err := storage.NewKVStoreSQL("sqlite3", storePath)
	if err != nil {
		t.Fatalf("NewKVStoreSQL() err = %v", err)
	}
	t.Cleanup(func() { st.Close() })
	return storage.NewDAL(st)
}

func TestUpdateOnChainBalanceAdvancesJobForOsp(t *testing.T) {
	oldChainID := config.ChainId
	config.ChainId = big.NewInt(1)
	t.Cleanup(func() { config.ChainId = oldChainID })

	self := ctype.Hex2Addr("1111111111111111111111111111111111111111")
	peer := ctype.Hex2Addr("2222222222222222222222222222222222222222")
	ledgerAddr := ctype.Hex2Addr("3333333333333333333333333333333333333333")
	cid := ctype.Hex2Cid("abc123")
	dal := newTestDAL(t)

	token := utils.GetTokenInfoFromAddress(ctype.NativeTokenAddr)
	openResp := &rpc.OpenChannelResponse{}
	onChainBalance := &structs.OnChainBalance{}
	simplex := &rpc.SignedSimplexState{}
	if err := dal.InsertChan(cid, peer, token, ledgerAddr, 1, openResp, onChainBalance, 0, 0, 0, 0, simplex, simplex); err != nil {
		t.Fatalf("InsertChan() err = %v", err)
	}

	p := &Processor{
		nodeConfig:  testNodeConfig{self: self, ledger: testContract{addr: ledgerAddr}},
		selfAddress: self,
		dal:         dal,
		callbacks:   make(map[string]Callback),
		runningJobs: make(map[string]bool),
		keepMonitor: true,
		enableJobs:  false,
	}

	withdrawHash := p.generateWithdrawHash(cid, ledgerAddr, 7)
	job := &structs.CooperativeWithdrawJob{
		WithdrawHash: withdrawHash,
		State:        structs.CooperativeWithdrawWaitTx,
		TxHash:       "0xtesttx",
		LedgerAddr:   ledgerAddr,
	}
	if err := dal.PutCooperativeWithdrawJob(withdrawHash, job); err != nil {
		t.Fatalf("PutCooperativeWithdrawJob() err = %v", err)
	}

	cb := &testWithdrawCallback{
		txHash: make(chan string, 1),
		err:    make(chan string, 1),
	}
	p.registerCallback(withdrawHash, cb)

	event := &ledgerbinding.CelerLedgerCooperativeWithdraw{
		ChannelId:       cid,
		Receiver:        self,
		Deposits:        [2]*big.Int{big.NewInt(5), big.NewInt(5)},
		Withdrawals:     [2]*big.Int{big.NewInt(1), big.NewInt(0)},
		SeqNum:          big.NewInt(7),
		WithdrawnAmount: big.NewInt(1),
	}

	p.updateOnChainBalance(cid, self, peer, event, "0xtesttx")

	select {
	case txHash := <-cb.txHash:
		if txHash != "0xtesttx" {
			t.Fatalf("OnWithdraw() txHash = %q, want 0xtesttx", txHash)
		}
	case errMsg := <-cb.err:
		t.Fatalf("OnError() err = %q, want withdraw success", errMsg)
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for withdraw callback")
	}

	updatedJob, err := dal.GetCooperativeWithdrawJob(withdrawHash)
	if err != nil {
		t.Fatalf("GetCooperativeWithdrawJob() err = %v", err)
	}
	if updatedJob.State != structs.CooperativeWithdrawSucceeded {
		t.Fatalf("job state = %d, want %d", updatedJob.State, structs.CooperativeWithdrawSucceeded)
	}

	updatedBalance, found, err := dal.GetOnChainBalance(cid)
	if err != nil {
		t.Fatalf("GetOnChainBalance() err = %v", err)
	}
	if !found {
		t.Fatal("GetOnChainBalance() did not find updated balance")
	}
	if updatedBalance.MyWithdrawal == nil || updatedBalance.MyWithdrawal.Cmp(big.NewInt(1)) != 0 {
		t.Fatalf("MyWithdrawal = %v, want 1", updatedBalance.MyWithdrawal)
	}
}
