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

type RespRequest struct {
	command RespCommand
	args    []string
}

var supportedCommands = []RespCommand{RESP_GET, RESP_SET}

// Perform basic validation and Build a RespRequest from an incoming command.
func newRespRequest(bytes []byte) (*RespRequest, error) {
	cmd := string(bytes)

	// 1. Must be a bulk string array starting with * and ending with \r\n.
	if len(cmd) < len(suffix) || cmd[0] != '*' || cmd[len(cmd)-len(suffix):] != suffix {
		return nil, errors.New("invalid command")
	}

	// 2. Generate our array of command segments from the bulk string array.
	bulk_array := strings.Fields(cmd[:len(cmd)-len(suffix)])
	cmd_arr := []string{}
	for _, element := range bulk_array {
		if element[0] != DT_BULK_STRINGS && element[0] != DT_ARRAYS {
			cmd_arr = append(cmd_arr, element)
		}
	}
	cmd_verb := RespCommand(cmd_arr[0])

	// 3. The command must be supported by our current implementation
	if !slices.Contains(supportedCommands, cmd_verb) {
		return nil, fmt.Errorf("unknown command, %s", cmd_verb)
	}

	// 4. Each command processor is responsible for validating the args later on.
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
