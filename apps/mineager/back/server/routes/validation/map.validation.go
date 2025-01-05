package validation

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"mgarnier11/go/logger"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

type postMapRequest struct {
	Name        string
	Version     string
	Description string
	File        *[]byte
}

func validateMinecraftMap(fileBuffer *[]byte) error {
	zipReader, err := zip.NewReader(bytes.NewReader(*fileBuffer), int64(len(*fileBuffer)))
	if err != nil {
		return fmt.Errorf("failed to read ZIP file: %v", err)
	}

	containsFunc := func(find string) bool {
		return slices.ContainsFunc(zipReader.File, func(file *zip.File) bool {
			return strings.Contains(file.Name, find)
		})
	}

	if !containsFunc("level.dat") {
		return errors.New("missing level.dat file")
	}

	if !containsFunc("region") {
		return errors.New("missing region folder")
	}

	if !containsFunc("data") {
		return errors.New("missing data folder")
	}

	return nil
}

func ValidateMapPostRequest(r *http.Request) (*postMapRequest, error) {
	const maxUploadSize = 1 << 30 // 1 GB

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		return nil, errors.New("failed to parse form data, file may be too large")
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, errors.New("failed to get file from request")
	}
	defer file.Close()

	if filepath.Ext(fileHeader.Filename) != ".zip" {
		return nil, errors.New("only .zip files are allowed")
	}

	buf := make([]byte, fileHeader.Size)
	i, err := file.Read(buf)
	if err != nil {
		return nil, errors.New("failed to read file")
	}

	logger.Infof("Read %d bytes from file", i)

	contentType := http.DetectContentType(buf)
	if contentType != "application/zip" && contentType != "application/x-zip-compressed" {
		return nil, errors.New("invalid file type")
	}

	if err := validateMinecraftMap(&buf); err != nil {
		return nil, err
	}

	name := strings.ToLower(r.FormValue("name"))
	if err := validateName(name, "name"); err != nil {
		return nil, err // Return the error directly
	}

	version := r.FormValue("version")
	if err := validateVersion(version, "version", true); err != nil {
		return nil, err
	}

	description := r.FormValue("description")
	if len(description) > 500 {
		return nil, errors.New("description must be less than 500 characters")
	}

	return &postMapRequest{
		Name:        name,
		Version:     version,
		Description: description,
		File:        &buf,
	}, nil
}
