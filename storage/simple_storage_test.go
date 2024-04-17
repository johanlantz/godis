package storage

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSimpleGetWhenNotExisting(t *testing.T) {
	storage := NewSimpleStorage()
	v := storage.Get("masterKey")
	require.Condition(t, v.IsNull)
}

func TestSimpleSetGetSimpleString(t *testing.T) {
	storage := NewSimpleStorage()
	key := "masterKey"
	setValue := []byte("hello")
	dt := byte('+')
	storage.Set(key, Entry{DataType: dt, Value: []byte(setValue)})
	getValue := storage.Get(key)
	require.Equal(t, setValue, getValue.Value)
	require.Equal(t, dt, getValue.DataType)
}
