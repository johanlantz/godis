package storage

type Storage struct {
	readCh  chan readRequest
	writeCh chan writeRequest
	data    map[string]Entry
}

type readRequest struct {
	key      string
	resultCh chan<- Entry
}

type writeRequest struct {
	key   string
	value Entry
	done  chan<- struct{}
}

func NewStorage() *Storage {
	kv := &Storage{
		readCh:  make(chan readRequest),
		writeCh: make(chan writeRequest),
		data:    make(map[string]Entry),
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

func (kv *Storage) Get(key string) Entry {
	resultCh := make(chan Entry)
	kv.readCh <- readRequest{key: key, resultCh: resultCh}
	return <-resultCh
}

func (kv *Storage) Set(key string, value Entry) {
	done := make(chan struct{})
	kv.writeCh <- writeRequest{key: key, value: value, done: done}
	<-done
}
