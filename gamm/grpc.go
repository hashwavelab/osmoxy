package gamm

import (
	"log"
	"strconv"
	"time"

	"github.com/hashwavelab/osmoxy/pb"
)

func (_d *Dex) ExportPoolsSnapshot() []*pb.Pool {
	r := make([]*pb.Pool, 0)
	_d.pools.Range(func(k, v interface{}) bool {
		index, assets, fee := v.(*Pool).export(true)
		r = append(r, &pb.Pool{
			Index:  index,
			Assets: assets,
			Fee:    fee,
		})
		return true
	})
	return r
}

func (D *Dex) SubscribePoolsUpdate(stream pb.Osmoxy_SubscribePoolsUpdateServer) error {
	ch := make(chan interface{})
	D.broadcaster.Register(ch, 100)
	defer log.Println("unsub")
	defer D.broadcaster.Unregister(ch)
	for u := range ch {
		updates := u.([]*Pool)
		r := &pb.PoolsUpdate{
			FromBlockNumber: D.proxy.GetLastBlockNumber(),
			ToBlockNumber:   D.proxy.GetLastBlockNumber(),
			Timestamp:       uint64(time.Now().UnixMilli()),
			Updates:         make([]*pb.PoolUpdate, 0),
		}
		for _, u := range updates {
			index, assets, _ := u.export(false)
			r.Updates = append(r.Updates, &pb.PoolUpdate{
				Index:  index,
				Assets: assets,
			})
		}
		if err := stream.Send(r); err != nil {
			log.Println("stream closed", err)
			return err
		}
	}
	return nil
}

// Legacy methods compatiable with UniV2 pairs:
func (_d *Dex) ExportUniV2PairsSnapshot() []*pb.UniV2Pair {
	r := make([]*pb.UniV2Pair, 0)
	_d.pools.Range(func(k, v interface{}) bool {
		pool := v.(*Pool)
		if !pool.isUniV2() {
			return true
		}
		index, assets, fee := v.(*Pool).export(true)
		fIntOriginal, err := strconv.Atoi(fee)
		if err != nil {
			return true
		}
		fInt := uint32(fIntOriginal / 100000000000000)
		r = append(r, &pb.UniV2Pair{
			PairAddress: index,
			Token0:      assets[0].Denom,
			Token1:      assets[1].Denom,
			Reserve0:    assets[0].Amount,
			Reserve1:    assets[1].Amount,
			FeeRev:      10000 - fInt,
		})
		return true
	})
	return r
}

func (D *Dex) SubscribeUniV2PairsUpdate(stream pb.Osmoxy_SubscribeUniV2PairsUpdateServer) error {
	ch := make(chan interface{})
	D.broadcaster.Register(ch, 100)
	defer log.Println("unsub")
	defer D.broadcaster.Unregister(ch)
	for u := range ch {
		updates := u.([]*Pool)
		r := &pb.UniV2PairsUpdate{
			FromBlockNumber: D.proxy.GetLastBlockNumber(),
			ToBlockNumber:   D.proxy.GetLastBlockNumber(),
			Timestamp:       uint64(time.Now().UnixMilli()),
			Univ2Updates:    make([]*pb.UniV2PairUpdate, 0),
		}
		for _, u := range updates {
			if !u.isUniV2() {
				continue
			}
			index, assets, _ := u.export(false)
			r.Univ2Updates = append(r.Univ2Updates, &pb.UniV2PairUpdate{
				PairAddress: index,
				Reserve0:    assets[0].Amount,
				Reserve1:    assets[1].Amount,
			})
		}
		if err := stream.Send(r); err != nil {
			log.Println("stream closed", err)
			return err
		}
	}
	return nil
}
