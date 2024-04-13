// CommandProcessor interface implementation for RESP.
// Invoked by the network layer
package main

import (
	"errors"
	"fmt"
	"strings"
)

type RespCommandProcessor struct{}

func newRespCommandProcessor() *RespCommandProcessor {
	return &RespCommandProcessor{}
}

func (rcp RespCommandProcessor) processCommand(bytes []byte) []byte {
	cmd := string(bytes)
	cmd_arr := strings.Split(cmd, " ")
	err := validateRespCommand(cmd_arr)

	if err != nil {
		response := newRespResponse(RESP_SIMPLE_ERROR, []string{"ERR"}).marshalToBytes()
		return response
	}

	return nil
}

// List of supported commands and the minimum required number of arguments for them to be valid.
var supportedRespCommands = map[RespCommand]int{RESP_GET: 1, RESP_SET: 2}

// Perform minimal validation to ensure the RESP command is implemented and that the caller has
// provided at least the minimum amount of required arguments.
func validateRespCommand(cmd_arr []string) error {
	cmd := RespCommand(cmd_arr[0])

	if len(cmd) == 0 {
		return errors.New("empty resp command is not allowed")
	} else if len(cmd) < len(Suffix) {
		return errors.New("malformed resp command")
	} else if string(cmd[len(cmd)-len(Suffix):]) != Suffix {
		return errors.New("resp command is missing terminating newline")
	}
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
