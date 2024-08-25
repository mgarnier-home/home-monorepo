package proxies

import (
	"context"
	"log"
	"net"
	"sync"
	"time"
)

type UDPProxy struct {
	ListenAddr string
	TargetAddr string
	UDPConn    *net.UDPConn
	Cancel     context.CancelFunc
}

func StartUDPProxy(config *UDPProxy, wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	targetAddr, err := net.ResolveUDPAddr("udp", config.TargetAddr)
	if err != nil {
		log.Fatalf("Failed to resolve target address %s: %v", config.TargetAddr, err)
	}

	listenAddr, err := net.ResolveUDPAddr("udp", config.ListenAddr)
	if err != nil {
		log.Fatalf("Failed to resolve listen address %s: %v", config.ListenAddr, err)
	}

	conn, err := net.ListenUDP("udp", listenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", config.ListenAddr, err)
	}
	config.UDPConn = conn
	defer conn.Close()

	log.Printf("UDP proxy listening on %s, forwarding to %s", config.ListenAddr, config.TargetAddr)

	buf := make([]byte, 2048)
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down UDP proxy on %s", config.ListenAddr)
			return
		default:
			conn.SetReadDeadline(time.Now().Add(1 * time.Second)) // Use timeout to allow checking for context cancelation
			n, clientAddr, err := conn.ReadFromUDP(buf)
			if err != nil {
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					continue
				}
				log.Printf("Failed to read from client: %v", err)
				continue
			}

			_, err = conn.WriteToUDP(buf[:n], targetAddr)
			if err != nil {
				log.Printf("Failed to write to target %s: %v", config.TargetAddr, err)
				continue
			}

			n, _, err = conn.ReadFromUDP(buf)
			if err != nil {
				log.Printf("Failed to read from target: %v", err)
				continue
			}

			_, err = conn.WriteToUDP(buf[:n], clientAddr)
			if err != nil {
				log.Printf("Failed to write to client: %v", err)
				continue
			}
		}
	}
}
