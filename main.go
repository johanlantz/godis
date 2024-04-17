package main

import (
	"github.com/johanlantz/redis/network"
	"github.com/johanlantz/redis/resp"
	"github.com/johanlantz/redis/storage"
)

func main() {

	storage := storage.NewStorage()
	network.StartServer(network.DefaultConfig(), resp.NewRespCommandProcessor(storage))
}
