package utils

import (
	"bytes"
	"crypto/rand"
	"io"
	"strings"
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
