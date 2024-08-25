package proxies

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/utils"
	"net"
	"sync"
)

type TCPProxy struct {
	ListenAddr     string
	TargetAddr     string
	HostStarted    func() (bool, error)
	StartHost      func() error
	PacketReceived func() error

	stop context.CancelFunc
}

type CustomWriter struct {
	io.Writer
	onWrite func(int)
}

func (cw *CustomWriter) Write(p []byte) (int, error) {
	n, err := cw.Writer.Write(p)
	if err == nil {
		cw.onWrite(n)
	}
	return n, err
}

func NewTCPProxy(hostConfig *config.HostConfig, proxyConfig *config.ProxyConfig, hostStarted func() (bool, error), startHost func() error, packetReceived func() error) *TCPProxy {
	tcpProxy := &TCPProxy{
		ListenAddr:     fmt.Sprintf("%s:%d", "0.0.0.0", proxyConfig.ListenPort),
		TargetAddr:     fmt.Sprintf("%s:%d", hostConfig.Ip, proxyConfig.TargetPort),
		HostStarted:    hostStarted,
		StartHost:      startHost,
		PacketReceived: packetReceived,
	}

	return tcpProxy
}

func (proxy *TCPProxy) Start(wg *sync.WaitGroup, ctx context.Context, stop context.CancelFunc) {
	proxy.stop = stop
	defer wg.Done()

	listener, err := net.Listen("tcp", proxy.ListenAddr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", proxy.ListenAddr, err)
	}
	defer listener.Close()

	log.Printf("TCP proxy listening on %s, forwarding to %s", proxy.ListenAddr, proxy.TargetAddr)

	stopChan := make(chan struct{})
	go func() {
		<-ctx.Done()
		close(stopChan)
		listener.Close() // This will unblock the listener.Accept() call
	}()

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			select {
			case <-stopChan:
				log.Printf("Shutting down TCP proxy on %s", proxy.ListenAddr)
				return
			default:
				log.Printf("Failed to accept connection on %s: %v", proxy.ListenAddr, err)
				continue
			}
		}
		log.Printf("Accepted connection from %s", clientConn.RemoteAddr())

		go proxy.handleTCPConnection(clientConn, proxy.TargetAddr)
	}
}

func (proxy *TCPProxy) Stop() {
	proxy.stop()
}

func (proxy *TCPProxy) shouldForwardProxy(clientConn net.Conn) (bool, error) {
	started, err := proxy.HostStarted()
	if err != nil {
		return false, err
	}

	if !started {
		reader := bufio.NewReader(clientConn)
		peek, err := reader.Peek(reader.Buffered())

		if err != nil {
			return false, err
		}

		if utils.IsHTTPRequest(peek) {
			request := string(peek)

			if utils.CheckRequestHeader(request, "Status", "true") {
				clientConn.Close()
				return false, nil
			}
		}

		err = proxy.StartHost()

		if err != nil {
			return false, err
		}
	}

	return true, nil
}

func (proxy *TCPProxy) handleTCPConnection(clientConn net.Conn, targetAddr string) {
	defer clientConn.Close()

	forwardProxy, err := proxy.shouldForwardProxy(clientConn)

	if err != nil {
		log.Printf("Failed to check host status: %v", err)
		return
	}

	if !forwardProxy {
		log.Printf("Dropping connection")
		return
	}

	log.Printf("Forwarding connection to %s", proxy.TargetAddr)

	targetConn, err := net.Dial("tcp", targetAddr)
	if err != nil {
		log.Printf("Failed to connect to target %s: %v", targetAddr, err)
		return
	}
	defer targetConn.Close()

	done := make(chan struct{})

	onDataTransfer := func(bytesTransferred int) {
		log.Printf("Data transferred: %d bytes", bytesTransferred)

		proxy.PacketReceived()
	}

	clientToTargetWriter := &CustomWriter{Writer: targetConn, onWrite: onDataTransfer}

	go func() {
		defer close(done)
		_, err := io.Copy(clientToTargetWriter, clientConn)
		if err != nil {
			log.Printf("Error copying from client to target: %v", err)
		}
	}()

	_, err = io.Copy(clientConn, targetConn)
	if err != nil {
		log.Printf("Error copying from target to client: %v", err)
	}

	<-done

	log.Println("Closing TCP connection")

}
