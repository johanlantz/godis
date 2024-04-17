// Simplest storage possible but it requires on someone else guaranteeing
// sequential execution. In the current solution this is managed by the
// processor.
package storage

type SimpleStorage struct {
	data map[string]Entry
}

func NewSimpleStorage() *SimpleStorage {
	return &SimpleStorage{data: make(map[string]Entry)}
}

func (kv *SimpleStorage) Get(key string) Entry {
	return kv.data[key]
}

func (kv *SimpleStorage) Set(key string, value Entry) {
	kv.data[key] = value
}
