package gamm

import (
	"log"
	"math/big"
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

func (d *Dex) InitQueryingPoolsAfterEveryNewBlock() {
	ch := make(chan interface{})
	d.proxy.SubscribeNewBlock(ch, 0)
	go func() {
		for range ch {
			start := time.Now()
			resp, err := d.proxy.GetPools()
			if err != nil {
				log.Println("get pools err", err)
				continue
			}
			updatedPools := d.handlePools(resp)
			d.broadcaster.Submit(updatedPools)
			log.Println("Pools updated", len(updatedPools), time.Since(start))
		}
	}()
}

func (d *Dex) handlePools(resp *types.QueryPoolsResponse) []*Pool {
	updatedPools := make([]*Pool, 0)
	for _, pool := range resp.Pools {
		up := d.updatePool(pool)
		if up != nil {
			updatedPools = append(updatedPools, up)
		}
	}
	return updatedPools
}

func (d *Dex) updatePool(p *codecTypes.Any) *Pool {
	var pool proto.Pool
	err := cdc.Amino.UnmarshalBinaryBare(p.Value, &pool)
	if err != nil {
		return nil
	}
	//
	pLoaded, ok := d.pools.Load(pool.Id)
	if !ok {
		np := newPool(&pool)
		d.pools.Store(pool.Id, np)
		return np
	} else {
		return pLoaded.(*Pool).update(&pool)
	}
}

func newPool(p *proto.Pool) *Pool {
	feeRat, ok := new(big.Rat).SetString(p.PoolParams.SwapFee.String())
	if !ok {
		log.Println("New pool fee converts to rat failed", p)
		return nil
	}
	pool := &Pool{
		Id:         p.Id,
		PoolAssets: make([]*PoolAsset, 0),
		FeeN:       feeRat.Num().Uint64(),
		FeeD:       feeRat.Denom().Uint64(),
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
