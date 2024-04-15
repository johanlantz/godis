// Resp protocol models for requests and responses
package resp

import (
	"errors"
	"fmt"
	"slices"
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

var supportedRespCommands = []RespCommand{RESP_GET, RESP_SET}

// Build a generic RespRequest struct from an incoming command.
func newRespRequest(bytes []byte) (*RespRequest, error) {
	cmd := string(bytes)

	// 1. \r\n is always required
	if len(cmd) < len(suffix) || cmd[len(cmd)-len(suffix):] != suffix {
		return nil, errors.New("invalid command, missing terminating newline")
	}

	// 2. Generate our array of command segments. Unlike Split, Fields removes
	// multiple whitespaces automatically.
	cmd_arr := strings.Fields(cmd[:len(cmd)-len(suffix)])
	cmd_verb := RespCommand(cmd_arr[0])

	// 3. The command must be supported by our current implementation
	if !slices.Contains(supportedRespCommands, cmd_verb) {
		return nil, fmt.Errorf("unknown command, %s", cmd_verb)
	}

	// 4. Each command handler is responsible for its args validation later on.
	cmd_args := []string{}
	if len(cmd_arr) > 1 {
		cmd_args = cmd_arr[1:]
	}

	return &RespRequest{command: cmd_verb, args: cmd_args}, nil
}

type ResponseDataType byte

type RespResponse struct {
	t    ResponseDataType
	args []string
}

func newRespResponse(responseType ResponseDataType, args []string) *RespResponse {
	return &RespResponse{t: responseType, args: args}
}

func (rr RespResponse) marshalToBytes() []byte {
	bytes := []byte{byte(rr.t)}
	bytes = fmt.Append(bytes, strings.Join(rr.args, " "))
	bytes = fmt.Appendf(bytes, suffix)
	return bytes
}
