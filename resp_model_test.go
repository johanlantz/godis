package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// func TestInvalidCommand(t *testing.T) {
// 	_, err := newRespRequest([]byte("SETI\r\nr"))
// 	require.Error(t, err)
// 	_, err = newRespRequest([]byte("\r\nr"))
// 	require.Error(t, err)
// }

func TestGetMissingParamters(t *testing.T) {
	_, err := newRespRequest([]string{"GET\r\nr"})
	require.Error(t, err)

	_, err = newRespRequest([]string{"GET ", "\r\nr"})
	require.Error(t, err)
}

func TestSetMissingParamters(t *testing.T) {
	_, err := newRespRequest([]string{"SET\r\nr"})
	require.Error(t, err)

	_, err = newRespRequest([]string{"SET ", "masterKey\r\nr"})
	require.Error(t, err)
}

func TestBuildGetCommand(t *testing.T) {
	cmd, err := newRespRequest([]string{"GET ", "masterKey\r\n"})
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_GET)
	require.Equal(t, len(cmd.args), 1)
}

func TestBuildSetCommand(t *testing.T) {
	cmd, err := newRespRequest([]string{"SET ", "masterKey", "abc123\r\n"})
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_SET)
	require.Equal(t, len(cmd.args), 2)
}
