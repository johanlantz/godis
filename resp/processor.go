// CommandProcessor interface implementation for RESP.
// It accepts a byte slice from the network layer, ensures
// that the model creation is correct before passing it on
// to the next layer. Once processed, it provdes the network
// layer with a generic byte slice response in return.
package resp

import (
	"errors"
	"strconv"

	"github.com/johanlantz/redis/storage"
)

type RespCommandProcessor struct{}

func NewRespCommandProcessor() *RespCommandProcessor {
	return &RespCommandProcessor{}
}

func (rcp RespCommandProcessor) ProcessCommand(bytes []byte) []byte {
	request, err := newRespRequest(bytes)
	var response *RespResponse

	if err != nil {
		response = newRespResponse(DT_SIMPLE_ERROR, []string{RESP_ERR, err.Error()})
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
		response = newRespResponse(DT_SIMPLE_ERROR, []string{RESP_ERR, err.Error()})
	}
	return response.marshalToBytes()
}

func process_get(request *RespRequest) (*RespResponse, error) {
	if len(request.args) < 1 {
		return nil, errors.New("get command requires at least the key parameter")
	}
	entry := storage.Get(request.args[0])
	return newRespResponse(ResponseDataType(entry.DataType), []string{string(entry.Value)}), nil
}

func process_set(request *RespRequest) (*RespResponse, error) {
	if len(request.args) < 2 {
		return nil, errors.New("set command requires key and value parameters")
	}

	key := request.args[0]
	value := request.args[1]
	if _, err := strconv.Atoi(value); err == nil {
		storage.Set(key, storage.StorageEntry{DataType: DT_INTEGER, Value: []byte(value)})
	} else {
		storage.Set(key, storage.StorageEntry{DataType: DT_SIMPLE_STRING, Value: []byte(value)})
	}

	return newRespResponse(DT_SIMPLE_STRING, []string{RESP_OK}), nil
}
