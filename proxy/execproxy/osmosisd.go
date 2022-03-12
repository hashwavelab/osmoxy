package execproxy

import "os/exec"

type OsmosisdCommand struct {
	name  string
	cmds  []string
	flags []string
}

func NewOsmosisdCommand() *OsmosisdCommand {
	return &OsmosisdCommand{
		name:  "osmosisd",
		cmds:  []string{},
		flags: []string{},
	}
}

func (C *OsmosisdCommand) Execute() ([]byte, error) {
	args := append(C.cmds, C.flags...)
	return exec.Command(C.name, args...).Output()
}

func (C *OsmosisdCommand) SwapExactAmountIn(amountAndDenom, reqAmount string) *OsmosisdCommand {
	C.cmds = []string{"tx", "gamm", "swap-exact-amount-in", amountAndDenom, reqAmount}
	return C
}

func (C *OsmosisdCommand) SwapExactAmountOut(amountAndDenom, reqAmount string) *OsmosisdCommand {
	C.cmds = []string{"tx", "gamm", "swap-exact-amount-out", amountAndDenom, reqAmount}
	return C
}

func (C *OsmosisdCommand) AddRoute(poolId, denom string) *OsmosisdCommand {
	C.flags = append(C.flags, []string{"--swap-route-pool-ids", poolId}...)
	C.flags = append(C.flags, []string{"--swap-route-denoms", denom}...)
	return C
}

func (C *OsmosisdCommand) From(accAddress string) *OsmosisdCommand {
	C.flags = append(C.flags, []string{"--from", accAddress}...)
	return C
}

func (C *OsmosisdCommand) OsmosisChainId() *OsmosisdCommand {
	C.flags = append(C.flags, []string{"--chain-id", "osmosis-1"}...)
	return C
}

func (C *OsmosisdCommand) TestKeyringBackEnd() *OsmosisdCommand {
	C.flags = append(C.flags, []string{"--keyring-backend", "test"}...)
	return C
}

func (C *OsmosisdCommand) SkipConfirmation() *OsmosisdCommand {
	C.flags = append(C.flags, "--yes")
	return C
}
