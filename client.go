package byor

import (
	"net"
)

var clog = newLogger("[CLIENT]")

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
		clog.errorf("Writing to connection %v\n", wErr)
	}
	clog.debugf("Wrote %v bytes to connection", wr)

	data := make([]byte, 2048)
	br, rErr := tcpConn.Read(data)
	if rErr != nil {
		clog.errorf("Reading from connection %v", rErr)
	}
	clog.debugf("Read %v bytes from connection", br)
	clog.debugf("Message from server %v", string(data))

	return nil
}
