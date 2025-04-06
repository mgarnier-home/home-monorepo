package utils

import (
	"errors"
	"io"
	"os"
	"sync"
)

type ReadWriterAt interface {
	io.ReaderAt
	io.WriterAt
	io.Closer
	Stat() (os.FileInfo, error)
}

func copyChunk[srcType ReadWriterAt, dstType ReadWriterAt](srcFile srcType, dstFile dstType, offset int64, size int, wg *sync.WaitGroup, progressFunc func(int)) {
	defer wg.Done()

	const bufferSize = 32 * 1024 // 32 KB buffer size
	buffer := make([]byte, bufferSize)

	for remaining := size; remaining > 0; {
		readSize := bufferSize
		if remaining < bufferSize {
			readSize = remaining
		}

		n, err := srcFile.ReadAt(buffer[:readSize], offset)
		if err != nil && err != io.EOF {
			panic(err)
		}

		if n > 0 {
			_, err = dstFile.WriteAt(buffer[:n], offset)
			if err != nil {
				panic(err)
			}
			offset += int64(n)
			remaining -= n
		}

		if progressFunc != nil {
			progressFunc(n)
		}
	}
}

func copyFile[srcType ReadWriterAt, dstType ReadWriterAt](localFile srcType, remoteFile dstType, fileSize int64, progressFunc func(int)) {
	numChunks := 16 // Number of parallel chunks
	chunkSize := (fileSize + int64(numChunks) - 1) / int64(numChunks)

	var wg sync.WaitGroup
	for i := 0; i < numChunks; i++ {
		wg.Add(1)
		offset := int64(i) * chunkSize
		size := int(chunkSize)
		if offset+int64(size) > fileSize {
			size = int(fileSize - offset)
		}
		go copyChunk(localFile, remoteFile, offset, size, &wg, progressFunc)
	}

	wg.Wait()
}

func ParallelCopyFile(
	srcPath,
	dstPath string,
	openFile func(string) (ReadWriterAt, error),
	createFile func(string) (ReadWriterAt, error),
	progressFunc func(int, int64, int64),
) error {
	if openFile == nil || createFile == nil {
		return errors.New("open and create functions cannot be nil")
	}

	srcHandle, err := openFile(srcPath)
	if err != nil {
		return err
	}
	defer srcHandle.Close()

	dstHandle, err := createFile(dstPath)
	if err != nil {
		return err
	}
	defer dstHandle.Close()

	fileInfo, err := srcHandle.Stat()
	if err != nil {
		return err
	}

	totalSize := fileInfo.Size()
	copiedSize := int64(0)
	fileProgress := func(n int) {
		copiedSize += int64(n)

		if progressFunc != nil {
			progressFunc(n, copiedSize, totalSize)
		}
	}

	copyFile(srcHandle, dstHandle, fileInfo.Size(), fileProgress)

	return nil
}
