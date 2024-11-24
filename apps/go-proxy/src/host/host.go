package host

import (
	"context"
	"errors"
	"fmt"
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/docker"
	"mgarnier11/go-proxy/proxies"
	"mgarnier11/go-proxy/utils"
	"slices"
	"strings"
	"sync"
	"time"
)

type Host struct {
	Proxies        map[string]*proxies.TCPProxy
	Started        bool
	LastPacketDate time.Time
	Config         *config.HostConfig

	logger *utils.Logger

	stopping bool
	starting bool

	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewHost(hostConfig *config.HostConfig) *Host {
	ctx, cancel := context.WithCancel(context.Background())

	host := &Host{
		Config:    hostConfig,
		Proxies:   make(map[string]*proxies.TCPProxy),
		waitGroup: sync.WaitGroup{},
		ctx:       ctx,
		cancel:    cancel,
		logger:    utils.NewLogger(fmt.Sprintf("[%s]", strings.ToUpper(hostConfig.Name)), "%-10s ", nil),
	}

	host.logger.Infof("created")

	host.StartHost(nil)
	host.Started = true
	host.LastPacketDate = time.Now()

	dockerProxies, _ := docker.GetProxiesFromDocker(hostConfig.SSHUsername, hostConfig.Ip, host.logger)

	host.setupProxies(slices.Concat(dockerProxies, hostConfig.Proxies))

	go host.setupContainersListener()

	return host
}

func (host *Host) setupContainersListener() {
	host.waitGroup.Add(1)
	host.logger.Infof("listening for docker containers")

	defer func() {
		host.logger.Infof("stopped listening for docker containers")
		host.waitGroup.Done()
	}()

	ticker := time.NewTicker(50 * time.Second)

	for {
		select {
		case <-ticker.C:
			host.logger.Infof("started: %v", host.Started)
			if host.Started {
				host.logger.Infof("Getting proxies from docker")

				proxies, err := docker.GetProxiesFromDocker(host.Config.SSHUsername, host.Config.Ip, host.logger)

				if err != nil {
					host.logger.Errorf("failed to get proxies from docker: %v", err)
				} else {
					host.setupProxies(proxies)
				}
			}
		case <-host.ctx.Done():
			// Ensure we break out of the loop if the context is cancelled
			ticker.Stop()
			return
		}
	}

}

func (host *Host) setupProxies(proxyConfigs []*config.ProxyConfig) {
	for name := range host.Proxies {
		exists := slices.ContainsFunc(proxyConfigs, func(proxy *config.ProxyConfig) bool {
			return proxy.Name == name
		})

		if !exists {
			host.DisposeProxy(name)
		}
	}

	for _, proxyConfig := range proxyConfigs {
		if host.Proxies[proxyConfig.Name] != nil {
			host.logger.Debugf("%s already exists", proxyConfig.Name)
			continue
		}

		host.Proxies[proxyConfig.Name] = proxies.NewTCPProxy(&proxies.TCPProxyArgs{
			HostIp:         host.Config.Ip,
			ProxyConfig:    proxyConfig,
			HostStarted:    host.HostStarted,
			StartHost:      host.StartHost,
			PacketReceived: host.PacketReceived,
		}, host.logger)

		go host.Proxies[proxyConfig.Name].Start(&host.waitGroup)
	}
}

func (host *Host) HostStarted() (bool, error) {
	host.logger.Infof("Checking host status")

	pingSuccess, err := getHostStatus(host.Config.Ip)

	if err != nil {
		host.logger.Errorf("failed to check host status: %v", err)

		return false, err
	}

	return pingSuccess, nil
}

func (host *Host) StartHost(proxy *proxies.TCPProxy) error {
	if host.starting {
		host.logger.Infof("Already starting")

		return nil
	}

	host.starting = true
	defer func() {
		host.starting = false
	}()

	host.logger.Infof("Starting")

	if packet, err := newMagicPacket(host.Config.MacAddress); err == nil {
		packet.send("255.255.255.255")
		host.logger.Debugf("sent magic packet")
	} else {
		host.logger.Errorf("failed to send magic packet: %v", err)

		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	for {
		select {
		case <-ctx.Done(): // Context timeout or cancellation
			host.logger.Debugf("Context canceled or timed out")
			return errors.New("host took too long to start")
		default:
			hostPinged, err := host.HostStarted()

			host.logger.Debugf("Host pinged: %v", hostPinged)

			if err != nil {
				host.logger.Errorf("failed to check host status: %v", err)
			} else if hostPinged {
				host.logger.Infof("Started")

				host.Started = true

				return nil
			}

			time.Sleep(1 * time.Second)
		}
	}

}

func (host *Host) StopHost() {
	if host.stopping {
		host.logger.Infof("Already stopping")

		return
	}

	host.stopping = true
	defer func() {
		host.stopping = false
	}()

	host.logger.Infof("Stopping")

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	go func() {

		err := sendSSHCommand(ctx, host.Config, "sudo pm-suspend &")

		if err != nil {
			host.logger.Errorf("failed to stop host: %v", err)
		}
	}()

	for {
		select {
		case <-ctx.Done(): // Context timeout or cancellation
			host.logger.Warnf("Context canceled or timed out")
			return
		default:
			hostPinged, err := host.HostStarted()

			host.logger.Infof("Host pinged: %v", hostPinged)

			if err != nil {
				host.logger.Errorf("failed to check host status: %v", err)
			} else if !hostPinged {
				host.logger.Infof("Stopped")

				host.Started = false

				return
			}

			time.Sleep(1 * time.Second)
		}
	}

}

func (host *Host) PacketReceived(proxy *proxies.TCPProxy) error {
	host.LastPacketDate = time.Now()
	return nil
}

func (host *Host) DisposeProxy(proxyName string) {
	proxy := host.Proxies[proxyName]

	if proxy == nil {
		host.logger.Errorf("%s: proxy does not exist", proxyName)
		return
	}

	proxy.Stop()
	delete(host.Proxies, proxyName)

	host.logger.Infof("%s: disposed", proxyName)
}

func (host *Host) Dispose() {
	host.logger.Infof("disposing")

	host.cancel()

	for name := range host.Proxies {
		host.DisposeProxy(name)
	}

	host.waitGroup.Wait()

	host.logger.Infof("disposed")
}
