package server

import (
	"echoserver/config"
	"io"
	"log"
	"net"
	"strconv"
)

func handleClientConnection(c net.Conn, concurrent_clients *int) {
	defer c.Close()

	// Infinite loop to handle client connection and echo back messages i.e. accept messages from client and send back the same message to client
	for {
		cmd, err := readIncomingCommand(c)

		if err != nil {
			*concurrent_clients -= 1

			log.Println("[Asynchronous]Client Disconnected: ", c.RemoteAddr(), " | Total Clients: ", *concurrent_clients)
			if err == io.EOF {
				return
			}
			log.Println("[Asynchronous]Error Reading Command from Client: ", err)
		}
		log.Println("[Asynchronous]Received Command from Client: ", cmd)

		if err := respondToClient(c, cmd); err != nil {
			log.Println("[Asynchronous]Error Responding to Client: ", err)
			*concurrent_clients -= 1
			log.Println("[Asynchronous]Client Disconnected: ", c.RemoteAddr(), " | Total Clients: ", *concurrent_clients)
			return
		}
	}
}

func RunAsynchronousTCPServer(config config.Config) {
	log.Println("[Asynchronous]Starting Asynchronous TCP Server on ", config.Host, ":", config.Port)

	var concurrent_clients int = 0

	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))

	if err != nil {
		log.Fatal("[Asynchronous]Error Starting TCP Server: ", err)
		panic(err)
	}

	// infinite loop to accept incoming connections
	for {
		c, err := lsnr.Accept()

		if err != nil {
			log.Println("[Asynchronous]Error Accepting Connection: ", err)
			continue
		}

		concurrent_clients += 1
		log.Println("[Asynchronous]New Client Connected: ", c.RemoteAddr(), " | Total Clients: ", concurrent_clients)

		// spawing a new goroutine to handle the client connection and echo back messages i.e. accept messages from client and send back the same message to client
		go handleClientConnection(c, &concurrent_clients)
	}
}
