package storage

// The storage does not know anything about RESP so the
// caller must perform the necessary validations before
// using the storage since it only reads and writes bytes.
type StorageEntry struct {
	DataType byte
	Value    []byte
}

func (se StorageEntry) IsNull() bool {
	return se.DataType == 0
}
