package byor

import (
	"errors"
	"io"
	"net"
	"sync"
	"syscall"
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

				for {
					cmd, cmdErr := Decoder(conn)
					if cmdErr != nil {
						if errors.Is(cmdErr, io.EOF) {
							slog.errorf("Connection closed %v", cmdErr)
							return
						}
						slog.errorf("Reading from connection %v", cmdErr)
						continue
					}
					slog.debugf("Message from client %v", string(cmd))

					if wErr := Encoder(conn, "Hello from server!"); wErr != nil {
						if errors.Is(wErr, syscall.EPIPE) {
							slog.errorf("Connection closed %v", wErr)
							return
						}
						slog.errorf("Writing to connection %v", wErr)
						continue
					}
				}
			}(tcpConn)
		}
	}
}
