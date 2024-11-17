package proxies

import (
	"bufio"
	"context"
	"fmt"
	"goUtils"
	"io"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/utils"
	"net"
	"sync"

	"github.com/charmbracelet/log"
)

type TCPProxy struct {
	hostName       string
	Name           string
	ListenAddr     *net.TCPAddr
	ServerAddr     *net.TCPAddr
	HostStarted    func() (bool, error)
	StartHost      func(proxy *TCPProxy) error
	PacketReceived func(proxy *TCPProxy) error

	wg     sync.WaitGroup
	ctx    context.Context
	cancel context.CancelFunc
}

type TCPProxyArgs struct {
	HostName       string
	HostIp         string
	ProxyConfig    *config.ProxyConfig
	HostStarted    func() (bool, error)
	StartHost      func(proxy *TCPProxy) error
	PacketReceived func(proxy *TCPProxy) error
}

func NewTCPProxy(args *TCPProxyArgs) *TCPProxy {
	ctx, cancel := context.WithCancel(context.Background())

	listenAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", "0.0.0.0", args.ProxyConfig.ListenPort))
	if err != nil {
		log.Errorf("%-10s %-20s Failed to resolve listen TCP address %d: %v", args.HostName, args.ProxyConfig.Name, args.ProxyConfig.ListenPort, err)
		panic(err)
	}

	serverAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%d", args.HostIp, args.ProxyConfig.ServerPort))
	if err != nil {
		log.Errorf("%-10s %-20s Failed to resolve server TCP address %d: %v", args.HostName, args.ProxyConfig.Name, args.ProxyConfig.ServerPort, err)
		panic(err)
	}

	tcpProxy := &TCPProxy{
		hostName:       args.HostName,
		Name:           args.ProxyConfig.Name,
		ListenAddr:     listenAddr,
		ServerAddr:     serverAddr,
		HostStarted:    args.HostStarted,
		StartHost:      args.StartHost,
		PacketReceived: args.PacketReceived,
		wg:             sync.WaitGroup{},
		ctx:            ctx,
		cancel:         cancel,
	}

	log.Infof("%-10s %-20s TCP Proxy created: %s -> %s", tcpProxy.hostName, tcpProxy.Name, tcpProxy.ListenAddr, tcpProxy.ServerAddr)

	return tcpProxy
}

func (proxy *TCPProxy) Start(hostWaitGroup *sync.WaitGroup) {
	hostWaitGroup.Add(1)
	proxy.wg.Add(1)
	defer func() {
		log.Infof("%-10s %-20s TCP proxy stopped", proxy.hostName, proxy.Name)
		hostWaitGroup.Done()
		proxy.wg.Done()
	}()

	listener, err := net.ListenTCP("tcp", proxy.ListenAddr)
	if err != nil {
		log.Errorf("%-10s %-20s %s: Failed to start TCP proxy: %v", proxy.hostName, proxy.hostName, proxy.Name, err)
		return
	}
	defer listener.Close()

	log.Debugf("%-10s %-20s TCP proxy started on %s", proxy.hostName, proxy.Name, proxy.ListenAddr)

	stopChan := make(chan struct{})
	go func() {
		<-proxy.ctx.Done()
		log.Infof("%-10s %-20s Stopping TCP proxy on %s", proxy.hostName, proxy.Name, proxy.ListenAddr)
		close(stopChan)
		listener.Close() // Débloque listener.AcceptTCP(), qui passe dans le stopChan qui return la fonction, ce qui appelle tous les defer
	}()

	for {
		clientConn, err := listener.AcceptTCP()
		if err != nil {
			select {
			case <-stopChan:
				log.Infof("%-10s %-20s Stopped accepting connections on %s", proxy.hostName, proxy.Name, proxy.ListenAddr)
				return
			default:
				log.Errorf("%-10s %-20s Failed to accept connection: %v", proxy.hostName, proxy.Name, err)
				continue
			}
		}
		log.Debugf("%-10s %-20s Accepted connection from %s", proxy.hostName, proxy.Name, clientConn.RemoteAddr())

		proxy.wg.Add(1)
		go proxy.handleTCPConnection(clientConn)
	}
}

func (proxy *TCPProxy) Stop() {
	log.Infof("%-10s %-20s Stopping TCP proxy", proxy.hostName, proxy.Name)
	proxy.cancel()

	proxy.wg.Wait()
}

func (proxy *TCPProxy) shouldForwardProxy(clientConn *net.TCPConn) (bool, error) {
	if proxy.HostStarted == nil {
		return false, fmt.Errorf("%s: HostStarted function not set", proxy.Name)
	}

	started, err := proxy.HostStarted()
	if err != nil {
		return false, err
	}

	if !started {
		reader := bufio.NewReader(clientConn)
		peek, err := reader.Peek(goUtils.Min(1024, reader.Buffered()))

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

func (proxy *TCPProxy) handleTCPConnection(clientConn *net.TCPConn) {
	defer proxy.wg.Done()
	defer clientConn.Close()

	// Détermine si le proxy doit être redirigé vers la cible (host démarré + requete pas une requete de status)
	forwardProxy, err := proxy.shouldForwardProxy(clientConn)

	if err != nil {
		log.Errorf("%-10s %-20s Failed to determine if proxy should be forwarded: %v", proxy.hostName, proxy.Name, err)
		return
	}

	if !forwardProxy {
		log.Infof("%-10s %-20s Proxy not forwarded to server", proxy.hostName, proxy.Name)
		return
	}

	log.Infof("%-10s %-20s Proxy forwarded to server %s", proxy.hostName, proxy.Name, proxy.ServerAddr)

	// Connexion à la cible
	serverConn, err := net.DialTCP("tcp", nil, proxy.ServerAddr)
	if err != nil {
		log.Errorf("%-10s %-20s Failed to connect to server: %v", proxy.hostName, proxy.Name, err)
		return
	}
	defer serverConn.Close()

	// Fonction qui va être appelée à chaque fois que des données sont transférées du client vers le serveur
	onClientToServer := func(bytesTransferred int) {
		log.Debugf("%-10s %-20s ClientToServer: %d bytes", proxy.hostName, proxy.Name, bytesTransferred)

		if proxy.PacketReceived == nil {
			log.Debugf("%-10s %-20s PacketReceived function not set", proxy.hostName, proxy.Name)
			return
		}

		proxy.PacketReceived(proxy)
	}

	// Fonction qui va être appelée à chaque fois que des données sont transférées du serveur vers le client
	onServerToClient := func(bytesTransferred int) {
		log.Debugf("%-10s %-20s ServerToClient: %d bytes", proxy.hostName, proxy.Name, bytesTransferred)
	}

	clientToServerWriter := &goUtils.CustomWriter{Writer: serverConn, OnWrite: onClientToServer}
	serverToClientWriter := &goUtils.CustomWriter{Writer: clientConn, OnWrite: onServerToClient}

	// Channel pour savoir quand la copie client -> serveur est terminée
	doneCopyClientToServer := make(chan struct{})

	// Go routine pour copier les données du client vers le serveur
	go func() {
		defer close(doneCopyClientToServer)
		_, err := io.Copy(clientToServerWriter, clientConn)
		if err != nil {
			log.Debugf("%-10s %-20s Error copying from client to server: %v", proxy.hostName, proxy.Name, err)
		}
	}()

	// Channel pour savoir quand la copie serveur -> client est terminée
	doneCopyServerToClient := make(chan struct{})

	// Go routine pour copier les données du serveur vers le client
	go func() {
		defer close(doneCopyServerToClient)
		_, err = io.Copy(serverToClientWriter, serverConn)
		if err != nil {
			log.Debugf("%-10s %-20s Error copying from server to client: %v", proxy.hostName, proxy.Name, err)
		}
	}()

	select {

	case <-proxy.ctx.Done(): // Si le context global est annulé, on ferme les connexions client et serveur
		log.Infof("%-10s %-20s Context done", proxy.hostName, proxy.Name)
		clientConn.Close()
		serverConn.Close()

	case <-doneCopyClientToServer: // Si la copie client -> serveur est terminée, on ferme la connexion serveur
		serverConn.Close()
		log.Infof("%-10s %-20s Client to server copy done", proxy.hostName, proxy.Name)

	case <-doneCopyServerToClient: // Si la copie serveur -> client est terminée, on ferme la connexion client
		log.Infof("%-10s %-20s Server to client copy done", proxy.hostName, proxy.Name)
		clientConn.Close()
	}

	log.Infof("%-10s %-20s Connection closed", proxy.hostName, proxy.Name)
}
