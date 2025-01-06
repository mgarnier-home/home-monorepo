package dto

import "mgarnier11/mineager/server/objects/bo"

type ServerDto struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	Version string `json:"version"`
	Map     string `json:"map"`
}

func ServerBoToServerDto(serverBo *bo.ServerBo) *ServerDto {
	return &ServerDto{
		Name:    serverBo.Name,
		Url:     serverBo.Url,
		Version: serverBo.Version,
		Map:     serverBo.Map,
	}
}

func ServersBoToServersDto(serversBo []*bo.ServerBo) []*ServerDto {
	serversDto := make([]*ServerDto, 0)

	for _, serverBo := range serversBo {
		serversDto = append(serversDto, ServerBoToServerDto(serverBo))
	}

	return serversDto
}

type CreateServerRequestDto struct {
	HostName string `json:"hostName"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	NewMap   bool   `json:"newMap"`
	MapName  string `json:"mapName,omitempty"`
	Memory   string `json:"memory"`
	Url      string `json:"url"`
}

type DeleteServerRequestDto struct {
	Full bool `json:"full"`
}
