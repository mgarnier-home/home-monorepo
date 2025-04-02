package dto

import "mgarnier11.fr/go/mineager/server/objects/bo"

type MapDto struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func MapMapBoToMapDto(mapBo *bo.MapBo) *MapDto {
	return &MapDto{
		Name:        mapBo.Name,
		Version:     mapBo.Version,
		Description: mapBo.Description,
	}
}

func MapMapsBoToMapsDto(mapsBo []*bo.MapBo) []*MapDto {
	mapsDto := make([]*MapDto, 0)

	for _, mapBo := range mapsBo {
		mapsDto = append(mapsDto, MapMapBoToMapDto(mapBo))
	}

	return mapsDto
}

type CreateMapRequestDto struct {
	Name        string
	Version     string
	Description string
	File        *[]byte
}
