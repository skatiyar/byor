package byor

import (
	"log"
	"net"
)

func Client(port string) error {
	addr, addrErr := net.ResolveTCPAddr("tcp", port)
	if addrErr != nil {
		return addrErr
	}

	tcpConn, connErr := net.DialTCP("tcp", nil, addr)
	if connErr != nil {
		return connErr
	}
	defer tcpConn.Close()

	wr, wErr := tcpConn.Write([]byte("Hello from client!"))
	if wErr != nil {
		log.Fatalf("[ERROR] Writing to connection %v\n", wErr)
	}
	log.Printf("[DEBUG] Wrote %v bytes to connection", wr)

	data := make([]byte, 2048)
	br, rErr := tcpConn.Read(data)
	if rErr != nil {
		log.Printf("[ERROR] Reading from connection %v\n", rErr)
	}
	log.Printf("[DEBUG] Read %v bytes from connection\n", br)
	log.Printf("[DEBUG] Message from server %v\n", string(data))

	return nil
}
