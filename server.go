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
			wg.Add(1)
			go func(conn net.Conn) {
				defer conn.Close()
				defer wg.Done()
				if hErr := connHandler(conn); hErr != nil {
					slog.errorf("Connection closed %v", hErr)
					return
				}
			}(tcpConn)
		}
	}
}

func connHandler(conn net.Conn) error {
	for {
		data, dataErr := Decoder(conn)
		if dataErr != nil {
			if errors.Is(dataErr, io.EOF) {
				return dataErr
			}
			slog.errorf("Reading from connection %v", dataErr)
			continue
		}

		status, reply := requestHandler(data)
		res, resErr := composeRes(status, reply)
		if resErr != nil {
			slog.errorf("Composing response %v", resErr)
			res = validResponseOnComposeError()
		}

		if wErr := Encoder(conn, res); wErr != nil {
			if errors.Is(wErr, syscall.EPIPE) {
				return wErr
			}
			slog.errorf("Writing to connection %v", wErr)
			continue
		}
	}
}
