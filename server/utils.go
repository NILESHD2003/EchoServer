package server

import "net"

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
