// Command processor handling the business logic for RESP.
// Accepts and returns a byte slice on the processing channel.
// It is not meant to be instanciated, hence no struct here.
package resp

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/johanlantz/redis/storage"
)

// Open up for different kinds of storage in the future
type KVStorage interface {
	Get(key string) storage.Entry
	Set(key string, value storage.Entry)
	Delete(key string)
}

type RespFunc = func(request *RespRequest, kv KVStorage) (*RespResponse, error)

// Implementing new commands only requires adding an entry here.
var processors = map[RespCommand]RespFunc{
	RESP_GET:  process_get,
	RESP_SET:  process_set,
	RESP_INCR: process_incr,
	RESP_DEL:  process_del,
}

// Redis proccesses in a single thread. This "event loop" provides the
// same behaviour while offering concurrency for the incoming connections.
// It also means the storage does not have to worry about race conditions.
func StartCommandProcessor(requestChannel <-chan []byte, responseChannel chan<- []byte, storage KVStorage) {
	go func() {
		for request := range requestChannel {
			processCommand(request, responseChannel, storage)
		}
	}()
}

func processCommand(bytes []byte, responseChannel chan<- []byte, storage KVStorage) {
	request, err := newRespRequest(bytes, &processors)
	var response *RespResponse

	if err != nil {
		response = newRespResponse(DT_SIMPLE_ERROR, []string{RESP_ERR, err.Error()})
		responseChannel <- response.marshalToBytes()
		return
	}

	response, err = processors[request.command](request, storage)

	if err != nil {
		response = newRespResponse(DT_SIMPLE_ERROR, []string{RESP_ERR, err.Error()})
	}
	responseChannel <- response.marshalToBytes()
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

func process_incr(request *RespRequest, kv KVStorage) (*RespResponse, error) {
	if len(request.args) != 1 {
		return nil, errors.New("incr command requires only key argument")
	}
	key := request.args[0]
	entry := kv.Get(request.args[0])

	if entry.IsNull() {
		kv.Set(key, storage.Entry{DataType: DT_INTEGER, Value: []byte("1")})
	} else if entry.DataType != DT_INTEGER {
		return nil, errors.New("WRONGTYPE existing value for key is not an integer")
	} else {
		if stored, err := strconv.Atoi(string(entry.Value)); err == nil {
			kv.Set(key, storage.Entry{DataType: DT_INTEGER, Value: []byte(fmt.Sprint(stored + 1))})
		} else {
			return nil, errors.New("FATAL storage corrupt")
		}
	}
	return newRespResponse(DT_SIMPLE_STRING, []string{RESP_OK}), nil
}

func process_del(request *RespRequest, kv KVStorage) (*RespResponse, error) {
	if len(request.args) < 1 {
		return nil, errors.New("del command requires at least one key argument")
	}
	deleteCount := 0
	for _, arg := range request.args {
		entry := kv.Get(arg)
		if !entry.IsNull() {
			kv.Delete(arg)
			deleteCount++
		}
	}
	return newRespResponse(DT_INTEGER, []string{fmt.Sprint(deleteCount)}), nil
}
