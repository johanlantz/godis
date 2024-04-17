package resp

import (
	"testing"

	"github.com/johanlantz/redis/storage"
	"github.com/johanlantz/redis/utils"
	"github.com/stretchr/testify/require"
)

func TestInvalidCommand(t *testing.T) {
	processor := NewRespCommandProcessor(storage.NewStorage())
	response := processor.ProcessCommand(utils.MarshalToResp("SETI"))
	require.Contains(t, string(response), RESP_ERR)
}

func TestGetWithoutKey(t *testing.T) {
	processor := NewRespCommandProcessor(storage.NewStorage())
	response := processor.ProcessCommand(utils.MarshalToResp("GET"))
	require.Contains(t, string(response), RESP_ERR)
}

func TestSetWithoutKey(t *testing.T) {
	processor := NewRespCommandProcessor(storage.NewStorage())
	response := processor.ProcessCommand(utils.MarshalToResp("SET \r\n"))
	require.Contains(t, string(response), RESP_ERR)
}

func TestGetWhenNoValueStored(t *testing.T) {
	processor := NewRespCommandProcessor(storage.NewStorage())
	response := processor.ProcessCommand(utils.MarshalToResp("GET masterKey"))
	require.Contains(t, response, byte(DT_NULLS))
}

func TestSetWithoutValue(t *testing.T) {
	processor := NewRespCommandProcessor(storage.NewStorage())
	response := processor.ProcessCommand(utils.MarshalToResp("SET masterKey"))
	require.Contains(t, string(response), RESP_ERR)
}

func TestSet(t *testing.T) {
	processor := NewRespCommandProcessor(storage.NewStorage())
	response := processor.ProcessCommand(utils.MarshalToResp("SET masterKey myValue"))
	require.Equal(t, "+OK\r\n", string(response))
}

func TestSetGetString(t *testing.T) {
	processor := NewRespCommandProcessor(storage.NewStorage())
	response := processor.ProcessCommand(utils.MarshalToResp("SET masterKey myValue"))
	require.Equal(t, "+OK\r\n", string(response))

	response = processor.ProcessCommand(utils.MarshalToResp("GET masterKey"))
	require.Equal(t, "+myValue\r\n", string(response))
}

func TestInteger(t *testing.T) {
	processor := NewRespCommandProcessor(storage.NewStorage())
	response := processor.ProcessCommand(utils.MarshalToResp("SET myIntCounter 5"))
	require.Equal(t, "+OK\r\n", string(response))

	response = processor.ProcessCommand(utils.MarshalToResp("GET myIntCounter"))
	require.Equal(t, ":5\r\n", string(response))
}

func TestSetFloat(t *testing.T) {
	processor := NewRespCommandProcessor(storage.NewStorage())
	response := processor.ProcessCommand(utils.MarshalToResp("SET myFloatCounter 5.4"))
	require.Equal(t, "+OK\r\n", string(response))

	response = processor.ProcessCommand(utils.MarshalToResp("GET myFloatCounter"))
	require.Equal(t, ",5.4\r\n", string(response))
}

func TestBool(t *testing.T) {
	processor := NewRespCommandProcessor(storage.NewStorage())
	response := processor.ProcessCommand(utils.MarshalToResp("SET myBool true"))
	require.Equal(t, "+OK\r\n", string(response))

	response = processor.ProcessCommand(utils.MarshalToResp("GET myBool"))
	require.Equal(t, "#t\r\n", string(response))

	response = processor.ProcessCommand(utils.MarshalToResp("SET myBool false"))
	require.Equal(t, "+OK\r\n", string(response))

	response = processor.ProcessCommand(utils.MarshalToResp("GET myBool"))
	require.Equal(t, "#f\r\n", string(response))
}
