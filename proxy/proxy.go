package proxy

import (
	"log"
	"sync"
	"time"

	base "github.com/cosmos/cosmos-sdk/api/cosmos/base/tendermint/v1beta1"
	"github.com/hashwavelab/osmoxy/broadcast"
	"google.golang.org/grpc"
)

var (
	BlockQueryInterval time.Duration = 1 * time.Second
	QueryTimeOut       float64       = 5
)

type Proxy struct {
	sync.RWMutex
	conn                *grpc.ClientConn
	address             string
	blockChan           chan *base.GetLatestBlockResponse
	blockNumber         uint64
	newBlockBoardcaster broadcast.Broadcaster
}

func NewProxy(a string) *Proxy {
	p := &Proxy{
		address:             a,
		blockChan:           make(chan *base.GetLatestBlockResponse),
		newBlockBoardcaster: broadcast.NewBroadcaster(),
	}

	// setup the grpc connection
	conn, err := p.connect()
	if err != nil {
		log.Fatal(err)
	}
	p.conn = conn

	p.initBlockProcessor()
	return p
}

func (p *Proxy) initBlockProcessor() {
	go func() {
		for b := range p.blockChan {
			bn := b.Block.LastCommit.Height
			p.setLastBlock(b, uint64(bn))
			p.newBlockBoardcaster.Submit(b)
		}
	}()
}

func (p *Proxy) setLastBlock(b *base.GetLatestBlockResponse, bn uint64) {
	p.Lock()
	defer p.Unlock()
	// TODO: save the last block
	p.blockNumber = bn
}

func (p *Proxy) GetLastBlockNumber() uint64 {
	p.RLock()
	defer p.RUnlock()
	return p.blockNumber
}

func (p *Proxy) InitBlockSubscription() {
	log.Println("try connecting to:", p.address)
	var lastBlockNumber int64 = 0
	go func() {
		for {
			func() {
				defer time.Sleep(BlockQueryInterval)
				resp, err := p.GetLatestBlock()
				if err != nil {
					return
				}
				bn := resp.Block.LastCommit.Height
				if bn > lastBlockNumber {
					log.Println("New Block", bn)
					p.blockChan <- resp
					lastBlockNumber = bn
				}
			}()
		}
	}()
}

func (p *Proxy) SubscribeNewBlock(c chan interface{}, buffer int) {
	p.newBlockBoardcaster.Register(c, buffer)
}

func (p *Proxy) UnsubscribeNewBlock(c chan interface{}) {
	p.newBlockBoardcaster.Unregister(c)
}
