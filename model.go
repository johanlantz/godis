// Resp protocol models for requests and responses
package main

import (
	"errors"
	"fmt"
	"strings"
)

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

func isCommandSupported(value RespCommand) bool {
	supportedRespCommands := []RespCommand{RESP_GET, RESP_SET}
	for _, v := range supportedRespCommands {
		if v == value {
			return true
		}
	}
	return false
}

// Build a RespRequest struct from and incoming command.
func newRespRequest(bytes []byte) (*RespRequest, error) {
	cmd := string(bytes)

	// Perform only generic validations since each RESP command has its own
	// requirements for the args. This keeps the responsability of the model
	// simple while still allowing us to catch generic errors early.

	// 1. \r\n is always required
	if len(cmd) < len(suffix) || cmd[len(cmd)-len(suffix):] != suffix {
		return nil, errors.New("invalid command, missing terminating newline")
	}

	cmd_arr := strings.Split(cmd[:len(cmd)-len(suffix)], " ")
	cmd_verb := cmd_arr[0]

	// 2. The command must be supported by our current implementation
	if !isCommandSupported(RespCommand(cmd_verb)) {
		return nil, fmt.Errorf("unknown command, %s", cmd_verb)
	}

	// 3. Support args list but do not prevent empty args list either
	// Each command handle is responsible for the args later on.
	cmd_args := []string{}
	if len(cmd_arr) > 1 {
		cmd_args = cmd_arr[1:]
	}

	return &RespRequest{command: RespCommand(cmd_verb), args: cmd_args}, nil
}

type ResponseDataType byte

const (
	RESP_SIMPLE_STRING ResponseDataType = '+'
	RESP_SIMPLE_ERROR  ResponseDataType = '-'

	RESP_OK  = "OK"
	RESP_ERR = "ERR"
)

type RespResponse struct {
	t    ResponseDataType
	args []string
}

func newRespResponse(responseType ResponseDataType, args []string) *RespResponse {
	return &RespResponse{t: responseType, args: args}
}

func (rr RespResponse) marshalToBytes() []byte {
	return append([]byte{byte(rr.t)}, []byte(rr.args[0])...)
}
