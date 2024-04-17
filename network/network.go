// Network layer implementation
package network

import (
	"fmt"
	"log"
	"net"

	"github.com/johanlantz/redis/resp"
)

const defaultPort = 6379
const defaultProtocol = "tcp"
const defaultAddress = "localhost"

type ServerConfig struct {
	addr     string
	port     int
	protocol string
}

// Default server parameters for local testing purposes
func DefaultConfig() ServerConfig {
	return ServerConfig{
		addr:     defaultAddress,
		port:     defaultPort,
		protocol: defaultProtocol,
	}
}

func StartServer(config ServerConfig, storage resp.KVStorage) {
	listener, err := net.Listen(config.protocol, fmt.Sprintf("%s:%d", config.addr, config.port))
	if err != nil {
		log.Panic("Error starting server:", err.Error())
		return
	}
	defer listener.Close()

	processingChannel := make(chan []byte)
	resp.StartCommandProcessor(processingChannel, storage)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err.Error())
			return
		}
		go handleConnection(conn, processingChannel)
	}
}

func handleConnection(conn net.Conn, processingChannel chan []byte) {
	defer conn.Close()

	for {
		bytes := make([]byte, 1024)
		n, err := conn.Read(bytes)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		log.Printf("Received: %s\n", string(bytes[:n]))

		processingChannel <- bytes[:n]
		response := <-processingChannel

		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Println("Error writing:", err.Error())
			break
		}
	}
	log.Printf("closing connection")
}
