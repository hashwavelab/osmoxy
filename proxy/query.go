package proxy

import (
	"context"
	"errors"
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

func (P *Proxy) connect() (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), GRPCConnectTimeOut)
	defer cancel()
	conn, err := grpc.DialContext(ctx, P.address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	return conn, err
}

func (P *Proxy) GetPools() (*gamm.QueryPoolsResponse, error) {
	conn, err := P.connect()
	if err != nil {
		return nil, errors.New("cannot connect to grpc")
	}
	defer conn.Close()
	c := gamm.NewQueryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), GRPCQueryTimeOut)
	defer cancel()
	return c.Pools(ctx, &gamm.QueryPoolsRequest{Pagination: &query.PageRequest{Limit: 10000}})
}

func (P *Proxy) GetBalances(accAddress string) (*bank.QueryAllBalancesResponse, error) {
	conn, err := P.connect()
	if err != nil {
		return nil, errors.New("cannot connect to grpc")
	}
	defer conn.Close()
	c := bank.NewQueryClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), GRPCQueryTimeOut)
	defer cancel()
	return c.AllBalances(ctx, &bank.QueryAllBalancesRequest{Address: accAddress})
}

func (P *Proxy) GetLatestBlock() (*base.GetLatestBlockResponse, error) {
	conn, err := P.connect()
	if err != nil {
		return nil, errors.New("cannot connect to grpc")
	}
	defer conn.Close()
	c := base.NewServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), GRPCQueryTimeOut)
	defer cancel()
	return c.GetLatestBlock(ctx, &base.GetLatestBlockRequest{})
}

func (P *Proxy) GetTx(hash string) (*tx.GetTxResponse, error) {
	conn, err := P.connect()
	if err != nil {
		return nil, errors.New("cannot connect to grpc")
	}
	defer conn.Close()
	c := tx.NewServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), GRPCQueryTimeOut)
	defer cancel()
	return c.GetTx(ctx, &tx.GetTxRequest{Hash: hash})
}
