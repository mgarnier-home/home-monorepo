package utils

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-ping/ping"
	"gopkg.in/yaml.v3"
)

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func GenerateRandomData(size int) ([]byte, error) {
	data := make([]byte, size)
	_, err := rand.Read(data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func GenerateRandomString(size int) (string, error) {
	data, err := GenerateRandomData(size / 2)

	if err != nil {
		return "", err
	}

	return hex.EncodeToString(data), nil
}

type CustomWriter struct {
	io.Writer
	OnWrite func(int)
}

func (cw *CustomWriter) Write(p []byte) (int, error) {
	n, err := cw.Writer.Write(p)
	if err == nil {
		cw.OnWrite(n)
	}
	return n, err
}

func IsHTTPRequest(data []byte) bool {

	return bytes.Contains(data, []byte("HTTP"))
}

func CheckRequestHeader(request string, header string, value string) bool {
	return strings.Contains(request, header+": "+value)
}

func FilterFunc[S ~[]E, E any](s S, f func(E) bool) []E {
	var r []E
	for _, v := range s {
		if f(v) {
			r = append(r, v)
		}
	}
	return r
}

// CopyFolder copies a folder from src to dst
func CopyFolder(src string, dst string) error {
	// Walk through the source folder
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Build the destination path
		relativePath, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, relativePath)

		// Check if it's a directory
		if info.IsDir() {
			// Create the directory
			return os.MkdirAll(dstPath, info.Mode())
		}

		// Copy the file
		return CopyFile(path, dstPath)
	})
}

// CopyFile copies a single file from src to dst
func CopyFile(src string, dst string) error {
	// Open the source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Create the destination file
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy the contents
	_, err = io.Copy(dstFile, srcFile)
	if err != nil {
		return err
	}

	// Copy file permissions
	srcInfo, err := os.Stat(src)
	if err != nil {
		return err
	}
	return os.Chmod(dst, srcInfo.Mode())
}

func PingIp(ip string, timeout time.Duration) (bool, error) {
	pinger, err := ping.NewPinger(ip)

	if err != nil {
		return false, fmt.Errorf("failed to create pinger: %v", err)
	}

	pinger.Count = 1
	pinger.Timeout = timeout

	err = pinger.Run()
	if err != nil {
		return false, fmt.Errorf("failed to run pinger: %v", err)
	}

	return pinger.Statistics().PacketsRecv > 0, nil
}

func ReadYamlFile[T any](filePath string) (*T, error) {
	configFile, err := os.ReadFile(filePath)

	if err != nil {
		return nil, err
	}

	appConfig := new(T)

	err = yaml.Unmarshal(configFile, appConfig)

	if err != nil {
		return nil, err
	}

	return appConfig, nil
}

func GetDirSize(path string) (int64, error) {
	var size int64

	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		size += info.Size()
		return nil
	})

	if err != nil {
		return 0, err
	}

	return size, nil
}

func GetDateOfDay() string {
	return fmt.Sprintf("%d-%02d-%02d", time.Now().Year(), time.Now().Month(), time.Now().Day())
}

func RunPeriodic(ctx context.Context, interval time.Duration, fn func()) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			fn()
		}
	}
}

func GetAbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	executablePath, err := os.Executable()

	if err != nil {
		return ""
	}

	absolutePath := filepath.Join(filepath.Dir(executablePath), path)
	if !filepath.IsAbs(absolutePath) {
		return ""
	}
	return absolutePath
}
