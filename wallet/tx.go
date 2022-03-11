package wallet

import (
	tx "github.com/cosmos/cosmos-sdk/api/cosmos/tx/v1beta1"
	"github.com/hashwavelab/osmoxy/pb"
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
