package controllers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mgarnier11/go/logger"
	"mgarnier11/go/utils"
	"mgarnier11/mineager/server/database"
	"mgarnier11/mineager/server/objects/bo"

	"os"
	"path/filepath"
	"slices"
	"strings"
)

type MapController struct {
	mapRepository *database.MapRepository
}

func NewMapController() *MapController {
	return &MapController{
		mapRepository: database.CreateMapRepository(),
	}
}

func (controller *MapController) GetMaps() ([]*bo.MapBo, error) {
	maps, err := controller.mapRepository.GetMaps()

	if err != nil {
		return nil, fmt.Errorf("error getting maps: %v", err)
	}

	return maps, nil
}

func (controller *MapController) GetMap(name string) (*bo.MapBo, error) {
	mapRow, err := controller.mapRepository.GetMapByName(name)

	if err != nil {
		return nil, fmt.Errorf("error getting map: %v", err)
	}

	return mapRow, nil
}

func (controller *MapController) PostMap(name string, version string, description string, file *[]byte) (*bo.MapBo, error) {
	newMap, err := controller.mapRepository.CreateMap(name, version, description)

	if err != nil {
		return nil, fmt.Errorf("error creating map: %v", err)
	}

	sendError := func(err error) (*bo.MapBo, error) {
		controller.mapRepository.DeleteMapByName(newMap.Name)
		return nil, err
	}

	mapPath := getMapPath(newMap.Name)

	// Create zip reader
	zipReader, err := zip.NewReader(bytes.NewReader(*file), int64(len(*file)))
	if err != nil {
		return sendError(fmt.Errorf("failed to read ZIP file: %v", err))
	}

	// Get the index of the level.dat file
	levelsDatFileIndex := slices.IndexFunc(zipReader.File, func(file *zip.File) bool {
		return strings.Contains(file.Name, "level.dat")
	})

	if levelsDatFileIndex == -1 {
		return sendError(fmt.Errorf("missing level.dat file"))
	}

	levelsDatFile := zipReader.File[levelsDatFileIndex]

	// Get the fodler where are located all the map files
	worldFolder := filepath.Dir(levelsDatFile.Name)

	// Get all the files in the world folder
	worldFiles := utils.FilterFunc(zipReader.File, func(file *zip.File) bool {
		return strings.HasPrefix(file.Name, worldFolder)
	})

	// Create the map folder
	err = os.MkdirAll(mapPath, 0755)
	if err != nil {
		return sendError(fmt.Errorf("error creating map folder: %v", err))
	}

	for _, file := range worldFiles {
		logger.Infof("File: %s", file.Name)
	}

	// Extract all the files in the world folder
	for _, file := range worldFiles {
		filePath := filepath.Join(mapPath, strings.TrimPrefix(file.Name, worldFolder))

		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			return sendError(fmt.Errorf("error creating folder: %v", err))
		}

		if !file.FileInfo().IsDir() {
			// Create the file
			src, err := file.Open()
			if err != nil {
				return sendError(fmt.Errorf("error opening file: %v", err))
			}
			defer src.Close()

			osFile, err := os.Create(filePath)
			if err != nil {
				return sendError(fmt.Errorf("error creating file: %v", err))
			}
			defer osFile.Close()

			_, err = io.Copy(osFile, src)
			if err != nil {
				return sendError(fmt.Errorf("error copying file: %v", err))
			}
		}
	}

	logger.Infof("Map %s created", newMap.Name)

	return newMap, nil
}

func (controller *MapController) DeleteMap(name string) error {
	err := controller.mapRepository.DeleteMapByName(name)

	if err != nil {
		return fmt.Errorf("error deleting map: %v", err)
	}

	mapPath := getMapPath(name)

	os.RemoveAll(mapPath)

	return nil
}
