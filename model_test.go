package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMalformedGet(t *testing.T) {
	_, err := newRespRequest([]byte("GET"))
	require.Error(t, err)

	_, err = newRespRequest([]byte("GET \r\n "))
	require.Error(t, err)

	_, err = newRespRequest([]byte("GET \n"))
	require.Error(t, err)

	_, err = newRespRequest([]byte("GETmasterKey\r\n"))
	require.Error(t, err)
}

func TestMalformedSet(t *testing.T) {
	_, err := newRespRequest([]byte("SET masterKey"))
	require.Error(t, err)
}

func TestBuildGetCommand(t *testing.T) {
	// This is actually valid from the model perspective.
	// It will fails during processing.
	cmd, err := newRespRequest([]byte("GET \r\n"))
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_GET)
	require.Equal(t, len(cmd.args), 1)

	cmd, err = newRespRequest([]byte("GET masterKey\r\n"))
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
