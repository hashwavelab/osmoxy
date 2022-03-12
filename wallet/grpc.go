package wallet

import (
	"log"

	"github.com/hashwavelab/osmoxy/pb"
	"github.com/hashwavelab/osmoxy/tx"
)

func (W *Wallet) ExportBalances() []*pb.Asset {
	assets := make([]*pb.Asset, 0)
	W.RLock()
	defer W.RUnlock()
	for _, a := range W.assets {
		assets = append(assets, &pb.Asset{
			Denom:  a.Denom,
			Amount: a.Amount,
		})
	}
	log.Println("exporting", assets)
	return assets
}

func (W *Wallet) SubscribeBalances(stream pb.Osmoxy_SubscribeBalancesServer) error {
	ch := make(chan interface{})
	W.broadcaster.Register(ch, 100)
	defer W.broadcaster.Unregister(ch)
	for u := range ch {
		assets := u.([]*asset)
		r := &pb.Balances{
			Assets: make([]*pb.Asset, 0),
		}
		for _, a := range assets {
			r.Assets = append(r.Assets, &pb.Asset{
				Denom:  a.Denom,
				Amount: a.Amount,
			})
		}
		if err := stream.Send(r); err != nil {
			log.Println("stream closed", err)
			return err
		}
	}
	return nil
}

func (W *Wallet) Swap(p *pb.SwapParams) (*pb.SwapResult, error) {
	hash, err := tx.SwapUsingOsmosisd(W.proxy, p)
	if err != nil {
		return nil, err
	}
	resp, err := tx.WaitForTxResponse(W.proxy, hash)
	if err != nil {
		return nil, err
	}
	return tx.GetSwapResultFromTxResponse(resp)
}
