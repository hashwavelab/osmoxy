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

func (_c *GRPCurlCommand) Execute() ([]byte, error) {
	args := append(_c.flags, []string{_c.address, _c.service}...)
	return exec.Command(_c.name, args...).Output()
}

func (_c *GRPCurlCommand) MaxTime(t float64) *GRPCurlCommand {
	_c.flags = append(_c.flags, []string{"-max-time", fmt.Sprintf("%f", t)}...)
	return _c
}

func (_c *GRPCurlCommand) Data(data string) *GRPCurlCommand {
	_c.flags = append(_c.flags, []string{"-d", data}...)
	return _c
}

func (_c *GRPCurlCommand) Address(a string) *GRPCurlCommand {
	_c.address = a
	return _c
}

func (_c *GRPCurlCommand) Service(a string) *GRPCurlCommand {
	_c.service = a
	return _c
}
