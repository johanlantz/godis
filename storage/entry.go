package storage

// Storage is RESP agnostic so the entity works with bytes.
// We keep track of the DataType and the resp layer will
// know how to interpret it when responding to a request.
type Entry struct {
	DataType byte
	Value    []byte
}

func (se Entry) IsNull() bool {
	return se.DataType == 0
}
