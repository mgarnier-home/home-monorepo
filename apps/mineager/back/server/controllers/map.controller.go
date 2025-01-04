package controllers

import (
	"fmt"
	"io"
	"mgarnier11/mineager/config"
	"mgarnier11/mineager/server/models"
	"mime/multipart"
	"os"
)

type MapDto struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func mapBoToMapDto(mapBo *models.MapBo) *MapDto {
	return &MapDto{
		Name:        mapBo.Name,
		Version:     mapBo.Version,
		Description: mapBo.Description,
	}
}

func mapsBoToMapsDto(mapsBo []*models.MapBo) []*MapDto {
	mapsDto := make([]*MapDto, 0)

	for _, mapBo := range mapsBo {
		mapsDto = append(mapsDto, mapBoToMapDto(mapBo))
	}

	return mapsDto
}

func GetMaps() ([]*MapDto, error) {
	maps, err := models.GetMaps()

	if err != nil {
		return nil, fmt.Errorf("error getting maps: %v", err)
	}

	return mapsBoToMapsDto(maps), nil
}

func GetMap(name string) (*MapDto, error) {
	mapRow, err := models.GetMapByName(name)

	if err != nil {
		return nil, fmt.Errorf("error getting map: %v", err)
	}

	return mapBoToMapDto(mapRow), nil
}

func PostMap(name string, version string, description string, file multipart.File) (*MapDto, error) {
	newMap, err := models.CreateMap(name, version, description)

	if err != nil {
		return nil, fmt.Errorf("error creating map: %v", err)
	}

	mapPath := fmt.Sprintf("%s/%s.zip", config.Config.MapsFolderPath, newMap.Name)

	dst, err := os.Create(mapPath)
	if err != nil {
		models.DeleteMapByName(newMap.Name)
		return nil, fmt.Errorf("error creating map file: %v", err)
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		models.DeleteMapByName(newMap.Name)
		return nil, fmt.Errorf("error copying file: %v", err)
	}

	return mapBoToMapDto(newMap), nil
}

func DeleteMap(name string) error {
	err := models.DeleteMapByName(name)

	if err != nil {
		return fmt.Errorf("error deleting map: %v", err)
	}

	mapPath := fmt.Sprintf("%s/%s.zip", config.Config.MapsFolderPath, name)

	err = os.Remove(mapPath)

	if err != nil {
		return fmt.Errorf("error deleting map file: %v", err)
	}

	return nil
}
