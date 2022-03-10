package main

import (
	"context"
	"log"
	"net"

	"github.com/hashwavelab/osmoxy/gamm"
	"github.com/hashwavelab/osmoxy/pb"
	"github.com/hashwavelab/osmoxy/proxy"
	"google.golang.org/grpc"
)

const (
	port = ":9094"
)

type server struct {
	pb.UnimplementedOsmoxyServer
	proxy *proxy.Proxy
	dex   *gamm.Dex
}

func (s *server) GetPoolsSnapshot(ctx context.Context, in *pb.EmptyRequest) (*pb.PoolsSnapshot, error) {
	snapshot := &pb.PoolsSnapshot{
		BlockNumber: s.proxy.GetLastBlockNumber(),
		Pools:       s.dex.ExportPoolsSnapshot(),
	}
	return snapshot, nil
}

func (s *server) SubscribePoolsUpdate(in *pb.EmptyRequest, stream pb.Osmoxy_SubscribePoolsUpdateServer) error {
	return s.dex.SubscribePoolsUpdate(stream)
}

// Legacy method for compitability with UniV2:
func (s *server) GetUniV2PairsSnapshot(ctx context.Context, in *pb.EmptyRequest) (*pb.UniV2PairsSnapshot, error) {
	snapshot := &pb.UniV2PairsSnapshot{
		BlockNumber: s.proxy.GetLastBlockNumber(),
		Pairs:       s.dex.ExportUniV2PairsSnapshot(),
	}
	return snapshot, nil
}

func (s *server) SubscribeUniV2PairsUpdate(in *pb.EmptyRequest, stream pb.Osmoxy_SubscribeUniV2PairsUpdateServer) error {
	return s.dex.SubscribeUniV2PairsUpdate(stream)
}

func InitGrpcServer(p *proxy.Proxy, d *gamm.Dex) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOsmoxyServer(s, &server{proxy: p, dex: d})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
