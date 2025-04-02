package dto

import "mgarnier11.fr/go/mineager/server/objects/bo"

type ServerDto struct {
	Name    string `json:"name"`
	Url     string `json:"url"`
	Version string `json:"version"`
	Map     string `json:"map"`
	Status  string `json:"status"`
}

func MapServerBoToServerDto(serverBo *bo.ServerBo) *ServerDto {
	return &ServerDto{
		Name:    serverBo.Name,
		Url:     serverBo.Url,
		Version: serverBo.Version,
		Map:     serverBo.Map,
		Status:  serverBo.Status,
	}
}

func MapServersBoToServersDto(serversBo []*bo.ServerBo) []*ServerDto {
	serversDto := make([]*ServerDto, 0)

	for _, serverBo := range serversBo {
		serversDto = append(serversDto, MapServerBoToServerDto(serverBo))
	}

	return serversDto
}

type CreateServerRequestDto struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	NewMap  bool   `json:"newMap"`
	MapName string `json:"mapName,omitempty"`
	Memory  string `json:"memory"`
}

type DeleteServerRequestDto struct {
	Full bool `json:"full"`
}
