package main

import (
	"flag"

	"github.com/hashwavelab/osmoxy/gamm"
	"github.com/hashwavelab/osmoxy/proxy"
	"github.com/hashwavelab/osmoxy/wallet"
)

var (
	address    string
	accAddress string
	wallets    map[string]*wallet.Wallet = make(map[string]*wallet.Wallet)
)

func init() {
	a := flag.String("grpc", "localhost:9092", "gRPC server address")
	aa := flag.String("acc", "", "account address")
	flag.Parse()
	address = *a
	accAddress = *aa
}

func main() {
	p := proxy.NewProxy(address)
	p.InitBlockSubscription()

	dex := gamm.NewDex(p)
	dex.InitQueryingPoolsAfterEveryNewBlock()

	wallet := wallet.NewWallet(p, accAddress)
	wallet.InitQueryingBalancesAfterEveryNewBlock()
	wallets[accAddress] = wallet

	InitGrpcServer(p, dex, wallets)
}
