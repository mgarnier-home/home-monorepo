package dto

import "mgarnier11/mineager/config"

type HostDto struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
}

func DockerHostToHostDto(host *config.DockerHostConfig) *HostDto {
	return &HostDto{
		Name: host.Name,
		Ip:   host.Ip,
	}
}

func DockerHostsToHostsDto(hosts []*config.DockerHostConfig) []*HostDto {
	hostsDto := make([]*HostDto, 0)

	for _, host := range hosts {
		hostsDto = append(hostsDto, DockerHostToHostDto(host))
	}

	return hostsDto
}
