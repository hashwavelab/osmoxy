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
	Fee        string
}

func (P *Pool) update(p *proto.Pool) *Pool {
	P.Lock()
	defer P.Unlock()
	updated := false
	ul := len(p.PoolAssets)
	ol := len(P.PoolAssets)
	if ul != ol {
		log.Println("Pool update error - number of assets not match", P)
		return nil
	}
	for i, a := range p.PoolAssets {
		ua := a.Token.Amount.String()
		if ua != P.PoolAssets[i].Amount {
			P.PoolAssets[i].Amount = ua
			updated = true
		}
	}
	if updated {
		return P
	} else {
		return nil
	}
}

func (P *Pool) export(includeWeight bool) (uint64, []*pb.PoolAsset, string) {
	pae := make([]*pb.PoolAsset, 0)
	P.RLock()
	defer P.RUnlock()
	for _, a := range P.PoolAssets {
		ae := &pb.PoolAsset{
			Denom:  a.Denom,
			Amount: a.Amount,
		}
		if includeWeight {
			ae.Weight = a.Weight
		}
		pae = append(pae, ae)
	}
	return P.Id, pae, P.Fee
}

// Return true if the pool is a UniV2 type pool with two assets and equal weights.
func (P *Pool) isUniV2() bool {
	P.RLock()
	defer P.RUnlock()
	if len(P.PoolAssets) != 2 {
		return false
	}
	return P.PoolAssets[0].Weight == P.PoolAssets[1].Weight
}

type PoolAsset struct {
	Denom  string
	Amount string
	Weight string
}
