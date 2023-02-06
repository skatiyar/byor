package byor

import (
	"log"
	"net"
	"sync"
	"time"
)

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
	log.Printf("[DEBUG] Listening on port %v", port)

	for {
		select {
		case <-quit:
			log.Printf("[DEBUG] Got close siginal, waiting for connections to drain")
			tcpListener.Close()
			wg.Wait()
			return nil
		default:
			tcpListener.SetDeadline(time.Now().Add(1e9))
			tcpConn, connErr := tcpListener.AcceptTCP()
			if connErr != nil {
				if opErr, ok := connErr.(*net.OpError); !ok || !opErr.Timeout() {
					log.Printf("[ERROR] %v\n", connErr)
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
					log.Printf("[ERROR] Reading from connection %v\n", rErr)
				}
				log.Printf("[DEBUG] Read %v bytes from connection\n", br)
				log.Printf("[DEBUG] Message from client %v\n", string(data))

				wr, wErr := conn.Write([]byte("Hello from server!"))
				if wErr != nil {
					log.Printf("[ERROR] Writing to connection %v\n", wErr)
				}
				log.Printf("[DEBUG] Wrote %v bytes to connection", wr)
			}(tcpConn)
		}
	}
}
