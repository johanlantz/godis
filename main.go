package main

import (
	"github.com/johanlantz/redis/network"
	"github.com/johanlantz/redis/resp"
)

func main() {

	network.StartServer(network.DefaultConfig(), resp.NewRespCommandProcessor())
}
