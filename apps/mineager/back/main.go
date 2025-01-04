package main

import (
	"mgarnier11/go/logger"

	_ "net/http/pprof"
)

func main() {
	logger.InitAppLogger("mineager")
}
