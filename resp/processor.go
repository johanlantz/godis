// CommandProcessor interface implementation for RESP.
// It accepts a byte slice from the network layer and ensures
// that the model creation is correct before passing it on.
// Once processed, it provides the network
// layer with a generic byte slice response in return.
package resp

import (
	"errors"
	"strconv"

	"github.com/johanlantz/redis/storage"
)

// The processor handles the business logic for RESP signalling
// and requires a storage. The storage however can be implemented
// with different characterics if needed.
type KVStorage interface {
	Get(key string) storage.Entry
	Set(key string, value storage.Entry)
}

type RespCommandProcessor struct {
	storage KVStorage
}

func NewRespCommandProcessor(storage KVStorage) *RespCommandProcessor {
	return &RespCommandProcessor{storage: storage}
}

func (rcp *RespCommandProcessor) ProcessCommand(bytes []byte) []byte {
	request, err := newRespRequest(bytes)
	var response *RespResponse

	if err != nil {
		response = newRespResponse(DT_SIMPLE_ERROR, []string{RESP_ERR, err.Error()})
		return response.marshalToBytes()
	}

	switch request.command {
	case RESP_GET:
		response, err = rcp.process_get(request)
	case RESP_SET:
		response, err = rcp.process_set(request)
	default:
		err = errors.New("unknown command")
	}

	if err != nil {
		response = newRespResponse(DT_SIMPLE_ERROR, []string{RESP_ERR, err.Error()})
	}
	return response.marshalToBytes()
}

func (rcp *RespCommandProcessor) process_get(request *RespRequest) (*RespResponse, error) {
	if len(request.args) < 1 {
		return nil, errors.New("get command requires at least the key parameter")
	}
	entry := rcp.storage.Get(request.args[0])
	if entry.IsNull() {
		return newRespResponse(DT_NULLS, []string{}), nil
	}
	return newRespResponse(ResponseDataType(entry.DataType), []string{string(entry.Value)}), nil
}

func (rcp *RespCommandProcessor) process_set(request *RespRequest) (*RespResponse, error) {
	if len(request.args) < 2 {
		return nil, errors.New("set command requires key and value parameters")
	}
	key := request.args[0]
	value := request.args[1]
	if _, err := strconv.Atoi(value); err == nil {
		rcp.storage.Set(key, storage.Entry{DataType: DT_INTEGER, Value: []byte(value)})
	} else if _, err := strconv.ParseFloat(value, 64); err == nil {
		rcp.storage.Set(key, storage.Entry{DataType: DT_DOUBLES, Value: []byte(value)})
	} else if _, err := strconv.ParseBool(value); err == nil {
		rcp.storage.Set(key, storage.Entry{DataType: DT_BOOLEANS, Value: []byte{value[0]}})
	} else {
		rcp.storage.Set(key, storage.Entry{DataType: DT_SIMPLE_STRING, Value: []byte(value)})
	}

	return newRespResponse(DT_SIMPLE_STRING, []string{RESP_OK}), nil
}
