package server

import (
	"echoserver/config"
	"io"
	"log"
	"net"
	"strconv"
)

func readIncomingCommand(c net.Conn) (string, error) {
	var buf []byte = make([]byte, 512)
	n, err := c.Read(buf[:])

	if err != nil {
		return "", err
	}
	return string(buf[:n]), nil
}

func respondToClient(c net.Conn, cmd string) error {
	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}
	return nil
}

func RunSynchronousTCPServer(config config.Config) {
	log.Println("Starting Synchronous TCP Server on ", config.Host, ":", config.Port)

	var concurrent_clients int = 0

	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))

	if err != nil {
		log.Fatal("Error Starting TCP Server: ", err)
		panic(err)
	}

	// infinite loop to accept incoming connections
	for {
		c, err := lsnr.Accept()

		if err != nil {
			log.Println("Error Accepting Connection: ", err)
			continue
		}

		concurrent_clients += 1
		log.Println("New Client Connected: ", c.RemoteAddr(), " | Total Clients: ", concurrent_clients)

		// Infinite loop to handle client connection and echo back messages i.e. accept messages from client and send back the same message to client
		for {
			cmd, err := readIncomingCommand(c)

			if err != nil {
				c.Close()
				concurrent_clients -= 1
				log.Println("Client Disconnected: ", c.RemoteAddr(), " | Total Clients: ", concurrent_clients)
				if err == io.EOF {
					break
				}
				log.Println("Error Reading Command from Client: ", err)
			}
			log.Println("Received Command from Client: ", cmd)
			if err = respondToClient(c, cmd); err != nil {
				log.Println("Error Responding to Client: ", err)
				c.Close()
				concurrent_clients -= 1
				log.Println("Client Disconnected: ", c.RemoteAddr(), " | Total Clients: ", concurrent_clients)
				break
			}
		}
	}
}
