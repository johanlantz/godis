package resp

import (
	"testing"

	"github.com/johanlantz/redis/utils"
	"github.com/stretchr/testify/require"
)

func TestInvalidCommand(t *testing.T) {
	processor := NewRespCommandProcessor()
	response := processor.ProcessCommand(utils.MarshalToResp("SETI"))
	require.Contains(t, string(response), RESP_ERR)
}

func TestGetWithoutKey(t *testing.T) {
	processor := NewRespCommandProcessor()
	response := processor.ProcessCommand(utils.MarshalToResp("GET"))
	require.Contains(t, string(response), RESP_ERR)
}

func TestSetWithoutKey(t *testing.T) {
	processor := NewRespCommandProcessor()
	response := processor.ProcessCommand(utils.MarshalToResp("SET \r\n"))
	require.Contains(t, string(response), RESP_ERR)
}

func TestSetWithoutValue(t *testing.T) {
	processor := NewRespCommandProcessor()
	response := processor.ProcessCommand(utils.MarshalToResp("SET masterKey"))
	require.Contains(t, string(response), RESP_ERR)
}

func TestSet(t *testing.T) {
	processor := NewRespCommandProcessor()
	response := processor.ProcessCommand(utils.MarshalToResp("SET masterKey myValue"))
	require.Equal(t, "+OK\r\n", string(response))
}

func TestSetGet(t *testing.T) {
	processor := NewRespCommandProcessor()
	response := processor.ProcessCommand(utils.MarshalToResp("SET masterKey myValue"))
	require.Equal(t, "+OK\r\n", string(response))

	response = processor.ProcessCommand(utils.MarshalToResp("GET masterKey"))
	require.Equal(t, "+myValue\r\n", string(response))
}
