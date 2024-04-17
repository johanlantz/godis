// Command processor handling the business logic for RESP.
// Accepts and returns a byte slice on the processing channel.
// It is not meant to be instanciated, hence no struct here.
package resp

import (
	"errors"
	"strconv"

	"github.com/johanlantz/redis/storage"
)

// Open up for different kinds of storage in the future
type KVStorage interface {
	Get(key string) storage.Entry
	Set(key string, value storage.Entry)
}

type RespFunc = func(request *RespRequest, kv KVStorage) (*RespResponse, error)

// These are all our implemented commands. Implementing new ones only requires
// adding an entry here with the corresponding processor function.
var processors = map[RespCommand]RespFunc{
	RESP_GET: process_get,
	RESP_SET: process_set,
}

// Redis proccesses in a single thread. This "event loop" provides the
// same behaviour while offering concurrency for the incoming connections.
// It also means the storage does not have to worry about race conditions.
func StartCommandProcessor(processingChannel chan []byte, storage KVStorage) {
	go func() {
		for request := range processingChannel {
			processCommand(request, processingChannel, storage)
		}
	}()
}

func processCommand(bytes []byte, processingChannel chan []byte, storage KVStorage) {
	request, err := newRespRequest(bytes)
	var response *RespResponse

	if err != nil {
		response = newRespResponse(DT_SIMPLE_ERROR, []string{RESP_ERR, err.Error()})
		processingChannel <- response.marshalToBytes()
		return
	}

	response, err = processors[request.command](request, storage)

	if err != nil {
		response = newRespResponse(DT_SIMPLE_ERROR, []string{RESP_ERR, err.Error()})
	}
	processingChannel <- response.marshalToBytes()
}

func process_get(request *RespRequest, kv KVStorage) (*RespResponse, error) {
	if len(request.args) < 1 {
		return nil, errors.New("get command requires at least the key parameter")
	}
	entry := kv.Get(request.args[0])
	if entry.IsNull() {
		return newRespResponse(DT_NULLS, []string{}), nil
	}
	return newRespResponse(ResponseDataType(entry.DataType), []string{string(entry.Value)}), nil
}

func process_set(request *RespRequest, kv KVStorage) (*RespResponse, error) {
	if len(request.args) < 2 {
		return nil, errors.New("set command requires key and value parameters")
	}
	key := request.args[0]
	value := request.args[1]
	if _, err := strconv.Atoi(value); err == nil {
		kv.Set(key, storage.Entry{DataType: DT_INTEGER, Value: []byte(value)})
	} else if _, err := strconv.ParseFloat(value, 64); err == nil {
		kv.Set(key, storage.Entry{DataType: DT_DOUBLES, Value: []byte(value)})
	} else if _, err := strconv.ParseBool(value); err == nil {
		kv.Set(key, storage.Entry{DataType: DT_BOOLEANS, Value: []byte{value[0]}})
	} else {
		kv.Set(key, storage.Entry{DataType: DT_SIMPLE_STRING, Value: []byte(value)})
	}
	return newRespResponse(DT_SIMPLE_STRING, []string{RESP_OK}), nil
}
