package resp

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/johanlantz/redis/storage"
	"github.com/johanlantz/redis/utils"
	"github.com/stretchr/testify/require"
)

var requestChannel = make(chan NetworkRequest)
var responseChannel = make(chan []byte)

func setup() {
	StartCommandProcessor(requestChannel, storage.NewSimpleStorage())
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestInvalidCommand(t *testing.T) {
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("SETI")}
	response := <-responseChannel
	require.Contains(t, string(response), RESP_ERR)
}

func TestGetWithoutKey(t *testing.T) {
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("GET")}
	response := <-responseChannel
	require.Contains(t, string(response), RESP_ERR)
}

func TestSetWithoutKey(t *testing.T) {
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("SET \r\n")}
	response := <-responseChannel
	require.Contains(t, string(response), RESP_ERR)
}

func TestGetWhenNoValueStored(t *testing.T) {
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("GET masterKey")}
	response := <-responseChannel
	require.Contains(t, response, byte(DT_NULLS))
}

func TestSetWithoutValue(t *testing.T) {
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("SET masterKey")}
	response := <-responseChannel
	require.Contains(t, string(response), RESP_ERR)
}

func TestSet(t *testing.T) {
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("SET masterKey myValue")}
	response := <-responseChannel
	require.Equal(t, "+OK\r\n", string(response))
}

func TestSetGetString(t *testing.T) {
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("SET masterKey myValue")}
	response := <-responseChannel
	require.Equal(t, "+OK\r\n", string(response))

	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("GET masterKey")}
	response = <-responseChannel
	require.Equal(t, "+myValue\r\n", string(response))
}

func TestInteger(t *testing.T) {
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("SET myIntCounter 5")}
	response := <-responseChannel
	require.Equal(t, "+OK\r\n", string(response))

	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("GET myIntCounter")}
	response = <-responseChannel
	require.Equal(t, ":5\r\n", string(response))
}

func TestSetFloat(t *testing.T) {
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("SET myFloatCounter 5.4")}
	response := <-responseChannel
	require.Equal(t, "+OK\r\n", string(response))

	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("GET myFloatCounter")}
	response = <-responseChannel
	require.Equal(t, ",5.4\r\n", string(response))
}

func TestBool(t *testing.T) {
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("SET myBool true")}
	response := <-responseChannel
	require.Equal(t, "+OK\r\n", string(response))

	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("GET myBool")}
	response = <-responseChannel
	require.Equal(t, "#t\r\n", string(response))

	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("SET myBool false")}
	response = <-responseChannel
	require.Equal(t, "+OK\r\n", string(response))

	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("GET myBool")}
	response = <-responseChannel
	require.Equal(t, "#f\r\n", string(response))
}

func TestIncrWithNilValue(t *testing.T) {
	var response []byte
	for i := 0; i < 15; i++ {
		requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("INCR myKey")}
		response = <-responseChannel
		require.Equal(t, "+OK\r\n", string(response))
	}
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("GET myKey")}
	response = <-responseChannel
	require.Equal(t, ":15\r\n", string(response))
}

func TestIncrWithStartValue(t *testing.T) {
	var response []byte
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("SET myKey 99")}
	response = <-responseChannel
	require.Equal(t, "+OK\r\n", string(response))

	for i := 0; i < 5; i++ {
		requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("INCR myKey")}
		response = <-responseChannel
		require.Equal(t, "+OK\r\n", string(response))
	}
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("GET myKey")}
	response = <-responseChannel
	require.Equal(t, ":104\r\n", string(response))
}

func TestIncrWithIncorrectValueType(t *testing.T) {
	var response []byte
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("SET myStringKey hello")}
	response = <-responseChannel
	require.Equal(t, "+OK\r\n", string(response))

	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("INCR myStringKey")}
	response = <-responseChannel
	require.Contains(t, string(response), "WRONGTYPE")
}

func TestConcurrency(t *testing.T) {
	var wg sync.WaitGroup

	var response []byte
	count := 100
	wg.Add(count)

	for i := 0; i < count; i++ {
		respCh := make(chan []byte)
		go func() {
			defer wg.Done()
			requestChannel <- NetworkRequest{ResponseChannel: respCh, Data: utils.MarshalToResp("INCR TestConcurrencyKey")}
			response = <-respCh
			require.Equal(t, "+OK\r\n", string(response))
		}()
	}
	wg.Wait()
	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("GET TestConcurrencyKey")}
	response = <-responseChannel
	require.Equal(t, ":100\r\n", string(response))
}

func TestDelete(t *testing.T) {
	var response []byte

	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("DEL missingKey")}
	response = <-responseChannel
	require.Contains(t, string(response), "0")

	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("SET myStringKey hello")}
	response = <-responseChannel
	require.Equal(t, "+OK\r\n", string(response))

	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("DEL myStringKey")}
	response = <-responseChannel
	require.Contains(t, string(response), "1")

	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp("DEL myStringKey")}
	response = <-responseChannel
	require.Contains(t, string(response), "0")
}

func TestMultiDelete(t *testing.T) {
	var response []byte
	count := 5
	var keyList []string
	for i := 0; i < count; i++ {
		requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp((fmt.Sprintf("SET myStringKey%d hello", i)))}
		response = <-responseChannel
		require.Equal(t, "+OK\r\n", string(response))
		keyList = append(keyList, fmt.Sprintf("myStringKey%d", i))
	}

	requestChannel <- NetworkRequest{ResponseChannel: responseChannel, Data: utils.MarshalToResp(fmt.Sprintf("DEL %s", strings.Join(keyList, " ")))}
	response = <-responseChannel
	require.Contains(t, string(response), fmt.Sprint(count))
}
