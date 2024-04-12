// Resp protocol parsing and construction
package main

import (
	"errors"
	"fmt"
	"strings"
)

type ResponseDataType byte

type RespResponse struct {
	t    ResponseDataType
	args []string
}

var suffix = "\r\n"

type RespCommand string

const (
	RESP_GET RespCommand = "GET"
	RESP_SET RespCommand = "SET"
)

type RespRequest struct {
	command RespCommand
	args    []string
}

// List of supported commands and the minimum required number of arguments for them to be valid.
var supportedRespCommands = map[RespCommand]int{RESP_GET: 1, RESP_SET: 2}

// Perform minimal validation to ensure the RESP command is implemented and that the called has
// provided at least the minimum amount of required arguments.
func validateRespCommand(cmd_arr []string) error {
	cmd := RespCommand(cmd_arr[0])

	if supportedRespCommands[cmd] == 0 {
		return fmt.Errorf("unsupported command: %s", cmd)
	}
	arg_count := len(cmd_arr) - 1
	min_args_required := supportedRespCommands[cmd]

	if arg_count < min_args_required {
		return fmt.Errorf("insuficient number of arguments, expected %d, received %d", min_args_required, arg_count)
	}

	return nil
}

// Build a RespRequest struct from incoming byte array
func newRespRequest(bytes []byte) (*RespRequest, error) {
	cmd := string(bytes)

	if len(cmd) == 0 {
		return nil, errors.New("empty resp command is not allowed")
	} else if len(cmd) < len(suffix) {
		return nil, errors.New("malformed resp command")
	} else if cmd[len(cmd)-len(suffix):] != suffix {
		return nil, errors.New("resp command is missing terminating newline")
	}

	cmd_arr := strings.Split(cmd, " ")
	err := validateRespCommand(cmd_arr)

	if err != nil {
		return nil, err
	}

	return &RespRequest{command: RespCommand(cmd_arr[0]), args: cmd_arr[1:]}, nil
}
