package main

import (
	"github.com/johanlantz/redis/network"
	"github.com/johanlantz/redis/storage"
)

func main() {
	storage := storage.NewSimpleStorage()
	network.StartServer(network.DefaultConfig(), storage)
}
