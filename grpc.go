package main

import (
	"context"
	"errors"
	"log"
	"net"

	"github.com/hashwavelab/osmoxy/gamm"
	"github.com/hashwavelab/osmoxy/pb"
	"github.com/hashwavelab/osmoxy/proxy"
	"github.com/hashwavelab/osmoxy/wallet"
	"google.golang.org/grpc"
)

const (
	port = ":9094"
)

type server struct {
	pb.UnimplementedOsmoxyServer
	proxy   *proxy.Proxy
	dex     *gamm.Dex
	wallets map[string]*wallet.Wallet
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

func (s *server) GetBalances(ctx context.Context, in *pb.AddressRequest) (*pb.Balances, error) {
	wallet, ok := s.wallets[in.Address]
	if !ok {
		return nil, errors.New("wallet of given address is not found")
	}
	bals := &pb.Balances{
		Assets: wallet.ExportBalances(),
	}
	return bals, nil
}

func (s *server) SubscribeBalances(in *pb.AddressRequest, stream pb.Osmoxy_SubscribeBalancesServer) error {
	wallet, ok := s.wallets[in.Address]
	if !ok {
		return errors.New("wallet of given address is not found")
	}
	return wallet.SubscribeBalances(stream)
}

func (s *server) SubmitSwap(ctx context.Context, in *pb.SwapParams) (*pb.SwapResult, error) {
	log.Println("Received swap request:", in)
	wallet, ok := s.wallets[in.WalletAddress]
	if !ok {
		return nil, errors.New("wallet of given address is not found")
	}
	return wallet.Swap(in)
}

func InitGrpcServer(p *proxy.Proxy, d *gamm.Dex, wallets map[string]*wallet.Wallet) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterOsmoxyServer(s, &server{proxy: p, dex: d, wallets: wallets})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
