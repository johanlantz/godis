package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetWhenNotExisting(t *testing.T) {
	storage := NewChannelStorage()
	v := storage.Get("masterKey")
	require.Condition(t, v.IsNull)
}

func TestSetGetSimpleString(t *testing.T) {
	storage := NewChannelStorage()
	key := "masterKey"
	setValue := []byte("hello")
	dt := byte('+')
	storage.Set(key, Entry{DataType: dt, Value: []byte(setValue)})
	getValue := storage.Get(key)
	require.Equal(t, setValue, getValue.Value)
	require.Equal(t, dt, getValue.DataType)
}
