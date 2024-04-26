package resp

import (
	"testing"

	"github.com/johanlantz/redis/utils"
	"github.com/stretchr/testify/require"
)

func TestMalformedGet(t *testing.T) {

	_, err := newRespRequest(utils.MarshalToResp("GET"), &processors)
	require.NoError(t, err)

	_, err = newRespRequest(utils.MarshalToResp("GETmasterKey"), &processors)
	require.Error(t, err)
}

func TestMalformedSet(t *testing.T) {
	_, err := newRespRequest(utils.MarshalToResp("SETmasterKey value"), &processors)
	require.Error(t, err)
}

func TestBuildGetCommand(t *testing.T) {
	cmd, err := newRespRequest(utils.MarshalToResp("GET"), &processors)
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_GET)
	require.Equal(t, len(cmd.args), 0)

	cmd, err = newRespRequest(utils.MarshalToResp("GET masterKey\r\n"), &processors)
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_GET)
	require.Equal(t, len(cmd.args), 1)

	cmd, err = newRespRequest(utils.MarshalToResp("GET    masterKey    \r\n"), &processors)
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_GET)
	require.Equal(t, len(cmd.args), 1)
}

func TestBuildSetCommand(t *testing.T) {
	cmd, err := newRespRequest(utils.MarshalToResp("SET masterKey abc123\r\n"), &processors)
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_SET)
	require.Equal(t, len(cmd.args), 2)
}

func TestSetCommandWithQuotes(t *testing.T) {
	cmd, err := newRespRequest(utils.MarshalToResp("SET masterKey \"abc123\"\r\n"), &processors)
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_SET)
	require.Equal(t, len(cmd.args), 2)
}

func TestSetCommandWithQuotesAndSpaces(t *testing.T) {
	cmd, err := newRespRequest([]byte("*3\r\n$3\r\nSET\r\n$9\r\nmasterKey\r\n$7\r\nabc 123\r\n"), &processors)
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_SET)
	require.Equal(t, len(cmd.args), 2)
}

func TestSetCommandWithNewlinesAndDollarSign(t *testing.T) {
	cmd, err := newRespRequest([]byte("*3\r\n$3\r\nSET\r\n$2\r\ngg\r\n$13\r\nmy\\r\\n$12\\r\\n\r\n"), &processors)
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_SET)
	require.Equal(t, len(cmd.args), 2)
}

func TestQuotesSetCommandWithNewlinesAndDollarSign(t *testing.T) {
	cmd, err := newRespRequest([]byte("*3\r\n$3\r\nSET\r\n$2\r\ngg\r\n$13\r\nmy\\r\\n$12\\r\\n\r\n"), &processors)
	require.NoError(t, err)
	require.Equal(t, cmd.command, RESP_SET)
	require.Equal(t, len(cmd.args), 2)
}
