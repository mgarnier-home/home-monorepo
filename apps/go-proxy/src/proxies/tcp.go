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
	"time"
)

type TCPProxy struct {
	ListenAddr     *net.TCPAddr
	TargetAddr     *net.TCPAddr
	HostStarted    func(proxy *TCPProxy) (bool, error)
	StartHost      func(proxy *TCPProxy) error
	PacketReceived func(proxy *TCPProxy) error

	ctx    context.Context
	cancel context.CancelFunc
}

type TCPProxyArgs struct {
	HostConfig     *config.HostConfig
	ProxyConfig    *config.ProxyConfig
	HostStarted    func(proxy *TCPProxy) (bool, error)
	StartHost      func(proxy *TCPProxy) error
	PacketReceived func(proxy *TCPProxy) error
}

func NewTCPProxy(ctx context.Context, args *TCPProxyArgs) *TCPProxy {
	ctx, cancel := context.WithCancel(ctx)

	listenAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", args.ProxyConfig.ListenPort))
	if err != nil {
		log.Printf("Failed to resolve listen TCP address %d: %v", args.ProxyConfig.ListenPort, err)
		panic(err)
	}

	targetAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", args.HostConfig.Ip, args.ProxyConfig.TargetPort))
	if err != nil {
		log.Printf("Failed to resolve target TCP address %d: %v", args.ProxyConfig.TargetPort, err)
		panic(err)
	}

	tcpProxy := &TCPProxy{
		ListenAddr:     listenAddr,
		TargetAddr:     targetAddr,
		HostStarted:    args.HostStarted,
		StartHost:      args.StartHost,
		PacketReceived: args.PacketReceived,
		ctx:            ctx,
		cancel:         cancel,
	}

	log.Printf("TCP Proxy created: %s -> %s", tcpProxy.ListenAddr, tcpProxy.TargetAddr)

	return tcpProxy
}

func (proxy *TCPProxy) Start(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		log.Printf("TCP proxy stopped")
	}()

	listener, err := net.ListenTCP("tcp", proxy.ListenAddr)
	if err != nil {
		log.Printf("Failed to listen on %s: %v", proxy.ListenAddr, err)
		return
	}
	defer listener.Close()

	log.Printf("TCP proxy listening on %s, forwarding to %s", proxy.ListenAddr, proxy.TargetAddr)

	stopChan := make(chan struct{})
	go func() {
		<-proxy.ctx.Done()
		log.Printf("Context done")
		close(stopChan)
		listener.Close() // This will unblock the listener.Accept() call
	}()

	for {
		clientConn, err := listener.AcceptTCP()
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

		wg.Add(1)
		go proxy.handleTCPConnection(wg, clientConn)
	}
}

func (proxy *TCPProxy) Stop() {
	log.Printf("Stop called on TCP proxy")
	proxy.cancel()
}

func (proxy *TCPProxy) shouldForwardProxy(clientConn *net.TCPConn) (bool, error) {
	if proxy.HostStarted == nil {
		return false, fmt.Errorf("HostStarted function not set")
	}

	started, err := proxy.HostStarted(proxy)
	if err != nil {
		return false, err
	}

	if !started {
		reader := bufio.NewReader(clientConn)
		peek, err := reader.Peek(utils.Min(1024, reader.Buffered()))

		if err != nil && err != io.EOF {
			return false, fmt.Errorf("failed to peek data: %w", err)
		}

		if utils.IsHTTPRequest(peek) {
			request := string(peek)

			if utils.CheckRequestHeader(request, "Status", "true") {
				clientConn.Close()
				return false, nil
			}
		}

		if proxy.StartHost == nil {
			return false, fmt.Errorf("StartHost function not set")
		}

		err = proxy.StartHost(proxy)

		if err != nil {
			return false, fmt.Errorf("failed to start host: %w", err)
		}
	}

	return true, nil
}

func (proxy *TCPProxy) handleTCPConnection(wg *sync.WaitGroup, clientConn *net.TCPConn) {
	defer wg.Done()
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

	targetConn, err := net.DialTCP("tcp", nil, proxy.TargetAddr)
	if err != nil {
		log.Printf("Failed to connect to target %s: %v", proxy.TargetAddr, err)
		return
	}
	defer targetConn.Close()

	// targetConn.SetKeepAlive(true)

	// targetConn.SetLinger(0)
	targetConn.SetKeepAliveConfig(net.KeepAliveConfig{
		Enable:   true,
		Idle:     5 * time.Second,
		Interval: 1 * time.Second,
		Count:    3,
	})

	clientConn.SetKeepAliveConfig(net.KeepAliveConfig{
		Enable:   true,
		Idle:     5 * time.Second,
		Interval: 1 * time.Second,
		Count:    3,
	})
	// clientConn.SetLinger(0)

	onClientToTarget := func(bytesTransferred int) {
		// log.Printf("ClientToTarget: %d bytes", bytesTransferred)

		if proxy.PacketReceived == nil {
			log.Printf("Error calling PacketReceived")
			return
		}

		proxy.PacketReceived(proxy)

	}

	onTargetToClient := func(bytesTransferred int) {
		// log.Printf("TargetToClient: %d bytes", bytesTransferred)

	}

	clientToTargetWriter := &utils.CustomWriter{Writer: targetConn, OnWrite: onClientToTarget}
	targetToClientWriter := &utils.CustomWriter{Writer: clientConn, OnWrite: onTargetToClient}

	doneCopyClientToTarget := make(chan struct{})
	go func() {
		defer close(doneCopyClientToTarget)
		_, err := io.Copy(clientToTargetWriter, clientConn)
		if err != nil {
			log.Printf("Error copying from client to target: %v", err)
		}

		log.Printf("Client to target copy done")
	}()

	doneCopyTargetToClient := make(chan struct{})
	go func() {
		defer close(doneCopyTargetToClient)
		_, err = io.Copy(targetToClientWriter, targetConn)
		if err != nil {
			log.Printf("Error copying from target to client: %v", err)
		}

		log.Printf("Target to client copy done")
	}()

	select {
	case <-proxy.ctx.Done():
		log.Printf("Context done, force closing open connections")
		clientConn.Close()
		targetConn.Close()
	case <-doneCopyClientToTarget:
		targetConn.Close()
		log.Printf("Client to target copy done")
	case <-doneCopyTargetToClient:
		log.Printf("Target to client copy done")
	}

	log.Println("Handle TCP connections exiting")

}
