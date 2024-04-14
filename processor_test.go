package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvalidCommand(t *testing.T) {
	processor := newRespCommandProcessor()
	response := processor.processCommand([]byte("SETI\r\n"))
	require.Contains(t, string(response), RESP_ERR)
}

func TestGetWithoutKey(t *testing.T) {
	processor := newRespCommandProcessor()
	response := processor.processCommand([]byte("GET \r\n"))
	require.Contains(t, string(response), RESP_ERR)
}

func TestSetWithoutKey(t *testing.T) {
	processor := newRespCommandProcessor()
	response := processor.processCommand([]byte("SET \r\n"))
	require.Contains(t, string(response), RESP_ERR)
}

func TestSetWithoutValue(t *testing.T) {
	processor := newRespCommandProcessor()
	response := processor.processCommand([]byte("SET masterKey\r\n"))
	require.Contains(t, string(response), RESP_ERR)
}

func TestSet(t *testing.T) {
	processor := newRespCommandProcessor()
	response := processor.processCommand([]byte("SET masterKey myValue\r\n"))
	require.Equal(t, "+OK\r\n", string(response))
}

func TestSetGet(t *testing.T) {
	processor := newRespCommandProcessor()
	response := processor.processCommand([]byte("SET masterKey myValue\r\n"))
	require.Equal(t, "+OK\r\n", string(response))

	response = processor.processCommand([]byte("GET masterKey\r\n"))
	require.Equal(t, "+myValue\r\n", string(response))
}
