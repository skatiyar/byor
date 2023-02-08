package byor

import (
	"net"
	"sync"
	"time"
)

var slog = newLogger("[SERVER]")

func Server(port string) error {
	var wg sync.WaitGroup
	quit := catchSignals()

	addr, addrErr := net.ResolveTCPAddr("tcp", port)
	if addrErr != nil {
		return addrErr
	}

	tcpListener, listenerErr := net.ListenTCP("tcp", addr)
	if listenerErr != nil {
		return listenerErr
	}
	slog.debugf("Listening on port %v", port)

	for {
		select {
		case <-quit:
			slog.debugf("Got close siginal, waiting for connections to drain")
			tcpListener.Close()
			wg.Wait()
			return nil
		default:
			tcpListener.SetDeadline(time.Now().Add(1e9))
			tcpConn, connErr := tcpListener.AcceptTCP()
			if connErr != nil {
				if opErr, ok := connErr.(*net.OpError); !ok || !opErr.Timeout() {
					slog.errorf("%v", connErr)
				}
				continue
			}
			go func(conn net.Conn) {
				wg.Add(1)
				defer conn.Close()
				defer wg.Done()

				data := make([]byte, 2048)
				br, rErr := conn.Read(data)
				if rErr != nil {
					slog.errorf("Reading from connection %v", rErr)
				}
				slog.debugf("Read %v bytes from connection", br)
				slog.debugf("Message from client %v", string(data))

				wr, wErr := conn.Write([]byte("Hello from server!"))
				if wErr != nil {
					slog.errorf("Writing to connection %v", wErr)
				}
				slog.debugf("Wrote %v bytes to connection", wr)
			}(tcpConn)
		}
	}
}
