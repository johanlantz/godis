// CommandProcessor interface implementation for RESP.
// It accepts a byte slice from the network layer, ensures
// that the model creation is correct before passing it on
// to the next layer. Once processed, it provdes the network
// layer with a generic byte slice response in return.
package main

import "errors"

type RespCommandProcessor struct{}

func newRespCommandProcessor() *RespCommandProcessor {
	return &RespCommandProcessor{}
}

func (rcp RespCommandProcessor) processCommand(bytes []byte) []byte {
	request, err := newRespRequest(bytes)
	var response *RespResponse

	if err != nil {
		response = newRespResponse(RESP_SIMPLE_ERROR, []string{RESP_ERR, err.Error()})
		return response.marshalToBytes()
	}

	switch request.command {
	case RESP_GET:
		response, err = process_get(request)
	case RESP_SET:
		response, err = process_set(request)
	default:
		err = errors.New("unknown command")
	}

	if err != nil {
		response = newRespResponse(RESP_SIMPLE_ERROR, []string{RESP_ERR, err.Error()})
	}
	return response.marshalToBytes()
}

func process_get(request *RespRequest) (*RespResponse, error) {
	if len(request.args) < 1 {
		return nil, errors.New("get command requires at least the key parameter")
	}
	// TODO, add storage
	return newRespResponse(RESP_SIMPLE_STRING, []string{RESP_OK}), nil
}

func process_set(request *RespRequest) (*RespResponse, error) {
	if len(request.args) < 2 {
		return nil, errors.New("set command requires key and value parameters")
	}
	// TODO, add storage
	return newRespResponse(RESP_SIMPLE_STRING, []string{RESP_OK}), nil
}
