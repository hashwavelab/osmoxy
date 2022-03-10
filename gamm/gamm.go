package gamm

import (
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/hashwavelab/osmoxy/broadcast"
	"github.com/hashwavelab/osmoxy/proxy"
)

type Dex struct {
	proxy       *proxy.Proxy
	pools       sync.Map //map[uint32]*Pool
	broadcaster broadcast.Broadcaster
}

func NewDex(proxy *proxy.Proxy) *Dex {
	dex := &Dex{
		proxy:       proxy,
		broadcaster: broadcast.NewBroadcaster(),
	}
	return dex
}

func (_d *Dex) InitQueryingPoolsAfterEveryNewBlock(ch chan interface{}) {
	go func() {
		for range ch {
			start := time.Now()
			pools, err := _d.proxy.GetPools()
			if err != nil {
				log.Println("get pools err", err)
				continue
			}
			updatedPools := _d.handlePools(pools)
			_d.broadcaster.Submit(updatedPools)
			log.Println("Pools updated", len(updatedPools), time.Since(start))
		}
	}()
}

func (_d *Dex) handlePools(j *simplejson.Json) []*Pool {
	updatedPools := make([]*Pool, 0)
	pools := j.Get("pools").MustArray()
	for i := range pools {
		up := _d.updatePool(j.Get("pools").GetIndex(i))
		if up != nil {
			updatedPools = append(updatedPools, up)
		}
	}
	return updatedPools
}

func (_d *Dex) updatePool(j *simplejson.Json) *Pool {
	i := j.Get("id").MustString()
	index, err := strconv.Atoi(i)
	if err != nil {
		return nil
	}
	p, ok := _d.pools.Load(index)
	if !ok {
		np := newPool(j)
		_d.pools.Store(index, np)
		return np
	} else {
		return p.(*Pool).update(j)
	}
}

func newPool(j *simplejson.Json) *Pool {
	l := len(j.Get("poolAssets").MustArray())
	fee := j.Get("poolParams").Get("swapFee").MustString()
	pair := &Pool{
		PairIndex:  j.Get("id").MustString(),
		PoolAssets: make([]*PoolAsset, 0),
		Fee:        fee,
	}
	for i := 0; i < l; i++ {
		pair.PoolAssets = append(pair.PoolAssets, &PoolAsset{
			Denom:  j.Get("poolAssets").GetIndex(i).Get("token").Get("denom").MustString(),
			Amount: j.Get("poolAssets").GetIndex(i).Get("token").Get("amount").MustString(),
			Weight: j.Get("poolAssets").GetIndex(i).Get("weight").MustString(),
		})
	}
	return pair
}
