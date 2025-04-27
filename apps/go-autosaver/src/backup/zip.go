package backup

import (
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"

	"mgarnier11.fr/go/libs/utils"
)

// zipFolder creates a zip archive of the given folder.
func ZipFolder(
	folderPath,
	zipPath string,
	progressFunc func(
		fileName string,
		written int,
		fileWritten int64,
		fileSize int64,
		totalWritten int64,
		totalSize int64,
	),
) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return fmt.Errorf("failed to create zip file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	totalWritten := int64(0)
	totalSize, err := utils.GetDirSize(folderPath)
	if err != nil {
		return fmt.Errorf("failed to get folder size: %w", err)
	}

	err = filepath.Walk(folderPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk through folder: %w", err)
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return fmt.Errorf("failed to create zip header: %w", err)
		}

		header.Name, _ = filepath.Rel(filepath.Clean(folderPath), filePath)
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return fmt.Errorf("failed to create zip writer: %w", err)
		}

		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()
			fileSize := info.Size()
			_, err = utils.CopyWithProgress(writer, file, func(written int, fileWritten int64) {
				totalWritten += int64(written)
				if progressFunc != nil {
					progressFunc(
						filePath,
						written,
						fileWritten,
						fileSize,
						totalWritten,
						totalSize,
					)
				}
			})

			if err != nil {
				return fmt.Errorf("failed to write file to zip: %w", err)
			}
		}
		return nil
	})

	if err != nil {
		return fmt.Errorf("failed to walk through folder: %w", err)
	}

	return nil
}
