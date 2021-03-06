package gamm

import (
	"log"
	"strconv"
	"time"

	"github.com/hashwavelab/osmoxy/pb"
)

func (d *Dex) ExportPoolsSnapshot() []*pb.Pool {
	r := make([]*pb.Pool, 0)
	d.pools.Range(func(k, v interface{}) bool {
		id, assets, feeN, feeD := v.(*Pool).export(true)
		r = append(r, &pb.Pool{
			Id:     id,
			Assets: assets,
			FeeN:   feeN,
			FeeD:   feeD,
		})
		return true
	})
	return r
}

func (d *Dex) SubscribePoolsUpdate(stream pb.Osmoxy_SubscribePoolsUpdateServer) error {
	ch := make(chan interface{})
	d.broadcaster.Register(ch, 100)
	defer d.broadcaster.Unregister(ch)
	for u := range ch {
		updates := u.([]*Pool)
		r := &pb.PoolsUpdate{
			FromBlockNumber: d.proxy.GetLastBlockNumber(),
			ToBlockNumber:   d.proxy.GetLastBlockNumber(),
			Timestamp:       uint64(time.Now().UnixMilli()),
			Updates:         make([]*pb.PoolUpdate, 0),
		}
		for _, u := range updates {
			id, assets, _, _ := u.export(false)
			r.Updates = append(r.Updates, &pb.PoolUpdate{
				Id:     id,
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

// Legacy method to be compatiable with UniV2 pairs:
func (d *Dex) ExportUniV2PairsSnapshot() []*pb.UniV2Pair {
	r := make([]*pb.UniV2Pair, 0)
	d.pools.Range(func(k, v interface{}) bool {
		pool := v.(*Pool)
		if !pool.isUniV2() {
			return true
		}
		id, assets, feeN, feeD := v.(*Pool).export(true)
		r = append(r, &pb.UniV2Pair{
			PairAddress: strconv.Itoa(int(id)),
			Token0:      assets[0].Denom,
			Token1:      assets[1].Denom,
			Reserve0:    assets[0].Amount,
			Reserve1:    assets[1].Amount,
			FeeN:        feeN,
			FeeD:        feeD,
		})
		return true
	})
	return r
}

func (d *Dex) SubscribeUniV2PairsUpdate(stream pb.Osmoxy_SubscribeUniV2PairsUpdateServer) error {
	ch := make(chan interface{})
	d.broadcaster.Register(ch, 100)
	defer log.Println("unsub")
	defer d.broadcaster.Unregister(ch)
	for u := range ch {
		updates := u.([]*Pool)
		r := &pb.UniV2PairsUpdate{
			FromBlockNumber: d.proxy.GetLastBlockNumber(),
			ToBlockNumber:   d.proxy.GetLastBlockNumber(),
			Timestamp:       uint64(time.Now().UnixMilli()),
			Univ2Updates:    make([]*pb.UniV2PairUpdate, 0),
		}
		for _, u := range updates {
			if !u.isUniV2() {
				continue
			}
			id, assets, _, _ := u.export(false)
			r.Univ2Updates = append(r.Univ2Updates, &pb.UniV2PairUpdate{
				PairAddress: strconv.Itoa(int(id)),
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
