package proxies

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/hostState"
	"mgarnier11/go/colors"
	"mgarnier11/go/logger"
	"mgarnier11/go/utils"
	"net"
	"strings"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type TCPProxy struct {
	Name           string
	ListenAddr     *net.TCPAddr
	ServerAddr     *net.TCPAddr
	StartHost      func() error
	PacketReceived func(proxyName string)

	logger *logger.Logger

	hostState *hostState.State
	wg        sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

type TCPProxyArgs struct {
	HostIp         string
	ProxyConfig    *config.ProxyConfig
	HostState      *hostState.State
	StartHost      func() error
	PacketReceived func(proxyName string)
}

func NewTCPProxy(args *TCPProxyArgs, hostLogger *logger.Logger) *TCPProxy {
	ctx, cancel := context.WithCancel(context.Background())

	logger := logger.NewLogger(fmt.Sprintf("[%s]", strings.ToUpper(args.ProxyConfig.Name)), "%-15s ", lipgloss.NewStyle().Foreground(lipgloss.Color(colors.GenerateHexColor(args.ProxyConfig.Name))), hostLogger)

	listenAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", args.ProxyConfig.ListenPort))
	if err != nil {
		logger.Errorf("Failed to resolve listen TCP address %d: %v", args.ProxyConfig.ListenPort, err)
		panic(err)
	}

	serverAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", args.HostIp, args.ProxyConfig.ServerPort))
	if err != nil {
		logger.Errorf("Failed to resolve server TCP address %d: %v", args.ProxyConfig.ServerPort, err)
		panic(err)
	}

	tcpProxy := &TCPProxy{
		Name:           args.ProxyConfig.Name,
		ListenAddr:     listenAddr,
		ServerAddr:     serverAddr,
		StartHost:      args.StartHost,
		PacketReceived: args.PacketReceived,
		logger:         logger,
		hostState:      args.HostState,
		wg:             sync.WaitGroup{},
		ctx:            ctx,
		cancel:         cancel,
	}

	logger.Infof("TCP Proxy created: %s -> %s", tcpProxy.ListenAddr, tcpProxy.ServerAddr)

	return tcpProxy
}

func (proxy *TCPProxy) Start(hostWaitGroup *sync.WaitGroup) {
	hostWaitGroup.Add(1)
	proxy.wg.Add(1)
	defer func() {
		proxy.logger.Infof("TCP proxy stopped")
		hostWaitGroup.Done()
		proxy.wg.Done()
	}()

	listener, err := net.ListenTCP("tcp", proxy.ListenAddr)
	if err != nil {
		proxy.logger.Errorf("Failed to start TCP proxy: %v", err)
		return
	}
	defer listener.Close()

	proxy.logger.Debugf("TCP proxy started on %s", proxy.ListenAddr)

	stopChan := make(chan struct{})
	go func() {
		<-proxy.ctx.Done()
		proxy.logger.Infof("Stopping TCP proxy on %s", proxy.ListenAddr)
		close(stopChan)
		listener.Close() // Débloque listener.AcceptTCP(), qui passe dans le stopChan qui return la fonction, ce qui appelle tous les defer
	}()

	for {
		clientConn, err := listener.AcceptTCP()
		if err != nil {
			select {
			case <-stopChan:
				proxy.logger.Infof("Stopped accepting connections on %s", proxy.ListenAddr)
				return
			default:
				proxy.logger.Errorf("Failed to accept connection: %v", err)
				continue
			}
		}
		proxy.logger.Debugf("Accepted connection from %s", clientConn.RemoteAddr())

		proxy.wg.Add(1)
		go proxy.handleTCPConnection(clientConn)
	}
}

func (proxy *TCPProxy) Stop() {
	proxy.logger.Infof("Stopping TCP proxy")
	proxy.cancel()

	proxy.wg.Wait()
}

func (proxy *TCPProxy) shouldForwardProxy(clientConn *net.TCPConn) (bool, error) {
	if *proxy.hostState == hostState.Stopped || *proxy.hostState == hostState.Stopping {
		reader := bufio.NewReader(clientConn)
		peek, err := reader.Peek(utils.Min(1024, reader.Buffered()))

		if err != nil && err != io.EOF {
			return false, fmt.Errorf("failed to peek data: %v", err)
		}

		if utils.IsHTTPRequest(peek) {
			request := string(peek)

			proxy.logger.Debugf("Request: %s", request)

			if utils.CheckRequestHeader(request, "Status", "true") {
				proxy.logger.Debugf("Status request")
				clientConn.Close()
				return false, nil
			} else {
				proxy.logger.Debugf("Not a status request")
			}
		}

		err = proxy.StartHost()

		if err != nil {
			return false, fmt.Errorf("failed to start host: %v", err)
		}
	}

	hostStarted := hostState.WaitForState(proxy.hostState, hostState.Started, 20*time.Second)

	if !hostStarted {
		return false, fmt.Errorf("host took too long to start")
	}

	return true, nil
}

func (proxy *TCPProxy) handleTCPConnection(clientConn *net.TCPConn) {
	defer proxy.wg.Done()
	defer clientConn.Close()

	// Détermine si le proxy doit être redirigé vers la cible (host démarré + requete pas une requete de status)
	forwardProxy, err := proxy.shouldForwardProxy(clientConn)

	if err != nil {
		proxy.logger.Errorf("Failed to determine if proxy should be forwarded: %v", err)
		return
	}

	if !forwardProxy {
		proxy.logger.Verbosef("Proxy not forwarded to server, status request")
		return
	}

	proxy.logger.Infof("Proxy forwarded to server %s", proxy.ServerAddr)

	// Connexion à la cible
	serverConn, err := net.DialTCP("tcp", nil, proxy.ServerAddr)
	if err != nil {
		proxy.logger.Errorf("Failed to connect to server: %v", err)
		return
	}
	defer serverConn.Close()

	// Fonction qui va être appelée à chaque fois que des données sont transférées du client vers le serveur
	onClientToServer := func(bytesTransferred int) {
		proxy.logger.Verbosef("ClientToServer: %d bytes", bytesTransferred)

		if proxy.PacketReceived == nil {
			proxy.logger.Errorf("PacketReceived function not set")
			return
		}

		proxy.PacketReceived(proxy.Name)
	}

	// Fonction qui va être appelée à chaque fois que des données sont transférées du serveur vers le client
	onServerToClient := func(bytesTransferred int) {
		proxy.logger.Verbosef("ServerToClient: %d bytes", bytesTransferred)
	}

	clientToServerWriter := &utils.CustomWriter{Writer: serverConn, OnWrite: onClientToServer}
	serverToClientWriter := &utils.CustomWriter{Writer: clientConn, OnWrite: onServerToClient}

	// Channel pour savoir quand la copie client -> serveur est terminée
	doneCopyClientToServer := make(chan struct{})

	// Go routine pour copier les données du client vers le serveur
	go func() {
		defer close(doneCopyClientToServer)
		_, err := io.Copy(clientToServerWriter, clientConn)
		if err != nil {
			proxy.logger.Errorf("Error copying from client to server: %v", err)
		}
	}()

	// Channel pour savoir quand la copie serveur -> client est terminée
	doneCopyServerToClient := make(chan struct{})

	// Go routine pour copier les données du serveur vers le client
	go func() {
		defer close(doneCopyServerToClient)
		_, err = io.Copy(serverToClientWriter, serverConn)
		if err != nil {
			proxy.logger.Errorf("Error copying from server to client: %v", err)
		}
	}()

	select {

	case <-proxy.ctx.Done(): // Si le context global est annulé, on ferme les connexions client et serveur
		proxy.logger.Debugf("Context done")
		clientConn.Close()
		serverConn.Close()

	case <-doneCopyClientToServer: // Si la copie client -> serveur est terminée, on ferme la connexion serveur
		serverConn.Close()
		proxy.logger.Debugf("Client to server copy done")

	case <-doneCopyServerToClient: // Si la copie serveur -> client est terminée, on ferme la connexion client
		proxy.logger.Debugf("Server to client copy done")
		clientConn.Close()
	}

	proxy.logger.Debugf("Connection closed")
}
