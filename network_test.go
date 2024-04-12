package main

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type MockConn struct {
	buffer bytes.Buffer
}

func (c *MockConn) Read(b []byte) (int, error) {
	return c.buffer.Read(b)
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
	config := defaultConfig()
	require.Equal(t, config.addr, defaultAddress)
	require.Equal(t, config.port, defaultPort)
	require.Equal(t, config.protocol, defaultProtocol)
}

func TestHandleConnection(t *testing.T) {
	mockConn := MockConn{}
	_, err := mockConn.Write([]byte("Hello"))
	require.NoError(t, err)
	handleConnection(&mockConn)
}
