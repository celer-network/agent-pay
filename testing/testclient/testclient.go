// Copyright 2018-2025 Celer Network

package main

import (
	"flag"
	"io/ioutil"
	"math/big"

	"github.com/celer-network/agent-pay/celersdk"
	"github.com/celer-network/agent-pay/common"
	"github.com/celer-network/agent-pay/ctype"
	"github.com/celer-network/agent-pay/webapi"
	"github.com/celer-network/goutils/eth"
	"github.com/celer-network/goutils/log"
)

var (
	grpcPort  = flag.Int("port", -1, "gRPC server listening port")
	ksPath    = flag.String("keystore", "", "Path to keystore json file")
	cfgPath   = flag.String("config", "profile.json", "Path to config json file")
	dataPath  = flag.String("datadir", "", "Path to the local database")
	extSigner = flag.Bool("extsign", false, "if set, exercise the external signer interface")
)

func main() {
	flag.Parse()
	ksBytes, err := ioutil.ReadFile(*ksPath)
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := ioutil.ReadFile(*cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Infoln("start testclient on port", *grpcPort, "using ks", *ksPath)
	if *extSigner {
		addr, priv, err := eth.GetAddrPrivKeyFromKeystore(string(ksBytes), "")
		if err != nil {
			log.Fatal(err)
		}
		p := common.Bytes2Profile(cfg)
		signer, err := eth.NewSigner(priv, big.NewInt(p.ChainId))
		if err != nil {
			log.Fatal(err)
		}
		webapi.NewInternalApiServerWithExternalSigner(
			-1,
			*grpcPort,
			"http://localhost:*",
			ctype.Addr2Hex(addr),
			*dataPath,
			string(cfg[:]),
			&testExternalSigner{Signer: signer},
			nil).Start()
		return
	}
	webapi.NewInternalApiServer(
		-1,
		*grpcPort,
		"http://localhost:*",
		string(ksBytes[:]),
		"",
		*dataPath,
		string(cfg[:])).Start()
}

type testExternalSigner struct {
	eth.Signer
}

func (es *testExternalSigner) OnSignMessage(reqid int, msg []byte) {
	res, _ := es.SignEthMessage(msg)
	celersdk.PublishSignedResult(reqid, res)
}

func (es *testExternalSigner) OnSignTransaction(reqid int, rawtx []byte) {
	res, _ := es.SignEthTransaction(rawtx)
	celersdk.PublishSignedResult(reqid, res)
}
