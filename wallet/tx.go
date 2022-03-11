package wallet

import (
	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/hashwavelab/osmoxy/pb"
	"github.com/hashwavelab/osmoxy/proxy"
)

func getSwapResultFromTxResponse(resp *tx.GetTxResponse) (*pb.SwapResult, error) {
	status := false
	if resp.TxResponse.Code == 0 {
		status = true
	}
	return &pb.SwapResult{
		Status: status,
	}, nil
}

// For testing
func QuerySwapResultByHash(p *proxy.Proxy, hash string) (*pb.SwapResult, error) {
	resp, err := p.GetTx(hash)
	if err != nil {
		return nil, err
	}
	return getSwapResultFromTxResponse(resp)
}
