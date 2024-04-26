// Simplest storage possible but it requires someone else guaranteeing
// sequential access. In the current solution this is managed by the
// resp processor so this is ok.
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

func (kv *SimpleStorage) Delete(key string) {
	delete(kv.data, key)
}
