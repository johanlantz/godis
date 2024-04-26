// RESP protocol models for requests and responses
package resp

import (
	"errors"
	"fmt"
	"strings"
)

var suffix = "\r\n"

type RespCommand string

type RespRequest struct {
	command RespCommand
	args    []string
}

// Perform basic validation and build a RespRequest from an incoming command.
func newRespRequest(bytes []byte, processors *map[RespCommand]RespFunc) (*RespRequest, error) {
	cmd := string(bytes)

	// 1. Must be a bulk string array starting with * and ending with \r\n.
	if len(cmd) < len(suffix) || cmd[0] != '*' || !strings.HasSuffix(cmd, suffix) {
		return nil, errors.New("invalid command")
	}

	// 2. Generate our array of command segments from the bulk string array.
	bulkArray := strings.Split(cmd[:len(cmd)-len(suffix)], suffix)
	cmdArray := []string{}
	for _, element := range bulkArray {
		if element[0] != DT_BULK_STRINGS && element[0] != DT_ARRAYS {
			cmdArray = append(cmdArray, element)
		}
	}
	cmdVerb := RespCommand(cmdArray[0])

	// 3. The command must be supported by our current implementation
	cmdSupported := false
	for key := range *processors {
		if cmdVerb == key {
			cmdSupported = true
		}
	}
	if !cmdSupported {
		return nil, fmt.Errorf("unknown command, %s", cmdVerb)
	}

	// 4. Each command processor is responsible for validating the args later on.
	cmd_args := []string{}
	if len(cmdArray) > 1 {
		cmd_args = cmdArray[1:]
	}

	return &RespRequest{command: cmdVerb, args: cmd_args}, nil
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
