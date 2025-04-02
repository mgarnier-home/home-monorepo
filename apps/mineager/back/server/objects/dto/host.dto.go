package dto

import (
	"mgarnier11.fr/go/mineager/server/objects/bo"
)

type HostDto struct {
	Name string `json:"name"`
	Ip   string `json:"ip"`
	Ping bool   `json:"ping"`
}

func MapHostBoToHostDto(host *bo.HostBo) *HostDto {
	return &HostDto{
		Name: host.Name,
		Ip:   host.Ip,
		Ping: host.Ping,
	}
}

func MapHostsBoHostsDto(hosts []*bo.HostBo) []*HostDto {
	hostsDto := make([]*HostDto, 0)

	for _, host := range hosts {
		hostsDto = append(hostsDto, MapHostBoToHostDto(host))
	}

	return hostsDto
}
