package network

import (
	"bytes"
	"io"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockConn struct {
	buffer bytes.Buffer
}

func (c *MockConn) Read(b []byte) (int, error) {
	if c.buffer.Len() > 0 {
		return c.buffer.Read(b)
	}
	return 0, io.EOF
}

func (c *MockConn) Write(b []byte) (int, error) {
	return c.buffer.Write(b)
}

func (c *MockConn) Close() error {
	return nil
}

func (c *MockConn) LocalAddr() net.Addr {
	return nil
}

func (c *MockConn) RemoteAddr() net.Addr {
	return nil
}

func (c *MockConn) SetDeadline(t time.Time) error {
	return nil
}

func (c *MockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (c *MockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	require.Equal(t, config.addr, defaultAddress)
	require.Equal(t, config.port, defaultPort)
	require.Equal(t, config.protocol, defaultProtocol)
}

// func TestHandleConnection(t *testing.T) {
// 	storage := storage.NewStorage()
// 	mockConn := MockConn{}
// 	_, err := mockConn.Write(utils.MarshalToResp("SET masterKey myValue"))
// 	require.NoError(t, err)
// 	handleConnection(&mockConn, resp.NewRespCommandProcessor(storage))

// 	time.Sleep(time.Millisecond * 100)
// 	bytes := make([]byte, 1024)
// 	n, err := mockConn.Read(bytes)
// 	require.NoError(t, err)
// 	require.Greater(t, n, 0)
// 	require.Contains(t, string(bytes[:n]), "+OK\r\n")
// 	mockConn.Close()
// }
