// A sample storage implementation with different channels for read and write.
// Not used anymore since the processor handles things sequentially now.
package storage

type ChannelStorage struct {
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

func NewChannelStorage() *ChannelStorage {
	kv := &ChannelStorage{
		readCh:  make(chan readRequest),
		writeCh: make(chan writeRequest),
		data:    make(map[string]Entry),
	}
	go kv.processReadRequests()
	go kv.processWriteRequests()
	return kv
}

func (kv *ChannelStorage) processReadRequests() {
	for req := range kv.readCh {
		value := kv.data[req.key]
		req.resultCh <- value
	}
}

func (kv *ChannelStorage) processWriteRequests() {
	for req := range kv.writeCh {
		kv.data[req.key] = req.value
		close(req.done)
	}
}

func (kv *ChannelStorage) Get(key string) Entry {
	resultCh := make(chan Entry)
	kv.readCh <- readRequest{key: key, resultCh: resultCh}
	return <-resultCh
}

func (kv *ChannelStorage) Set(key string, value Entry) {
	done := make(chan struct{})
	kv.writeCh <- writeRequest{key: key, value: value, done: done}
	<-done
}
