package pool

import (
	"time"

	cosmostypes "github.com/cosmos/cosmos-sdk/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v7/x/gamm/types"
)

type Pool struct {
	Address            string                                 `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty" yaml:"address"`
	Id                 uint64                                 `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	PoolParams         PoolParams                             `protobuf:"bytes,3,opt,name=poolParams,proto3" json:"poolParams" yaml:"pool_params"`
	FuturePoolGovernor string                                 `protobuf:"bytes,4,opt,name=future_pool_governor,json=futurePoolGovernor,proto3" json:"future_pool_governor,omitempty" yaml:"future_pool_governor"`
	TotalShares        cosmostypes.Coin                       `protobuf:"bytes,5,opt,name=totalShares,proto3" json:"totalShares" yaml:"total_shares"`
	PoolAssets         []types.PoolAsset                      `protobuf:"bytes,6,rep,name=poolAssets,proto3" json:"poolAssets" yaml:"pool_assets"`
	TotalWeight        github_com_cosmos_cosmos_sdk_types.Int `protobuf:"bytes,7,opt,name=totalWeight,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Int" json:"totalWeight" yaml:"total_weight"`
}

type PoolParams struct {
	SwapFee                  github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,1,opt,name=swapFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"swapFee" yaml:"swap_fee"`
	ExitFee                  github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=exitFee,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"exitFee" yaml:"exit_fee"`
	SmoothWeightChangeParams *SmoothWeightChangeParams              `protobuf:"bytes,3,opt,name=smoothWeightChangeParams,proto3" json:"smoothWeightChangeParams,omitempty" yaml:"smooth_weight_change_params"`
}

type SmoothWeightChangeParams struct {
	StartTime          time.Time         `protobuf:"bytes,1,opt,name=start_time,json=startTime,proto3,stdtime" json:"start_time" yaml:"start_time"`
	Duration           time.Duration     `protobuf:"bytes,2,opt,name=duration,proto3,stdduration" json:"duration,omitempty" yaml:"duration"`
	InitialPoolWeights []types.PoolAsset `protobuf:"bytes,3,rep,name=initialPoolWeights,proto3" json:"initialPoolWeights" yaml:"initial_pool_weights"`
	TargetPoolWeights  []types.PoolAsset `protobuf:"bytes,4,rep,name=targetPoolWeights,proto3" json:"targetPoolWeights" yaml:"target_pool_weights"`
}
