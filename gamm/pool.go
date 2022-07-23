package gamm

import (
	"log"
	"sync"

	"github.com/hashwavelab/osmoxy/gamm/proto"
	"github.com/hashwavelab/osmoxy/pb"
)

type Pool struct {
	// fixed
	sync.RWMutex
	Id         uint64
	PoolAssets []*PoolAsset
	FeeN       uint64
	FeeD       uint64
}

func (p *Pool) update(pl *proto.Pool) *Pool {
	p.Lock()
	defer p.Unlock()
	updated := false
	ul := len(pl.PoolAssets)
	ol := len(p.PoolAssets)
	if ul != ol {
		log.Println("Pool update error - number of assets not match", p)
		return nil
	}
	for i, a := range pl.PoolAssets {
		ua := a.Token.Amount.String()
		if ua != p.PoolAssets[i].Amount {
			p.PoolAssets[i].Amount = ua
			updated = true
		}
	}
	if updated {
		return p
	} else {
		return nil
	}
}

func (p *Pool) export(includeWeight bool) (uint64, []*pb.PoolAsset, uint64, uint64) {
	pae := make([]*pb.PoolAsset, 0)
	p.RLock()
	defer p.RUnlock()
	for _, a := range p.PoolAssets {
		ae := &pb.PoolAsset{
			Denom:  a.Denom,
			Amount: a.Amount,
		}
		if includeWeight {
			ae.Weight = a.Weight
		}
		pae = append(pae, ae)
	}
	return p.Id, pae, p.FeeN, p.FeeD
}

// Return true if the pool is a UniV2 type pool with two assets and equal weights.
func (p *Pool) isUniV2() bool {
	p.RLock()
	defer p.RUnlock()
	if len(p.PoolAssets) != 2 {
		return false
	}
	return p.PoolAssets[0].Weight == p.PoolAssets[1].Weight
}

type PoolAsset struct {
	Denom  string
	Amount string
	Weight string
}
