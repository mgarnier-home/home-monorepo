package proxies

import (
	"context"
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

func sendAndCheckBytes(tcpProxy *TCPProxy, data []byte) (dataReceived []byte, err error) {
	listener, err := net.Listen("tcp", tcpProxy.TargetAddr)
	if err != nil {
		return nil, err
	}
	defer listener.Close()

	sender, err := net.Dial("tcp", tcpProxy.ListenAddr)
	if err != nil {
		return nil, err
	}
	defer sender.Close()

	sender.Write(data)

	// Accept the connection from the listener
	conn, err := listener.Accept()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}

	return buf[:n], nil
}

func TestNewTCPProxy(t *testing.T) {

	tcpProxy := NewTCPProxy(hostConfig, proxyConfig, nil, nil, nil)

	assert.NotNil(t, tcpProxy)
	assert.Equal(t, "0.0.0.0:8080", tcpProxy.ListenAddr)
	assert.Equal(t, "localhost:8081", tcpProxy.TargetAddr)
}

func TestSendData(t *testing.T) {
	tcpProxyMock := new(MockedTCPProxy)
	tcpProxyMock.On("HostStarted").Return(true, nil)
	tcpProxyMock.On("StartHost").Return(nil)
	tcpProxyMock.On("PacketReceived").Return(nil)
	tcpProxy := NewTCPProxy(hostConfig, proxyConfig, tcpProxyMock.HostStarted, tcpProxyMock.StartHost, tcpProxyMock.PacketReceived)

	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	go tcpProxy.Start(&wg, ctx, cancel)

	dataReceived, err := sendAndCheckBytes(tcpProxy, []byte("Hello World"))

	time.Sleep(500 * time.Millisecond)

	tcpProxy.Stop()
	wg.Wait()

	assert.Nil(t, err)
	assert.Equal(t, "Hello World", string(dataReceived))

}

func TestTCPProxyStop(t *testing.T) {
	tcpProxy := NewTCPProxy(hostConfig, proxyConfig, nil, nil, nil)

	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	go tcpProxy.Start(&wg, ctx, cancel)

	time.Sleep(500 * time.Millisecond)

	tcpProxy.Stop()
	wg.Wait()

	_, err := sendAndCheckBytes(tcpProxy, []byte("Hello World"))

	assert.NotNil(t, err)

}

func TestStartingHost(t *testing.T) {
	tcpProxyMock := new(MockedTCPProxy)
	tcpProxyMock.On("HostStarted").Return(true, nil)
	tcpProxyMock.On("StartHost").Return(nil)
	tcpProxy := NewTCPProxy(hostConfig, proxyConfig, tcpProxyMock.HostStarted, tcpProxyMock.StartHost, tcpProxyMock.PacketReceived)

	var wg sync.WaitGroup
	wg.Add(1)
	ctx, cancel := context.WithCancel(context.Background())

	go tcpProxy.Start(&wg, ctx, cancel)

	dataReceived, err := sendAndCheckBytes(tcpProxy, []byte("Hello World"))

	time.Sleep(500 * time.Millisecond)

	tcpProxy.Stop()
	wg.Wait()

	assert.Nil(t, err)
	assert.Equal(t, "Hello World", string(dataReceived))
	// tcpProxyMock.AssertExpectations(t)
	tcpProxyMock.AssertNumberOfCalls(t, "HostStarted", 1)
	tcpProxyMock.AssertNumberOfCalls(t, "StartHost", 0)
	// tcpProxyMock.AssertNumberOfCalls(t, "PacketReceived", 0)
}
