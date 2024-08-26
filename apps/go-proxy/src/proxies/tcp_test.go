package proxies

import (
	"bytes"
	"context"
	"io"
	"log"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/utils"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockedTCPProxy struct {
	mock.Mock
}

func setupTest(hostStartedReturn bool) (*sync.WaitGroup, context.Context, *bytes.Buffer, *MockedTCPProxy) {
	var logs bytes.Buffer
	multiWriter := io.MultiWriter(os.Stdout, &logs)
	log.SetOutput(multiWriter)

	var wg sync.WaitGroup
	wg.Add(1)
	ctx := context.Background()

	tcpProxyMock := new(MockedTCPProxy)
	tcpProxyMock.On("HostStarted").Return(hostStartedReturn, nil)
	tcpProxyMock.On("StartHost").Return(nil)
	tcpProxyMock.On("PacketReceived").Return(nil)

	return &wg, ctx, &logs, tcpProxyMock
}

func (m *MockedTCPProxy) HostStarted(proxy *TCPProxy) (bool, error) {
	args := m.Called()
	return args.Bool(0), args.Error(1)
}

func (m *MockedTCPProxy) StartHost(proxy *TCPProxy) error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockedTCPProxy) PacketReceived(proxy *TCPProxy) error {
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
	tcpProxy := NewTCPProxy(context.TODO(), &TCPProxyArgs{
		ProxyConfig: proxyConfig,
		HostConfig:  hostConfig,
	})

	assert.NotNil(t, tcpProxy)
	assert.Equal(t, "0.0.0.0:8080", tcpProxy.ListenAddr)
	assert.Equal(t, "localhost:8081", tcpProxy.TargetAddr)
}

func TestSendData(t *testing.T) {
	wg, ctx, logs, tcpProxyMock := setupTest(true)

	tcpProxy := NewTCPProxy(ctx, &TCPProxyArgs{
		ProxyConfig:    proxyConfig,
		HostConfig:     hostConfig,
		HostStarted:    tcpProxyMock.HostStarted,
		StartHost:      tcpProxyMock.StartHost,
		PacketReceived: tcpProxyMock.PacketReceived,
	})

	go tcpProxy.Start(wg)

	dataSent, dataReceived, err := sendAndReceiveBytes(tcpProxy, []byte("This is the data"), []byte("Here is the response"))

	tcpProxy.Stop()
	wg.Wait()

	assert.Nil(t, err)
	assert.Equal(t, "This is the data", string(dataSent))
	assert.Equal(t, "Here is the response", string(dataReceived))
	assert.NotRegexp(t, `(?i)error`, logs.String())
	tcpProxyMock.AssertNumberOfCalls(t, "HostStarted", 1)
	tcpProxyMock.AssertNumberOfCalls(t, "PacketReceived", 1)
}

func TestTCPProxyStop(t *testing.T) {
	longData, _ := utils.GenerateRandomData(1024 * 1024)

	wg, ctx, logs, _ := setupTest(true)

	tcpProxy := NewTCPProxy(ctx, &TCPProxyArgs{
		ProxyConfig: proxyConfig,
		HostConfig:  hostConfig,
	})

	go tcpProxy.Start(wg)

	tcpProxy.Stop()
	wg.Wait()

	_, _, err := sendAndReceiveBytes(tcpProxy, longData, []byte("Here is the response"))

	assert.NotNil(t, err)
	assert.NotRegexp(t, `(?i)error`, logs.String())
}

func TestTCPProxyForceStop(t *testing.T) {
	wg, ctx, logs, _ := setupTest(true)

	tcpProxy := NewTCPProxy(ctx, &TCPProxyArgs{
		ProxyConfig: proxyConfig,
		HostConfig:  hostConfig,
		HostStarted: func(proxy *TCPProxy) (bool, error) {
			return true, nil
		},
		PacketReceived: func(proxy *TCPProxy) error {
			proxy.Stop()
			return nil
		},
	})

	go tcpProxy.Start(wg)

	dataSent, dataReceived, err := sendAndReceiveBytes(tcpProxy, []byte("This is the data"), []byte("Here is the response"))

	time.Sleep(500 * time.Millisecond)

	tcpProxy.Stop()
	wg.Wait()

	assert.Nil(t, err)
	assert.Equal(t, "This is the data", string(dataSent))
	assert.Equal(t, "Here is the response", string(dataReceived))
	assert.Regexp(t, `(?i)Error copying from client to target: writeto tcp 127.0.0.1:8080`, logs.String())
}

func TestStartingHost(t *testing.T) {
	wg, ctx, logs, tcpProxyMock := setupTest(false)

	tcpProxy := NewTCPProxy(ctx, &TCPProxyArgs{
		ProxyConfig:    proxyConfig,
		HostConfig:     hostConfig,
		HostStarted:    tcpProxyMock.HostStarted,
		StartHost:      tcpProxyMock.StartHost,
		PacketReceived: tcpProxyMock.PacketReceived,
	})

	go tcpProxy.Start(wg)

	dataSent, dataReceived, err := sendAndReceiveBytes(tcpProxy, []byte("This is the data"), []byte("Here is the response"))

	time.Sleep(500 * time.Millisecond)

	tcpProxy.Stop()
	wg.Wait()

	assert.Nil(t, err)
	assert.Equal(t, "This is the data", string(dataSent))
	assert.Equal(t, "Here is the response", string(dataReceived))
	assert.NotRegexp(t, `(?i)error`, logs.String())
	tcpProxyMock.AssertNumberOfCalls(t, "HostStarted", 1)
	tcpProxyMock.AssertNumberOfCalls(t, "StartHost", 1)
}
