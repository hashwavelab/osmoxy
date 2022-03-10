package proxy

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/hashwavelab/osmoxy/broadcast"
)

var (
	BlockQueryInterval time.Duration = 1 * time.Second
	QueryTimeOut       float64       = 5
)

type Proxy struct {
	sync.RWMutex
	address             string
	blockChan           chan *simplejson.Json
	blockNumber         uint64
	newBlockBoardcaster broadcast.Broadcaster
}

func NewProxy(a string) *Proxy {
	p := &Proxy{
		address:             a,
		blockChan:           make(chan *simplejson.Json),
		newBlockBoardcaster: broadcast.NewBroadcaster(),
	}
	p.initBlockProcessor()
	return p
}

func (_p *Proxy) initBlockProcessor() {
	lastBlockNumber := 0
	go func() {
		for b := range _p.blockChan {
			bn, err := strconv.Atoi(b.Get("block").Get("lastCommit").Get("height").MustString())
			if err != nil {
				continue
			}
			if bn != lastBlockNumber {
				_p.setLastBlock(b, uint64(bn))
				_p.newBlockBoardcaster.Submit(b)
				lastBlockNumber = bn
			}
		}
	}()
}

func (_p *Proxy) setLastBlock(b *simplejson.Json, bn uint64) {
	_p.Lock()
	defer _p.Unlock()
	_p.blockNumber = bn
}

func (_p *Proxy) InitBlockSubscription() {
	log.Println("try connecting to:", _p.address)
	lastBlockNumber := 0
	go func() {
		for {
			func() {
				defer time.Sleep(BlockQueryInterval)
				b, err := _p.GetLatestBlock()
				if err != nil {
					return
				}
				bn, err := strconv.Atoi(b.Get("block").Get("lastCommit").Get("height").MustString())
				if err != nil {
					return
				}
				if bn > lastBlockNumber {
					_p.blockChan <- b
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
