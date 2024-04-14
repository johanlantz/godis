// Resp protocol models for requests and responses
package main

var Suffix = "\r\n"

type RespCommand string

const (
	RESP_GET RespCommand = "GET"
	RESP_SET RespCommand = "SET"
)

type RespRequest struct {
	command RespCommand
	args    []string
}

// Build a RespRequest struct. This relies on that priamry validation has
// been performed beforehand to ensure the command is implemented and has
// the minimum required arg count. The reason for not validating here
// is that:
//  1. This keeps the models responsabilities to only being a model
//  2. Even if the command is implemented and the arg count is correct,
//     the command can still be errounous and we will catch that when trying to execute.
func newRespRequest(cmd_arr []string) (*RespRequest, error) {
	return &RespRequest{command: RespCommand(cmd_arr[0]), args: cmd_arr[1:]}, nil
}

type ResponseDataType byte

const (
	RESP_SIMPLE_STRING ResponseDataType = '+'
	RESP_SIMPLE_ERROR  ResponseDataType = '-'
)

var RESP_ERR = "ERR"
var RESP_OK = "OK"

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
