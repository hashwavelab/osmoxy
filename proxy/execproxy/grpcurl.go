package execproxy

import (
	"fmt"
	"os/exec"
)

type GRPCurlCommand struct {
	name    string
	flags   []string
	address string
	service string
}

func NewPlaintextGrpcurlCommand() *GRPCurlCommand {
	return &GRPCurlCommand{
		name:  "grpcurl",
		flags: []string{"--plaintext"},
	}
}

func (C *GRPCurlCommand) Execute() ([]byte, error) {
	args := append(C.flags, []string{C.address, C.service}...)
	return exec.Command(C.name, args...).Output()
}

func (C *GRPCurlCommand) MaxTime(t float64) *GRPCurlCommand {
	C.flags = append(C.flags, []string{"-max-time", fmt.Sprintf("%f", t)}...)
	return C
}

func (C *GRPCurlCommand) Data(data string) *GRPCurlCommand {
	C.flags = append(C.flags, []string{"-d", data}...)
	return C
}

func (C *GRPCurlCommand) Address(a string) *GRPCurlCommand {
	C.address = a
	return C
}

func (C *GRPCurlCommand) Service(a string) *GRPCurlCommand {
	C.service = a
	return C
}
