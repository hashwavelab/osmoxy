package main

import (
	"flag"

	"github.com/hashwavelab/osmoxy/gamm"
	"github.com/hashwavelab/osmoxy/proxy"
)

var address string

func init() {
	a := flag.String("a", "localhost:9092", "gRPC server address")
	flag.Parse()
	address = *a
}

func main() {
	p := proxy.NewProxy(address)
	p.InitBlockSubscription()

	dex := gamm.NewDex(p)
	ch := make(chan interface{})
	p.SubscribeNewBlock(ch, 0)
	dex.InitQueryingPoolsAfterEveryNewBlock(ch)

	InitGrpcServer(p, dex)
}
