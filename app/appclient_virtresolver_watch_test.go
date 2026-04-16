package app

import (
	"errors"
	"math/big"
	"sync"
	"testing"

	"github.com/celer-network/agent-pay/chain"
	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/config"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/goutils/eth/monitor"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type testStateCallback struct{}

func (testStateCallback) OnDispute(int) {}

type fakeMonitorService struct {
	mu          sync.Mutex
	monitorHits int
	removed     []monitor.CallbackID
}

func (m *fakeMonitorService) GetCurrentBlockNumber() *big.Int { return big.NewInt(1) }
func (m *fakeMonitorService) RegisterDeadline(monitor.Deadline) monitor.CallbackID {
	return 0
}
func (m *fakeMonitorService) Monitor(_ *monitor.Config, _ func(monitor.CallbackID, types.Log) bool) (monitor.CallbackID, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.monitorHits++
	return monitor.CallbackID(123), nil
}
func (m *fakeMonitorService) RemoveDeadline(monitor.CallbackID) {}
func (m *fakeMonitorService) RemoveEvent(id monitor.CallbackID) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.removed = append(m.removed, id)
}
func (m *fakeMonitorService) Close() {}

type fakeContract struct{}

func (fakeContract) GetAddr() ctype.Addr             { return ctype.ZeroAddr }
func (fakeContract) GetABI() string                  { return "" }
func (fakeContract) GetETHClient() *ethclient.Client { return nil }
func (fakeContract) SendTransaction(*bind.TransactOpts, string, ...interface{}) (*types.Transaction, error) {
	return nil, errors.New("not implemented")
}
func (fakeContract) CallFunc(interface{}, string, ...interface{}) error {
	return errors.New("not implemented")
}
func (fakeContract) WatchEvent(string, *bind.WatchOpts, <-chan bool) (types.Log, error) {
	return types.Log{}, errors.New("not implemented")
}
func (fakeContract) FilterEvent(string, *bind.FilterOpts, interface{}) (*chain.EventIterator, error) {
	return nil, errors.New("not implemented")
}
func (fakeContract) ParseEvent(string, types.Log, interface{}) error {
	return errors.New("not implemented")
}

type fakeNodeConfig struct{ virt chain.Contract }

func (f fakeNodeConfig) GetOnChainAddr() ctype.Addr                           { return ctype.ZeroAddr }
func (f fakeNodeConfig) GetEthPoolAddr() ctype.Addr                           { return ctype.ZeroAddr }
func (f fakeNodeConfig) GetEthConn() *ethclient.Client                        { return nil }
func (f fakeNodeConfig) GetRPCAddr() string                                   { return "" }
func (f fakeNodeConfig) GetSvrName() string                                   { return "" }
func (f fakeNodeConfig) GetWalletContract() chain.Contract                    { return nil }
func (f fakeNodeConfig) GetLedgerContract() chain.Contract                    { return nil }
func (f fakeNodeConfig) GetLedgerContractOn(ctype.Addr) chain.Contract        { return nil }
func (f fakeNodeConfig) GetAllLedgerContracts() map[ctype.Addr]chain.Contract { return nil }
func (f fakeNodeConfig) GetLedgerContractOf(ctype.CidType) chain.Contract     { return nil }
func (f fakeNodeConfig) GetVirtResolverContract() chain.Contract              { return f.virt }
func (f fakeNodeConfig) GetPayResolverContract() chain.Contract               { return nil }
func (f fakeNodeConfig) GetPayRegistryContract() chain.Contract               { return nil }
func (f fakeNodeConfig) GetRouterRegistryContract() chain.Contract            { return nil }
func (f fakeNodeConfig) GetCheckInterval(string) uint64                       { return 0 }

func TestNewAppChannelOnVirtualContract_SharedDeployWatch(t *testing.T) {
	oldChainID := config.ChainId
	config.ChainId = big.NewInt(1)
	t.Cleanup(func() { config.ChainId = oldChainID })

	ms := &fakeMonitorService{}
	nc := fakeNodeConfig{virt: fakeContract{}}

	c := NewAppClient(nc, nil, nil, ms, nil, nil)
	cb := common.StateCallback(testStateCallback{})

	cid1, err := c.NewAppChannelOnVirtualContract([]byte{0x01}, []byte{0x02}, 1, 0, cb)
	if err != nil {
		t.Fatalf("first NewAppChannelOnVirtualContract failed: %v", err)
	}
	cid2, err := c.NewAppChannelOnVirtualContract([]byte{0x01}, []byte{0x02}, 2, 0, cb)
	if err != nil {
		t.Fatalf("second NewAppChannelOnVirtualContract failed: %v", err)
	}
	if cid1 == cid2 {
		t.Fatalf("expected different cids, got both %s", cid1)
	}

	ms.mu.Lock()
	gotMonitorHits := ms.monitorHits
	ms.mu.Unlock()
	if gotMonitorHits != 1 {
		t.Fatalf("expected 1 Monitor() call, got %d", gotMonitorHits)
	}

	c.DeleteAppChannel(cid1)
	c.DeleteAppChannel(cid2)

	ms.mu.Lock()
	removed := append([]monitor.CallbackID(nil), ms.removed...)
	ms.mu.Unlock()
	if len(removed) != 1 {
		t.Fatalf("expected shared watch removed once, got %d removals: %v", len(removed), removed)
	}
}
