package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetWhenNotExisting(t *testing.T) {
	v := Get("masterKey")
	require.Condition(t, v.IsNull)
}

func TestSetGetSimpleString(t *testing.T) {
	key := "masterKey"
	setValue := []byte("hello")
	dt := byte('+')
	Set(key, StorageEntry{DataType: dt, Value: []byte(setValue)})
	getValue := Get(key)
	require.Equal(t, setValue, getValue.Value)
	require.Equal(t, dt, getValue.DataType)
}
