package gamm

import (
	"log"
	"sync"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	codecTypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/hashwavelab/osmoxy/broadcast"
	"github.com/hashwavelab/osmoxy/gamm/proto"
	"github.com/hashwavelab/osmoxy/proxy"
	"github.com/osmosis-labs/osmosis/v7/x/gamm/types"
)

var (
	amino = codec.NewLegacyAmino()
	cdc   = codec.NewAminoCodec(amino)
)

type Dex struct {
	proxy       *proxy.Proxy
	pools       sync.Map //map[uint64]*Pool
	broadcaster broadcast.Broadcaster
}

func NewDex(proxy *proxy.Proxy) *Dex {
	dex := &Dex{
		proxy:       proxy,
		broadcaster: broadcast.NewBroadcaster(),
	}
	return dex
}

func (D *Dex) InitQueryingPoolsAfterEveryNewBlock() {
	ch := make(chan interface{})
	D.proxy.SubscribeNewBlock(ch, 0)
	go func() {
		for range ch {
			start := time.Now()
			resp, err := D.proxy.GetPools()
			if err != nil {
				log.Println("get pools err", err)
				continue
			}
			updatedPools := D.handlePools(resp)
			D.broadcaster.Submit(updatedPools)
			log.Println("Pools updated", len(updatedPools), time.Since(start))
		}
	}()
}

func (D *Dex) handlePools(resp *types.QueryPoolsResponse) []*Pool {
	updatedPools := make([]*Pool, 0)
	for _, pool := range resp.Pools {
		up := D.updatePool(pool)
		if up != nil {
			updatedPools = append(updatedPools, up)
		}
	}
	return updatedPools
}

func (D *Dex) updatePool(p *codecTypes.Any) *Pool {
	var pool proto.Pool
	err := cdc.Amino.UnmarshalBinaryBare(p.Value, &pool)
	if err != nil {
		return nil
	}
	//
	pLoaded, ok := D.pools.Load(pool.Id)
	if !ok {
		np := newPool(&pool)
		D.pools.Store(pool.Id, np)
		return np
	} else {
		return pLoaded.(*Pool).update(&pool)
	}
}

func newPool(p *proto.Pool) *Pool {
	pool := &Pool{
		Id:         p.Id,
		PoolAssets: make([]*PoolAsset, 0),
		Fee:        p.PoolParams.SwapFee.String(),
	}
	for _, a := range p.PoolAssets {
		pool.PoolAssets = append(pool.PoolAssets, &PoolAsset{
			Denom:  a.Token.Denom,
			Amount: a.Token.Amount.String(),
			Weight: a.Weight.String(),
		})
	}
	return pool
}
