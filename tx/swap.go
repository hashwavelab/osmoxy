package tx

import (
	"errors"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/hashwavelab/osmoxy/pb"
	"github.com/hashwavelab/osmoxy/proxy"
	"github.com/hashwavelab/osmoxy/proxy/execproxy"
)

var (
	TransactionTimeout = 30 * time.Second
)

func SwapUsingOsmosisd(proxy *proxy.Proxy, params *pb.SwapParams) (string, error) {
	cmd := execproxy.NewOsmosisdCommand()
	if params.SwapExactTo {
		cmd = cmd.SwapExactAmountIn(params.ExactAmount+params.ExactDenom, params.ReqAmount)
	} else {
		cmd = cmd.SwapExactAmountOut(params.ExactAmount+params.ExactDenom, params.ReqAmount)
	}
	for _, route := range params.SwapRoutes {
		cmd = cmd.AddRoute(strconv.FormatInt(int64(route.PoolId), 10), route.Denom)
	}
	bytes, err := cmd.From(params.WalletAddress).OsmosisChainId().TestKeyringBackEnd().SkipConfirmation().Execute()
	if err != nil {
		log.Println("swap error:", err)
		return "", nil
	}
	str := strings.TrimSpace(string(bytes))
	tx := str[len(strings.TrimSpace(string(bytes)))-64:]
	return tx, nil
}

func WaitForTxResponse(proxy *proxy.Proxy, hash string) (*tx.GetTxResponse, error) {
	ch := make(chan interface{})
	proxy.SubscribeNewBlock(ch, 1)
	defer proxy.UnsubscribeNewBlock(ch)
	timec := time.After(TransactionTimeout)
	for {
		select {
		case <-ch:
			resp, err := proxy.GetTx(hash)
			if err != nil {
				continue
			} else {
				return resp, nil
			}
		case <-timec:
			log.Println("ExecuteRPDV2Trade ERROR -99 timeout")
			return nil, errors.New("wait for transaction to be mined timeout")
		}
	}
}

func GetSwapResultFromTxResponse(resp *tx.GetTxResponse) (*pb.SwapResult, error) {
	status := false
	if resp.TxResponse.Code == 0 {
		status = true
	}
	var gasFlow pb.Asset
	if len(resp.Tx.AuthInfo.Fee.Amount) == 0 {
		gasFlow = pb.Asset{
			Denom:  "uosmo",
			Amount: "0",
		}
	}
	gasFlow = pb.Asset{
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
						if _, ok := assetFlowMap[n]; !ok {
							assetFlowMap[n] = -a
						} else {
							assetFlowMap[n] -= a
						}
					case "tokens_out":
						a, n := separator(data.Value)
						if _, ok := assetFlowMap[n]; !ok {
							assetFlowMap[n] = a
						} else {
							assetFlowMap[n] += a
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
			// the second part of data is asset name
			name := data[i:]
			// the first part of data is amount
			amount, err := strconv.Atoi(data[0:i])
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
	return GetSwapResultFromTxResponse(resp)
}
