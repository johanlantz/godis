package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestEmptyCommand(t *testing.T) {
	_, err := newRespRequest([]byte(""))
	require.Error(t, err)
}

func TestInvalidCommand(t *testing.T) {
	_, err := newRespRequest([]byte("SETI\r\nr"))
	require.Error(t, err)
	_, err = newRespRequest([]byte("\r\nr"))
	require.Error(t, err)
}

func TestGetMissingParamters(t *testing.T) {
	_, err := newRespRequest([]byte("GET\r\nr"))
	require.Error(t, err)

	_, err = newRespRequest([]byte("GET \r\nr"))
	require.Error(t, err)
}

func TestSetMissingParamters(t *testing.T) {
	_, err := newRespRequest([]byte("SET\r\nr"))
	require.Error(t, err)

	_, err = newRespRequest([]byte("SET masterKey\r\nr"))
	require.Error(t, err)
}

func TestBuildGetCommand(t *testing.T) {
	cmd, err := newRespRequest([]byte("GET masterKey\r\n"))
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_GET)
	require.Equal(t, len(cmd.args), 1)
}

func TestBuildSetCommand(t *testing.T) {
	cmd, err := newRespRequest([]byte("SET masterKey abc123\r\n"))
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_SET)
	require.Equal(t, len(cmd.args), 2)
}
