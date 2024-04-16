// Network layer implementation
package network

import (
	"fmt"
	"log"
	"net"
)

// The network layer is only concerned about managing connections and
// bytestreams, it has no notion of what RESP is.
type CommandProcessor interface {
	ProcessCommand(data []byte) []byte
}

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

func StartServer(config ServerConfig, cmdProc CommandProcessor) {
	listener, err := net.Listen(config.protocol, fmt.Sprintf("%s:%d", config.addr, config.port))
	if err != nil {
		log.Panic("Error starting server:", err.Error())
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err.Error())
			return
		}
		go handleConnection(conn, cmdProc)
	}
}

func handleConnection(conn net.Conn, cmdProc CommandProcessor) {
	defer conn.Close()

	for {
		bytes := make([]byte, 1024)
		n, err := conn.Read(bytes)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		log.Printf("Received: %s\n", string(bytes[:n]))

		response := cmdProc.ProcessCommand(bytes[:n])

		_, err = conn.Write([]byte(response))
		if err != nil {
			log.Println("Error writing:", err.Error())
			break
		}
	}
	log.Printf("closing connection")
}
