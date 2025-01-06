package controllers

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"mgarnier11/go/logger"
	"mgarnier11/go/sshutils"
	"mgarnier11/go/sshutils/sftp"
	"mgarnier11/go/utils"
	"mgarnier11/mineager/config"
	"mgarnier11/mineager/server/models"
	"net"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
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

func getMapPath(mapName string) string {
	return fmt.Sprintf("%s/%s", config.Config.MapsFolderPath, mapName)
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

func PostMap(name string, version string, description string, file *[]byte) (*MapDto, error) {
	newMap, err := models.CreateMap(name, version, description)

	if err != nil {
		return nil, fmt.Errorf("error creating map: %v", err)
	}

	sendError := func(err error) (*MapDto, error) {
		models.DeleteMapByName(newMap.Name)
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

	return mapBoToMapDto(newMap), nil
}

func DeleteMap(name string) error {
	err := models.DeleteMapByName(name)

	if err != nil {
		return fmt.Errorf("error deleting map: %v", err)
	}

	mapPath := getMapPath(name)

	os.RemoveAll(mapPath)

	return nil
}

func sendMapToHost(serverName string, mapName string, host *config.DockerHostConfig) error {
	sshAuth, err := sshutils.GetSSHKeyAuth(config.Config.SSHKeyPath)

	if err != nil {
		return fmt.Errorf("error getting ssh key auth: %v", err)
	}

	sshClient, err := ssh.Dial("tcp", net.JoinHostPort(host.Ip, host.SSHPort), &ssh.ClientConfig{
		User:            host.SSHUsername,
		Auth:            []ssh.AuthMethod{sshAuth},
		Timeout:         5 * time.Second,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})

	if err != nil {
		return fmt.Errorf("error connecting to ssh: %v", err)
	}
	defer sshClient.Close()

	session, err := sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("error creating session: %v", err)
	}

	command := fmt.Sprintf("mkdir -p %s/%s/world", host.MineagerPath, serverName)

	logger.Infof("Command: %s", command)

	err = session.Run(command)
	if err != nil {
		return fmt.Errorf("error creating server folder: %v", err)
	}

	session.Close()

	serverDestPath := models.GetServerDestPath(host, serverName)

	return sftp.LocalToRemote(sshClient, getMapPath(mapName), serverDestPath)
}
