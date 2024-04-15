package main

import "github.com/johanlantz/redis/resp"

func main() {

	startServer(defaultConfig(), resp.NewRespCommandProcessor())
}
