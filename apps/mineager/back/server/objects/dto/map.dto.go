package dto

import "mgarnier11/mineager/server/objects/bo"

type MapDto struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func MapBoToMapDto(mapBo *bo.MapBo) *MapDto {
	return &MapDto{
		Name:        mapBo.Name,
		Version:     mapBo.Version,
		Description: mapBo.Description,
	}
}

func MapsBoToMapsDto(mapsBo []*bo.MapBo) []*MapDto {
	mapsDto := make([]*MapDto, 0)

	for _, mapBo := range mapsBo {
		mapsDto = append(mapsDto, MapBoToMapDto(mapBo))
	}

	return mapsDto
}

type CreateMapRequestDto struct {
	Name        string
	Version     string
	Description string
	File        *[]byte
}
