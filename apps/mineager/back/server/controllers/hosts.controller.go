package controllers

import (
	"time"

	"mgarnier11.fr/go/mineager/config"
	"mgarnier11.fr/go/mineager/server/objects/bo"

	"mgarnier11.fr/go/libs/logger"
	"mgarnier11.fr/go/libs/utils"
)

type HostsController struct {
}

func NewHostsController() *HostsController {
	return &HostsController{}
}

func mapDockerHostToHostBo(host *config.DockerHostConfig) *bo.HostBo {
	ping, err := utils.PingIp(host.Ip, 500*time.Millisecond)

	if err != nil {
		logger.Errorf("failed to check host status: %v", err)
	}

	return &bo.HostBo{
		Name:         host.Name,
		Ip:           host.Ip,
		ProxyIp:      host.ProxyIp,
		SSHUsername:  host.SSHUsername,
		SSHPort:      host.SSHPort,
		StartPort:    host.StartPort,
		MineagerPath: host.MineagerPath,
		Ping:         ping,
	}
}

func (controller *HostsController) GetHosts() []*bo.HostBo {
	hosts := config.Config.AppConfig.DockerHosts

	boHosts := make([]*bo.HostBo, 0)

	for _, host := range hosts {
		boHosts = append(boHosts, mapDockerHostToHostBo(host))
	}

	return boHosts
}

func (controller *HostsController) GetHost(hostName string) (*bo.HostBo, error) {
	host, err := config.GetHost(hostName)
	if err != nil {
		return nil, err
	}

	return mapDockerHostToHostBo(host), nil
}
