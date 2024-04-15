package resp

import (
	"testing"

	"github.com/johanlantz/redis/utils"
	"github.com/stretchr/testify/require"
)

func TestMalformedGet(t *testing.T) {
	_, err := newRespRequest(utils.MarshalToResp("GET"))
	require.NoError(t, err)

	_, err = newRespRequest(utils.MarshalToResp("GETmasterKey"))
	require.Error(t, err)
}

func TestMalformedSet(t *testing.T) {
	_, err := newRespRequest(utils.MarshalToResp("SETmasterKey value"))
	require.Error(t, err)
}

func TestBuildGetCommand(t *testing.T) {
	cmd, err := newRespRequest(utils.MarshalToResp("GET"))
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_GET)
	require.Equal(t, len(cmd.args), 0)

	cmd, err = newRespRequest(utils.MarshalToResp("GET masterKey\r\n"))
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_GET)
	require.Equal(t, len(cmd.args), 1)

	cmd, err = newRespRequest(utils.MarshalToResp("GET    masterKey    \r\n"))
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_GET)
	require.Equal(t, len(cmd.args), 1)
}

func TestBuildSetCommand(t *testing.T) {
	cmd, err := newRespRequest(utils.MarshalToResp("SET masterKey abc123\r\n"))
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_SET)
	require.Equal(t, len(cmd.args), 2)
}
