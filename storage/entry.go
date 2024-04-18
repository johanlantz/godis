package storage

type Entry struct {
	DataType byte
	Value    []byte
}

func (se Entry) IsNull() bool {
	return se.DataType == 0
}
