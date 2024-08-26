package proxies

import (
	"context"
	"log"
	"mgarnier11/go-proxy/config"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedTCPProxy struct {
	mock.Mock
}

func (m *MockedTCPProxy) HostStarted() (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func (m *MockedTCPProxy) StartHost() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedTCPProxy) PacketReceived() error {
	args := m.Called()
	return args.Error(0)
}

var hostConfig = &config.HostConfig{
	Name: "test",
	Ip:   "localhost",
}

var proxyConfig = &config.ProxyConfig{
	ListenPort: 8080,
	TargetPort: 8081,
	Protocol:   "tcp",
	Name:       "test",
}

func sendAndReceiveBytes(tcpProxy *TCPProxy, data []byte, response []byte) (dataSent []byte, responseReceived []byte, err error) {
	listener, err := net.Listen("tcp", tcpProxy.TargetAddr)
	if err != nil {
		return nil, nil, err
	}
	defer listener.Close()

	sender, err := net.Dial("tcp", tcpProxy.ListenAddr)
	if err != nil {
		return nil, nil, err
	}
	defer sender.Close()

	sender.Write(data)

	// Accept the connection from the listener
	conn, err := listener.Accept()
	if err != nil {
		return nil, nil, err
	}
	defer conn.Close()

	bufTarget := make([]byte, 128)
	nTarget, err := conn.Read(bufTarget)
	if err != nil {
		return nil, nil, err
	}

	conn.Write(response)

	bufResponse := make([]byte, 128)
	nResponse, err := sender.Read(bufResponse)
	if err != nil {
		return nil, nil, err
	}

	return bufTarget[:nTarget], bufResponse[:nResponse], nil
}

func TestNewTCPProxy(t *testing.T) {
	tcpProxy := NewTCPProxy(nil, nil, &TCPProxyArgs{
		ProxyConfig: proxyConfig,
		HostConfig:  hostConfig,
	})

	assert.NotNil(t, tcpProxy)
	assert.Equal(t, "0.0.0.0:8080", tcpProxy.ListenAddr)
	assert.Equal(t, "localhost:8081", tcpProxy.TargetAddr)
}

func TestSendData(t *testing.T) {
	tcpProxyMock := new(MockedTCPProxy)
	tcpProxyMock.On("HostStarted").Return(true, nil)
	tcpProxyMock.On("StartHost").Return(nil)
	tcpProxyMock.On("PacketReceived").Return(nil)

	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	tcpProxy := NewTCPProxy(ctx, cancel, &TCPProxyArgs{
		ProxyConfig:    proxyConfig,
		HostConfig:     hostConfig,
		HostStarted:    tcpProxyMock.HostStarted,
		StartHost:      tcpProxyMock.StartHost,
		PacketReceived: tcpProxyMock.PacketReceived,
	})

	go tcpProxy.Start(&wg)

	dataSent, dataReceived, err := sendAndReceiveBytes(tcpProxy, []byte("This is the data"), []byte("Here is the response"))

	tcpProxy.Stop()
	wg.Wait()

	assert.Nil(t, err)
	assert.Equal(t, "This is the data", string(dataSent))
	assert.Equal(t, "Here is the response", string(dataReceived))
	tcpProxyMock.AssertNumberOfCalls(t, "HostStarted", 1)
	tcpProxyMock.AssertNumberOfCalls(t, "PacketReceived", 1)
}

func TestTCPProxyStop(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	tcpProxy := NewTCPProxy(ctx, cancel, &TCPProxyArgs{
		ProxyConfig: proxyConfig,
		HostConfig:  hostConfig,
	})

	go tcpProxy.Start(&wg)

	// time.Sleep(500 * time.Millisecond)

	tcpProxy.Stop()
	wg.Wait()

	_, _, err := sendAndReceiveBytes(tcpProxy, []byte("This is the data"), []byte("Here is the response"))

	assert.NotNil(t, err)
}

func TestTCPProxyForceStop(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	tcpProxy := NewTCPProxy(ctx, cancel, &TCPProxyArgs{
		ProxyConfig: proxyConfig,
		HostConfig:  hostConfig,
		HostStarted: func() (bool, error) {
			return true, nil
		},
		PacketReceived: func() error {
			cancel()
			return nil
		},
	})

	go tcpProxy.Start(&wg)

	dataSent, dataReceived, err := sendAndReceiveBytes(tcpProxy, []byte("This is the data"), []byte("Here is the response"))

	if err != nil {
		log.Println(err)
	}
	assert.Nil(t, err)
	assert.Equal(t, "This is the data", string(dataSent))
	assert.Equal(t, "Here is the response", string(dataReceived))

	time.Sleep(500 * time.Millisecond)

	tcpProxy.Stop()
	wg.Wait()

	// _, _, err := sendAndReceiveBytes(tcpProxy, []byte("This is the data"), []byte("Here is the response"))

	// assert.NotNil(t, err)
}

func TestStartingHost(t *testing.T) {
	tcpProxyMock := new(MockedTCPProxy)
	tcpProxyMock.On("HostStarted").Return(true, nil)
	tcpProxyMock.On("StartHost").Return(nil)
	tcpProxyMock.On("PacketReceived").Return(nil)

	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	tcpProxy := NewTCPProxy(ctx, cancel, &TCPProxyArgs{
		ProxyConfig:    proxyConfig,
		HostConfig:     hostConfig,
		HostStarted:    tcpProxyMock.HostStarted,
		StartHost:      tcpProxyMock.StartHost,
		PacketReceived: tcpProxyMock.PacketReceived,
	})

	go tcpProxy.Start(&wg)

	dataSent, dataReceived, err := sendAndReceiveBytes(tcpProxy, []byte("This is the data"), []byte("Here is the response"))

	tcpProxy.Stop()
	wg.Wait()

	assert.Nil(t, err)
	assert.Equal(t, "This is the data", string(dataSent))
	assert.Equal(t, "Here is the response", string(dataReceived))
	tcpProxyMock.AssertNumberOfCalls(t, "HostStarted", 1)
	tcpProxyMock.AssertNumberOfCalls(t, "StartHost", 0)
}
