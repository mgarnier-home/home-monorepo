package main

import (
	"mgarnier11/go-proxy/config"
	"mgarnier11/go-proxy/hostManager"
	"mgarnier11/go-proxy/server"
	"mgarnier11/go/logger"
	"runtime"
	"time"

	_ "net/http/pprof"
)

func main() {
	logger.InitAppLogger("")

	go func() {
		for {
			time.Sleep(time.Second * 5)

			// Analyzing goroutine leaks
			var stats runtime.MemStats
			runtime.ReadMemStats(&stats)
			logger.Infof("Number of Goroutines: %d", runtime.NumGoroutine())

			// buf := make([]byte, 1<<16) // Create a large buffer to capture stack traces
			// stackLen := runtime.Stack(buf, true)
			// fmt.Printf("=== Goroutine Stack Dump ===\n%s\n", buf[:stackLen])
		}
	}()

	server := server.NewServer(config.Config.ServerPort)

	go server.Start()

	for configFile := range config.SetupConfigListener() {
		hostManager.ConfigFileChanged(configFile)
	}
}
