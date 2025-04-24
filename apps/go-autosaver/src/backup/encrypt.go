package backup

import (
	"os"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"mgarnier11.fr/go/libs/utils"
)

func EncryptFileWithPassword(inputPath, outputPath, password string, progressCallback func(int, int64, int64)) error {
	messageReader, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer messageReader.Close()

	ciphertextWriter, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer ciphertextWriter.Close()

	passwordBytes := []byte(password)
	pgp := crypto.PGP()

	encHandle, err := pgp.Encryption().Password(passwordBytes).New()
	if err != nil {
		return err
	}
	writer, err := encHandle.EncryptingWriter(ciphertextWriter, crypto.Auto)
	if err != nil {
		return err
	}
	defer writer.Close()

	// Get the size of the input file for progress calculation
	fileInfo, err := messageReader.Stat()
	if err != nil {
		return err
	}
	totalSize := fileInfo.Size()

	utils.CopyWithProgress(writer, messageReader, func(written int, total int64) {
		if progressCallback != nil {
			progressCallback(written, total, totalSize)
		}
	})

	return nil
}
