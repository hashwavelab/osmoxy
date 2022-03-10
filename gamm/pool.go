package gamm

import (
	"log"
	"sync"

	"github.com/bitly/go-simplejson"
	"github.com/hashwavelab/osmoxy/pb"
)

type Pool struct {
	// fixed
	sync.RWMutex
	PairIndex  string
	PoolAssets []*PoolAsset
	Fee        string
}

func (_p *Pool) update(j *simplejson.Json) *Pool {
	_p.Lock()
	defer _p.Unlock()
	updated := false
	ul := len(j.Get("poolAssets").MustArray())
	ol := len(_p.PoolAssets)
	if ul != ol {
		log.Println("Pool update error - number of assets not match", _p)
		return nil
	}
	for i := 0; i < ul; i++ {
		ua := j.Get("poolAssets").GetIndex(i).Get("token").Get("amount").MustString()
		if ua != _p.PoolAssets[i].Amount {
			_p.PoolAssets[i].Amount = ua
			updated = true
		}
	}
	if updated {
		return _p
	} else {
		return nil
	}
}

func (_p *Pool) export(includeWeight bool) (string, []*pb.PoolAsset, string) {
	pae := make([]*pb.PoolAsset, 0)
	_p.RLock()
	defer _p.RUnlock()
	for _, a := range _p.PoolAssets {
		ae := &pb.PoolAsset{
			Denom:  a.Denom,
			Amount: a.Amount,
		}
		if includeWeight {
			ae.Weight = a.Weight
		}
		pae = append(pae, ae)
	}
	return _p.PairIndex, pae, _p.Fee
}

// Return true if the pool is a UniV2 type pool with two assets and equal weights.
func (_p *Pool) isUniV2() bool {
	_p.RLock()
	defer _p.RUnlock()
	if len(_p.PoolAssets) != 2 {
		return false
	}
	return _p.PoolAssets[0].Weight == _p.PoolAssets[1].Weight
}

type PoolAsset struct {
	Denom  string
	Amount string
	Weight string
}
