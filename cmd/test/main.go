package main

import (
	"flag"
	"log"

	"github.com/hashwavelab/osmoxy/proxy"
	"github.com/hashwavelab/osmoxy/tx"
	"github.com/hashwavelab/osmoxy/wallet"
)

var (
	address    string
	accAddress string
	hash       string
	wallets    map[string]*wallet.Wallet = make(map[string]*wallet.Wallet)
)

func init() {
	a := flag.String("grpc", "localhost:9092", "gRPC server address")
	aa := flag.String("acc", "", "account address")
	h := flag.String("hash", "", "tx hash")
	flag.Parse()
	address = *a
	accAddress = *aa
	hash = *h
}

func main() {
	p := proxy.NewProxy(address)
	p.InitBlockSubscription()

	r, err := tx.QuerySwapResultByHash(p, hash)
	log.Println(r, err)
}
