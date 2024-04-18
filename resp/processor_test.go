package resp

import (
	"os"
	"testing"

	"github.com/johanlantz/redis/storage"
	"github.com/johanlantz/redis/utils"
	"github.com/stretchr/testify/require"
)

var processingChannel = make(chan []byte)

func setup() {
	StartCommandProcessor(processingChannel, storage.NewSimpleStorage())
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestInvalidCommand(t *testing.T) {
	processingChannel <- utils.MarshalToResp("SETI")
	response := <-processingChannel
	require.Contains(t, string(response), RESP_ERR)
}

func TestGetWithoutKey(t *testing.T) {
	processingChannel <- utils.MarshalToResp("GET")
	response := <-processingChannel
	require.Contains(t, string(response), RESP_ERR)
}

func TestSetWithoutKey(t *testing.T) {
	processingChannel <- utils.MarshalToResp("SET \r\n")
	response := <-processingChannel
	require.Contains(t, string(response), RESP_ERR)
}

func TestGetWhenNoValueStored(t *testing.T) {
	processingChannel <- utils.MarshalToResp("GET masterKey")
	response := <-processingChannel
	require.Contains(t, response, byte(DT_NULLS))
}

func TestSetWithoutValue(t *testing.T) {
	processingChannel <- utils.MarshalToResp("SET masterKey")
	response := <-processingChannel
	require.Contains(t, string(response), RESP_ERR)
}

func TestSet(t *testing.T) {
	processingChannel <- utils.MarshalToResp("SET masterKey myValue")
	response := <-processingChannel
	require.Equal(t, "+OK\r\n", string(response))
}

func TestSetGetString(t *testing.T) {
	processingChannel <- utils.MarshalToResp("SET masterKey myValue")
	response := <-processingChannel
	require.Equal(t, "+OK\r\n", string(response))

	processingChannel <- utils.MarshalToResp("GET masterKey")
	response = <-processingChannel
	require.Equal(t, "+myValue\r\n", string(response))
}

func TestInteger(t *testing.T) {
	processingChannel <- utils.MarshalToResp("SET myIntCounter 5")
	response := <-processingChannel
	require.Equal(t, "+OK\r\n", string(response))

	processingChannel <- utils.MarshalToResp("GET myIntCounter")
	response = <-processingChannel
	require.Equal(t, ":5\r\n", string(response))
}

func TestSetFloat(t *testing.T) {
	processingChannel <- utils.MarshalToResp("SET myFloatCounter 5.4")
	response := <-processingChannel
	require.Equal(t, "+OK\r\n", string(response))

	processingChannel <- utils.MarshalToResp("GET myFloatCounter")
	response = <-processingChannel
	require.Equal(t, ",5.4\r\n", string(response))
}

func TestBool(t *testing.T) {
	processingChannel <- utils.MarshalToResp("SET myBool true")
	response := <-processingChannel
	require.Equal(t, "+OK\r\n", string(response))

	processingChannel <- utils.MarshalToResp("GET myBool")
	response = <-processingChannel
	require.Equal(t, "#t\r\n", string(response))

	processingChannel <- utils.MarshalToResp("SET myBool false")
	response = <-processingChannel
	require.Equal(t, "+OK\r\n", string(response))

	processingChannel <- utils.MarshalToResp("GET myBool")
	response = <-processingChannel
	require.Equal(t, "#f\r\n", string(response))
}

func TestIncrWithNilValue(t *testing.T) {
	var response []byte
	for i := 0; i < 15; i++ {
		processingChannel <- utils.MarshalToResp("INCR myKey")
		response = <-processingChannel
		require.Equal(t, "+OK\r\n", string(response))
	}
	processingChannel <- utils.MarshalToResp("GET myKey")
	response = <-processingChannel
	require.Equal(t, ":15\r\n", string(response))
}

func TestIncrWithStartValue(t *testing.T) {
	var response []byte
	processingChannel <- utils.MarshalToResp("SET myKey 99")
	response = <-processingChannel
	require.Equal(t, "+OK\r\n", string(response))

	for i := 0; i < 5; i++ {
		processingChannel <- utils.MarshalToResp("INCR myKey")
		response = <-processingChannel
		require.Equal(t, "+OK\r\n", string(response))
	}
	processingChannel <- utils.MarshalToResp("GET myKey")
	response = <-processingChannel
	require.Equal(t, ":104\r\n", string(response))
}

func TestIncrWithIncorrectValueType(t *testing.T) {
	var response []byte
	processingChannel <- utils.MarshalToResp("SET myKey hello")
	response = <-processingChannel
	require.Equal(t, "+OK\r\n", string(response))

	processingChannel <- utils.MarshalToResp("INCR myKey")
	response = <-processingChannel
	require.Contains(t, string(response), "WRONGTYPE")
}
