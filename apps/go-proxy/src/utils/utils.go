package utils

import (
	"bytes"
	"strings"
)

func IsHTTPRequest(data []byte) bool {
	return bytes.Contains(data, []byte("HTTP"))
}

func CheckRequestHeader(request string, header string, value string) bool {
	return strings.Contains(request, header+": "+value)
}
