package proxy

import (
	"log"
	"sync"
	"time"

	base "github.com/cosmos/cosmos-sdk/api/cosmos/base/tendermint/v1beta1"
	"github.com/hashwavelab/osmoxy/broadcast"
)

var (
	BlockQueryInterval time.Duration = 1 * time.Second
	QueryTimeOut       float64       = 5
)

type Proxy struct {
	sync.RWMutex
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
	p.initBlockProcessor()
	return p
}

func (_p *Proxy) initBlockProcessor() {
	go func() {
		for b := range _p.blockChan {
			bn := b.Block.LastCommit.Height
			_p.setLastBlock(b, uint64(bn))
			_p.newBlockBoardcaster.Submit(b)
		}
	}()
}

func (_p *Proxy) setLastBlock(b *base.GetLatestBlockResponse, bn uint64) {
	_p.Lock()
	defer _p.Unlock()
	// TODO: save the last block
	_p.blockNumber = bn
}

func (_p *Proxy) GetLastBlockNumber() uint64 {
	_p.RLock()
	defer _p.RUnlock()
	return _p.blockNumber
}

func (_p *Proxy) InitBlockSubscription() {
	log.Println("try connecting to:", _p.address)
	var lastBlockNumber int64 = 0
	go func() {
		for {
			func() {
				defer time.Sleep(BlockQueryInterval)
				resp, err := _p.GetLatestBlock()
				if err != nil {
					return
				}
				bn := resp.Block.LastCommit.Height
				if bn > lastBlockNumber {
					log.Println("New Block", bn)
					_p.blockChan <- resp
					lastBlockNumber = bn
				}
			}()
		}
	}()
}

func (_p *Proxy) SubscribeNewBlock(c chan interface{}, buffer int) {
	_p.newBlockBoardcaster.Register(c, buffer)
}

func (_p *Proxy) UnsubscribeNewBlock(c chan interface{}) {
	_p.newBlockBoardcaster.Unregister(c)
}
