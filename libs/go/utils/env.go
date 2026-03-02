package utils

import (
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

func InitEnvFromFile() {
	envFilePath := GetEnv("ENV_FILE_PATH", "./.env")

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}

	exPath := path.Dir(ex)

	if !path.IsAbs(envFilePath) {
		envFilePath = path.Join(exPath, envFilePath)
	}

	godotenv.Load(envFilePath)
}

func GetEnv[T bool | string | int](key string, defaultValue T) T {
	value := os.Getenv(key)

	if value == "" {
		if _, err := os.Stat("/run/secrets/" + key); err == nil {
			fileContent, err := os.ReadFile("/run/secrets/" + strings.ToLower(key))

			if err != nil {
				value = ""
			} else {
				value = string(fileContent)
			}
		}
	}

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
