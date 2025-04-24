package backup

import (
	"archive/zip"
	"os"
	"path/filepath"

	"mgarnier11.fr/go/libs/utils"
)

// zipFolder creates a zip archive of the given folder.
func ZipFolder(folderPath, zipPath string, progressFunc func(string, int, int64, int64)) error {
	zipFile, err := os.Create(zipPath)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	err = filepath.Walk(folderPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name, _ = filepath.Rel(filepath.Clean(folderPath), filePath)
		if info.IsDir() {
			header.Name += "/"
		} else {
			header.Method = zip.Deflate
		}

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		if !info.IsDir() {
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()
			fileSize := info.Size()
			_, err = utils.CopyWithProgress(writer, file, func(written int, total int64) {
				if progressFunc != nil {
					progressFunc(filePath, written, total, fileSize)
				}
			})

			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
