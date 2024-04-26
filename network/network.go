// Network layer implementation
package network

import (
	"bytes"
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

	requestChannel := make(chan resp.NetworkRequest)

	resp.StartCommandProcessor(requestChannel, storage)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err.Error())
			return
		}
		go handleConnection(conn, requestChannel)
	}
}

func handleConnection(conn net.Conn, requestChannel chan<- resp.NetworkRequest) {
	defer conn.Close()

	responseChannel := make(chan []byte)

	var buffer bytes.Buffer
	readBuffer := make([]byte, 10) // Only to demonstrate segmentation support
	for {

		n, err := conn.Read(readBuffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		log.Printf("Received: %s\n", string(readBuffer[:n]))

		buffer.Write(readBuffer[:n])

		request := resp.NetworkRequest{ResponseChannel: responseChannel, Data: buffer.Bytes()[:len(buffer.Bytes())]}

		requestChannel <- request
		response := <-responseChannel

		if len(response) > 0 {
			_, err = conn.Write([]byte(response))
			if err != nil {
				log.Println("Error writing:", err.Error())
				break
			}
			buffer.Reset()
		}

	}
	log.Printf("closing connection")
}
