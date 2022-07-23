package proxy

import (
	"context"
	"time"

	base "github.com/cosmos/cosmos-sdk/api/cosmos/base/tendermint/v1beta1"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/cosmos/cosmos-sdk/types/tx"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	gamm "github.com/osmosis-labs/osmosis/v7/x/gamm/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	GRPCConnectTimeOut time.Duration = 3 * time.Second
	GRPCQueryTimeOut   time.Duration = 5 * time.Second
)

func (p *Proxy) connect() (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPCConnectTimeOut)
	defer cancel()
	conn, err := grpc.DialContext(ctx, p.address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	return conn, err
}

func (p *Proxy) GetPools() (*gamm.QueryPoolsResponse, error) {
	c := gamm.NewQueryClient(p.conn)
	ctx, cancel := context.WithTimeout(context.Background(), GRPCQueryTimeOut)
	defer cancel()
	return c.Pools(ctx, &gamm.QueryPoolsRequest{Pagination: &query.PageRequest{Limit: 10000}})
}

func (p *Proxy) GetBalances(accAddress string) (*bank.QueryAllBalancesResponse, error) {
	c := bank.NewQueryClient(p.conn)
	ctx, cancel := context.WithTimeout(context.Background(), GRPCQueryTimeOut)
	defer cancel()
	return c.AllBalances(ctx, &bank.QueryAllBalancesRequest{Address: accAddress})
}

func (p *Proxy) GetLatestBlock() (*base.GetLatestBlockResponse, error) {
	c := base.NewServiceClient(p.conn)
	ctx, cancel := context.WithTimeout(context.Background(), GRPCQueryTimeOut)
	defer cancel()
	return c.GetLatestBlock(ctx, &base.GetLatestBlockRequest{})
}

func (p *Proxy) GetTx(hash string) (*tx.GetTxResponse, error) {
	c := tx.NewServiceClient(p.conn)
	ctx, cancel := context.WithTimeout(context.Background(), GRPCQueryTimeOut)
	defer cancel()
	return c.GetTx(ctx, &tx.GetTxRequest{Hash: hash})
}
