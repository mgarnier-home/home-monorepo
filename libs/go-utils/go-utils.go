package goUtils

import (
	"crypto/rand"
	"io"
	"os"
	"strconv"
)

func GetEnv[T bool | string | int](key string, defaultValue T) T {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}

	switch any(defaultValue).(type) {
	case bool:
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return any(boolValue).(T)
		}
	case int:
		if intValue, err := strconv.Atoi(value); err == nil {
			return any(intValue).(T)
		}
	case string:
		return any(value).(T)
	}
	return defaultValue
}

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
