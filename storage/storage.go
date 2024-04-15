package storage

type Storage struct {
	readCh  chan readRequest
	writeCh chan writeRequest
	data    map[string]StorageEntry
}

type readRequest struct {
	key      string
	resultCh chan<- StorageEntry
}

type writeRequest struct {
	key   string
	value StorageEntry
	done  chan<- struct{}
}

func NewStorage() *Storage {
	kv := &Storage{
		readCh:  make(chan readRequest),
		writeCh: make(chan writeRequest),
		data:    make(map[string]StorageEntry),
	}
	go kv.processReadRequests()
	go kv.processWriteRequests()
	return kv
}

func (kv *Storage) processReadRequests() {
	for req := range kv.readCh {
		value := kv.data[req.key]
		req.resultCh <- value
	}
}

func (kv *Storage) processWriteRequests() {
	for req := range kv.writeCh {
		kv.data[req.key] = req.value
		close(req.done)
	}
}

func (kv *Storage) Get(key string) StorageEntry {
	resultCh := make(chan StorageEntry)
	kv.readCh <- readRequest{key: key, resultCh: resultCh}
	return <-resultCh
}

func (kv *Storage) Set(key string, value StorageEntry) {
	done := make(chan struct{})
	kv.writeCh <- writeRequest{key: key, value: value, done: done}
	<-done
}
