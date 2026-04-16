// Copyright 2018-2025 Celer Network

package e2e

import (
	"os"
	"testing"
	"time"

	"github.com/celer-network/agent-pay/ctype"
	tf "github.com/celer-network/agent-pay/testing"
)

func setMultiOSP() ([]*tf.ServerController, error) {
	// osp1 is already registered in TestMain; register the additional multi-OSP routers here.
	if err := tf.RegisterRouters([]string{osp2Keystore, osp3Keystore, osp4Keystore, osp5Keystore}); err != nil {
		return nil, err
	}
	os.RemoveAll(sStoreDir)
	// Be careful: because the limit of test set up, state between the two osps is kept during tests.
	// Each tests need to reset state between the two osps.
	o1 := tf.StartServerController(outRootDir+toBuild["server"],
		"-profile", noProxyProfile,
		"-port", o1Port,
		"-storedir", sStoreDir,
		"-ks", ospKeystore,
		"-nopassword",
		"-rtc", rtConfigMultiOSP,
		"-defaultroute", osp3EthAddr,
		"-svrname", "o1",
		"-logcolor",
		"-logprefix", "o1_"+ospEthAddr[:4])

	o2 := tf.StartServerController(outRootDir+toBuild["server"],
		"-profile", noProxyProfile,
		"-port", o2Port,
		"-storedir", sStoreDir,
		"-adminrpc", o2AdminRPC,
		"-adminweb", o2AdminWeb,
		"-ks", osp2Keystore,
		"-nopassword",
		"-rtc", rtConfigMultiOSP,
		"-svrname", "o2",
		"-logcolor",
		"-logprefix", "o2_"+osp2EthAddr[:4])

	o3 := tf.StartServerController(outRootDir+toBuild["server"],
		"-profile", noProxyProfile,
		"-port", o3Port,
		"-storedir", sStoreDir,
		"-adminrpc", o3AdminRPC,
		"-adminweb", o3AdminWeb,
		"-ks", osp3Keystore,
		"-nopassword",
		"-rtc", rtConfigMultiOSP,
		"-defaultroute", osp4EthAddr,
		"-svrname", "o3",
		"-logcolor",
		"-logprefix", "o3_"+osp3EthAddr[:4])

	o4 := tf.StartServerController(outRootDir+toBuild["server"],
		"-profile", noProxyProfile,
		"-port", o4Port,
		"-storedir", sStoreDir,
		"-adminrpc", o4AdminRPC,
		"-adminweb", o4AdminWeb,
		"-ks", osp4Keystore,
		"-nopassword",
		"-rtc", rtConfigMultiOSP,
		"-defaultroute", ospEthAddr,
		"-svrname", "o4",
		"-logcolor",
		"-logprefix", "o4_"+osp4EthAddr[:4])

	o5 := tf.StartServerController(outRootDir+toBuild["server"],
		"-profile", noProxyProfile,
		"-port", o5Port,
		"-storedir", sStoreDir,
		"-adminrpc", o5AdminRPC,
		"-adminweb", o5AdminWeb,
		"-ks", osp5Keystore,
		"-nopassword",
		"-rtc", rtConfigMultiOSP,
		"-defaultroute", osp4EthAddr,
		"-svrname", "o5",
		"-logcolor",
		"-logprefix", "o5_"+osp5EthAddr[:4])

	time.Sleep(3 * time.Second)
	/* OSPs connect with each other with topology:
	  o4---o5
	 /  \    \
	o3---o1---o2
	*/
	if err := registerStreamWithRetry(o2AdminWeb, ctype.Hex2Addr(ospEthAddr), localhost+o1Port); err != nil {
		return nil, err
	}
	if err := registerStreamWithRetry(o3AdminWeb, ctype.Hex2Addr(ospEthAddr), localhost+o1Port); err != nil {
		return nil, err
	}
	if err := registerStreamWithRetry(o4AdminWeb, ctype.Hex2Addr(ospEthAddr), localhost+o1Port); err != nil {
		return nil, err
	}
	if err := registerStreamWithRetry(o4AdminWeb, ctype.Hex2Addr(osp3EthAddr), localhost+o3Port); err != nil {
		return nil, err
	}
	if err := registerStreamWithRetry(o4AdminWeb, ctype.Hex2Addr(osp5EthAddr), localhost+o5Port); err != nil {
		return nil, err
	}
	if err := registerStreamWithRetry(o5AdminWeb, ctype.Hex2Addr(osp2EthAddr), localhost+o2Port); err != nil {
		return nil, err
	}
	time.Sleep(4 * time.Second)

	return []*tf.ServerController{o1, o2, o3, o4, o5}, nil
}

func TestE2EMultiOSP(t *testing.T) {
	svrs, err := setMultiOSP()
	if err != nil {
		t.Fatal(err)
	}
	defer tearDownMultiSvr([]Killable{svrs[0], svrs[1], svrs[2], svrs[3], svrs[4]})

	// Be careful: because the limit of test set up, state between the two osps is kept during tests.
	// Each tests need to reset state between the two osps.
	t.Run("e2e-multiosp-open-channel", func(t *testing.T) {
		t.Run("multiOspOpenChannelPolicyTest", multiOspOpenChannelPolicyTest)
		t.Run("multiOspOpenChannelPolicyFallbackTest", multiOspOpenChannelPolicyFallbackTest)
		t.Run("multiOspOpenChannelTest", multiOspOpenChannelTest)
	})
	t.Run("e2e-multiosp-routing", func(t *testing.T) {
		t.Run("multiOspRouting", multiOspRouting(svrs...))
	})
	t.Run("e2e-multiosp-channel-migration", func(t *testing.T) {
		t.Run("migrateChannelBetweenOsps", migrateChannelBetweenOsps(svrs...))
	})
}
