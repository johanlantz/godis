package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInvalidCommand(t *testing.T) {
	processor := newRespCommandProcessor()
	response := processor.processCommand([]byte("SETI\r\nr"))
	require.Contains(t, response, "ERR")
}
