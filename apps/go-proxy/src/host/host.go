package host

import (
	"context"
	"fmt"
	"slices"
	"strings"
	"sync"
	"time"

	"mgarnier11.fr/go/libs/colors"
	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/ntfy"
	"mgarnier11.fr/go/libs/utils"

	"mgarnier11.fr/go/go-proxy/config"
	"mgarnier11.fr/go/go-proxy/docker"
	"mgarnier11.fr/go/go-proxy/hostState"
	"mgarnier11.fr/go/go-proxy/proxies"

	"github.com/charmbracelet/lipgloss"
)

type Host struct {
	Proxies             map[string]*proxies.TCPProxy
	State               hostState.State
	LastPacketDate      time.Time
	LastPacketProxyName string
	Config              *config.HostConfig

	logger *logger.Logger

	waitGroup sync.WaitGroup
	ctx       context.Context
	cancel    context.CancelFunc
}

func NewHost(hostConfig *config.HostConfig) *Host {

	ctx, cancel := context.WithCancel(context.Background())

	host := &Host{
		Config:    hostConfig,
		Proxies:   make(map[string]*proxies.TCPProxy),
		State:     hostState.Stopped,
		waitGroup: sync.WaitGroup{},
		ctx:       ctx,
		cancel:    cancel,
		logger: logger.
			NewLogger(
				fmt.Sprintf("[%s]",
					strings.ToUpper(hostConfig.Name)),
				"%-10s ",
				lipgloss.NewStyle().Foreground(lipgloss.Color(colors.GenerateHexColor(hostConfig.Name))),
				nil,
			),
	}

	host.logger.Infof("created")

	go host.setupHostLoop()

	host.StartHost("Starting proxy")
	host.LastPacketDate = time.Now()

	return host
}

func (host *Host) setupHostLoop() {
	host.waitGroup.Add(1)
	host.logger.Infof("starting host loop")

	defer func() {
		host.logger.Infof("stopping host loop")
		host.waitGroup.Done()
	}()

	stateTicker := time.NewTicker(1 * time.Second)
	defer stateTicker.Stop()

	dockerTicker := time.NewTicker(10 * time.Second)
	defer dockerTicker.Stop()

	inactivityTicker := time.NewTicker(15 * time.Second)
	defer inactivityTicker.Stop()

	for {
		select {
		case <-stateTicker.C:
			host.updateState()
		case <-dockerTicker.C:

			if host.State == hostState.Started {
				proxies := host.Config.Proxies

				host.logger.Infof("Getting proxies from docker")

				dockerProxies, err := docker.GetProxiesFromDocker(host.Config.SSHUsername, host.Config.Ip, host.Config.SSHPort, host.logger)

				if err != nil {
					host.logger.Errorf("failed to get proxies from docker: %v", err)
				}

				proxies = append(proxies, dockerProxies...)

				host.setupProxies(proxies)
			}
		case <-inactivityTicker.C:
			timeout := time.Duration(host.Config.MaxAliveTime) * time.Minute
			if host.Config.Autostop {
				if host.State == hostState.Started && time.Since(host.LastPacketDate) > timeout {
					host.logger.Infof("Host has been inactive for too long, stopping it")
					go host.StopHost()
				} else if host.State == hostState.Started {
					host.logger.Infof("Time remaining before inactivity timeout: %v", timeout-time.Since(host.LastPacketDate).Round(time.Second))
				} else if host.State == hostState.Stopped {
					host.logger.Infof("Server stopped since %v", time.Since(host.LastPacketDate.Add(timeout)).Round(time.Second))
				}
			} else {
				host.logger.Infof("Autostop is disabled")
			}

		case <-host.ctx.Done():
			// Ensure we break out of the loop if the context is cancelled

			return
		}
	}
}

func (host *Host) updateState() {
	pingSuccess, err := utils.PingIp(host.Config.Ip, 500*time.Millisecond)

	if err != nil {
		host.logger.Errorf("failed to check host status: %v", err)
		return
	}

	if host.State == hostState.Started && !pingSuccess {
		host.State = hostState.Stopped
	} else if host.State == hostState.Stopped && pingSuccess {
		host.State = hostState.Started
	} else if host.State == hostState.Starting && pingSuccess {
		host.State = hostState.Started
	} else if host.State == hostState.Stopping && !pingSuccess {
		host.State = hostState.Stopped
	}

}

func (host *Host) setupProxies(proxyConfigs []*config.ProxyConfig) {
	for key := range host.Proxies {
		exists := slices.ContainsFunc(proxyConfigs, func(proxy *config.ProxyConfig) bool {
			return proxy.Key == key
		})

		if !exists {
			host.DisposeProxy(key)
		}
	}

	for _, proxyConfig := range proxyConfigs {
		if host.Proxies[proxyConfig.Key] != nil {
			host.logger.Debugf("%s already exists", proxyConfig.Key)
			continue
		}

		host.Proxies[proxyConfig.Key] = proxies.NewTCPProxy(&proxies.TCPProxyArgs{
			HostIp:         host.Config.Ip,
			ProxyConfig:    proxyConfig,
			HostState:      &host.State,
			StartHost:      host.StartHost,
			PacketReceived: host.PacketReceived,
		}, host.logger)

		go host.Proxies[proxyConfig.Key].Start(&host.waitGroup)
	}
}

func (host *Host) StartHost(proxyName string) error {
	if host.State != hostState.Stopped {
		host.logger.Infof("Cannot start host, state is not stopped : %s", host.State.String())
		return nil
	}

	host.State = hostState.Starting

	if packet, err := newMagicPacket(host.Config.MacAddress); err == nil {
		packet.send("255.255.255.255")
		host.logger.Debugf("Sent magic packet to start host")
	} else {
		return fmt.Errorf("failed to send magic packet: %v", err)
	}

	err := ntfy.SendNotification("Proxy", fmt.Sprintf("Starting host %s\nRequest coming from %s", host.Config.Name, proxyName), "")

	if err != nil {
		host.logger.Warnf("failed to send notification: %v", err)
	}

	hostStarted := hostState.WaitForState(&host.State, hostState.Started, 20*time.Second)

	if !hostStarted {
		host.State = hostState.Stopped
		return fmt.Errorf("Host took too long to start")
	} else {
		return nil
	}
}

func (host *Host) StopHost() {
	if host.State != hostState.Started {
		host.logger.Infof("Cannot stop host, state is not started : %s", host.State.String())
		return
	}

	host.State = hostState.Stopping

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		err := sendSSHCommand(ctx, host.Config, "sudo pm-suspend &")

		if err != nil {
			host.logger.Errorf("failed to stop host: %v", err)
		}
	}()

	err := ntfy.SendNotification("Proxy", fmt.Sprintf("Stopping host %s", host.Config.Name), "")

	if err != nil {
		host.logger.Warnf("failed to send notification: %v", err)
	}

	hostStopped := hostState.WaitForState(&host.State, hostState.Stopped, 20*time.Second)

	if !hostStopped {
		host.State = hostState.Started
		host.logger.Errorf("Host took too long to stop")
	}
}

func (host *Host) PacketReceived(proxyName string) {
	host.LastPacketDate = time.Now()
	host.LastPacketProxyName = proxyName
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
