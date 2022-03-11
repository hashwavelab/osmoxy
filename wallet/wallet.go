package wallet

import (
	"log"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/hashwavelab/osmoxy/broadcast"
	"github.com/hashwavelab/osmoxy/proxy"
)

type Wallet struct {
	sync.RWMutex
	proxy       *proxy.Proxy
	address     string
	assets      []*asset
	broadcaster broadcast.Broadcaster
}

type asset struct {
	Denom  string
	Amount string
}

func NewWallet(proxy *proxy.Proxy, address string) *Wallet {
	w := &Wallet{
		proxy:       proxy,
		address:     address,
		assets:      make([]*asset, 0),
		broadcaster: broadcast.NewBroadcaster(),
	}
	return w
}

func (W *Wallet) InitQueryingBalancesAfterEveryNewBlock() {
	ch := make(chan interface{})
	W.proxy.SubscribeNewBlock(ch, 0)
	go func() {
		for range ch {
			start := time.Now()
			resp, err := W.proxy.GetBalances(W.address)
			if err != nil {
				log.Println("get balances err", err)
				continue
			}
			updatedBals := W.handleBalances(resp)
			W.broadcaster.Submit(updatedBals)
			log.Println("Balance updated", updatedBals, time.Since(start))
		}
	}()
}

func (W *Wallet) handleBalances(resp *types.QueryAllBalancesResponse) []*asset {
	assets := make([]*asset, 0)
	for _, b := range resp.Balances {
		assets = append(assets, &asset{
			Denom:  b.Denom,
			Amount: b.Amount.String(),
		})
	}
	W.Lock()
	defer W.Unlock()
	W.assets = assets
	return assets
}
