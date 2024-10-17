package proxies

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/utils"
	"net"
	"sync"

	"github.com/charmbracelet/log"
)

type TCPProxy struct {
	Name           string
	ListenAddr     *net.TCPAddr
	TargetAddr     *net.TCPAddr
	HostStarted    func(proxy *TCPProxy) (bool, error)
	StartHost      func(proxy *TCPProxy) error
	PacketReceived func(proxy *TCPProxy) error

	ctx    context.Context
	cancel context.CancelFunc
}

type TCPProxyArgs struct {
	HostIp         string
	ProxyConfig    *config.ProxyConfig
	HostStarted    func(proxy *TCPProxy) (bool, error)
	StartHost      func(proxy *TCPProxy) error
	PacketReceived func(proxy *TCPProxy) error
}

func NewTCPProxy(args *TCPProxyArgs) *TCPProxy {
	ctx, cancel := context.WithCancel(context.Background())

	listenAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", args.ProxyConfig.ListenPort))
	if err != nil {
		log.Errorf("%s: Failed to resolve listen TCP address %d: %v", args.ProxyConfig.Name, args.ProxyConfig.ListenPort, err)
		panic(err)
	}

	targetAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", args.HostIp, args.ProxyConfig.TargetPort))
	if err != nil {
		log.Errorf("%s: Failed to resolve target TCP address %d: %v", args.ProxyConfig.Name, args.ProxyConfig.TargetPort, err)
		panic(err)
	}

	tcpProxy := &TCPProxy{
		Name:           args.ProxyConfig.Name,
		ListenAddr:     listenAddr,
		TargetAddr:     targetAddr,
		HostStarted:    args.HostStarted,
		StartHost:      args.StartHost,
		PacketReceived: args.PacketReceived,
		ctx:            ctx,
		cancel:         cancel,
	}

	log.Infof("%s: TCP Proxy created: %s -> %s", tcpProxy.Name, tcpProxy.ListenAddr, tcpProxy.TargetAddr)

	return tcpProxy
}

func (proxy *TCPProxy) Start(wg *sync.WaitGroup) {
	defer func() {
		wg.Done()
		log.Infof("%s: TCP proxy stopped", proxy.Name)
	}()

	listener, err := net.ListenTCP("tcp", proxy.ListenAddr)
	if err != nil {
		log.Errorf("%s: Failed to start TCP proxy: %v", proxy.Name, err)
		return
	}
	defer listener.Close()

	log.Debugf("%s: TCP proxy started on %s", proxy.Name, proxy.ListenAddr)

	stopChan := make(chan struct{})
	go func() {
		<-proxy.ctx.Done()
		log.Infof("%s: Stopping TCP proxy on %s", proxy.Name, proxy.ListenAddr)
		close(stopChan)
		listener.Close() // This will unblock the listener.Accept() call
	}()

	for {
		clientConn, err := listener.AcceptTCP()
		if err != nil {
			select {
			case <-stopChan:
				log.Debugf("%s: Stopped accepting connections on %s", proxy.Name, proxy.ListenAddr)
				return
			default:
				log.Errorf("%s: Failed to accept connection: %v", proxy.Name, err)
				continue
			}
		}
		log.Debugf("%s: Accepted connection from %s", proxy.Name, clientConn.RemoteAddr())

		wg.Add(1)
		go proxy.handleTCPConnection(wg, clientConn)
	}
}

func (proxy *TCPProxy) Stop() {
	log.Infof("%s: Stopping TCP proxy", proxy.Name)
	proxy.cancel()
}

func (proxy *TCPProxy) shouldForwardProxy(clientConn *net.TCPConn) (bool, error) {
	if proxy.HostStarted == nil {
		return false, fmt.Errorf("%s: HostStarted function not set", proxy.Name)
	}

	started, err := proxy.HostStarted(proxy)
	if err != nil {
		return false, err
	}

	if !started {
		reader := bufio.NewReader(clientConn)
		peek, err := reader.Peek(utils.Min(1024, reader.Buffered()))

		if err != nil && err != io.EOF {
			return false, fmt.Errorf("%s: Failed to peek data: %v", proxy.Name, err)
		}

		if utils.IsHTTPRequest(peek) {
			request := string(peek)

			if utils.CheckRequestHeader(request, "Status", "true") {
				clientConn.Close()
				return false, nil
			}
		}

		if proxy.StartHost == nil {
			return false, fmt.Errorf("%s: StartHost function not set", proxy.Name)
		}

		err = proxy.StartHost(proxy)

		if err != nil {
			return false, fmt.Errorf("%s: Failed to start host: %v", proxy.Name, err)
		}
	}

	return true, nil
}

func (proxy *TCPProxy) handleTCPConnection(wg *sync.WaitGroup, clientConn *net.TCPConn) {
	defer wg.Done()

	forwardProxy, err := proxy.shouldForwardProxy(clientConn)

	if err != nil {
		log.Errorf("%s: Failed to determine if proxy should be forwarded: %v", proxy.Name, err)
		return
	}

	if !forwardProxy {
		log.Infof("%s: Proxy not forwarded to target", proxy.Name)
		return
	}

	log.Infof("%s: Proxy forwarded to target %s", proxy.Name, proxy.TargetAddr)

	targetConn, err := net.DialTCP("tcp", nil, proxy.TargetAddr)
	if err != nil {
		log.Errorf("%s: Failed to connect to target: %v", proxy.Name, err)
		return
	}
	defer targetConn.Close()

	onClientToTarget := func(bytesTransferred int) {
		log.Debugf("%s: ClientToTarget: %d bytes", proxy.Name, bytesTransferred)

		if proxy.PacketReceived == nil {
			log.Debugf("%s: PacketReceived function not set", proxy.Name)
			return
		}

		proxy.PacketReceived(proxy)

	}

	onTargetToClient := func(bytesTransferred int) {
		log.Debugf("%s: TargetToClient: %d bytes", proxy.Name, bytesTransferred)

	}

	clientToTargetWriter := &utils.CustomWriter{Writer: targetConn, OnWrite: onClientToTarget}
	targetToClientWriter := &utils.CustomWriter{Writer: clientConn, OnWrite: onTargetToClient}

	doneCopyClientToTarget := make(chan struct{})
	go func() {
		defer close(doneCopyClientToTarget)
		_, err := io.Copy(clientToTargetWriter, clientConn)
		if err != nil {
			log.Debugf("%s: Error copying from client to target: %v", proxy.Name, err)
		}
	}()

	doneCopyTargetToClient := make(chan struct{})
	go func() {
		defer close(doneCopyTargetToClient)
		_, err = io.Copy(targetToClientWriter, targetConn)
		if err != nil {
			log.Debugf("%s: Error copying from target to client: %v", proxy.Name, err)
		}
	}()

	select {
	case <-proxy.ctx.Done():
		log.Infof("%s: Context done", proxy.Name)
		clientConn.Close()
		targetConn.Close()
	case <-doneCopyClientToTarget:
		targetConn.Close()
		log.Infof("%s: Client to target copy done", proxy.Name)
	case <-doneCopyTargetToClient:
		log.Infof("%s: Target to client copy done", proxy.Name)
		clientConn.Close()
	}

	log.Infof("%s: Connection closed", proxy.Name)

}
