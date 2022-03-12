package wallet

import (
	"log"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/hashwavelab/osmoxy/pb"
	"github.com/hashwavelab/osmoxy/proxy"
)

func getSwapResultFromTxResponse(resp *tx.GetTxResponse) (*pb.SwapResult, error) {
	status := false
	if resp.TxResponse.Code == 0 {
		status = true
	}
	var gasFlow pb.Asset = pb.Asset{
		Denom:  resp.Tx.AuthInfo.Fee.Amount[0].Denom,
		Amount: resp.Tx.AuthInfo.Fee.Amount[0].Amount.String(),
	}
	var assetFlowArr []*pb.Asset
	assetFlowMap := make(map[string]int)
	// suppose to loop once
	for _, log := range resp.TxResponse.Logs {
		// supopose to loop once as well
		for _, event := range log.Events {
			if event.Type == "token_swapped" {
				for _, data := range event.Attributes {
					switch data.Key {
					case "tokens_in":
						a, n := separator(data.Value)
						if _, ok := assetFlowMap[n]; ok {
							assetFlowMap[n] -= a
						} else {
							assetFlowMap[n] = -a
						}
					case "tokens_out":
						a, n := separator(data.Value)
						if _, ok := assetFlowMap[n]; ok { // idk why there is a warning here
							assetFlowMap[n] += a
						} else {
							assetFlowMap[n] = a
						}
					}
				}
			}
		}
	}
	for token := range assetFlowMap {
		// ignore the asset which flow amount is 0
		if assetFlowMap[token] != 0 {
			assetFlow := pb.Asset{
				Denom:  token,
				Amount: strconv.Itoa(assetFlowMap[token]),
			}
			assetFlowArr = append(assetFlowArr, &assetFlow)
		}
	}
	return &pb.SwapResult{
		Status:     status,
		Hash:       resp.TxResponse.TxHash,
		AssetFlows: assetFlowArr,
		GasFlow:    &gasFlow,
		GasUsed:    uint64(resp.TxResponse.GasUsed),
	}, nil
}

func separator(data string) (int, string) {
	for i := 0; i < len(data); i++ {
		if '0' <= data[i] && data[i] <= '9' {
			continue
		} else {
			name := data[i:]                       // the second part of data is asset name
			amount, err := strconv.Atoi(data[0:i]) // the first part of data is amount
			if err != nil {
				log.Println(err)
			}
			return amount, name
		}
	}
	return 0, ""
}

// For testing
func QuerySwapResultByHash(p *proxy.Proxy, hash string) (*pb.SwapResult, error) {
	resp, err := p.GetTx(hash)
	if err != nil {
		return nil, err
	}
	return getSwapResultFromTxResponse(resp)
}
