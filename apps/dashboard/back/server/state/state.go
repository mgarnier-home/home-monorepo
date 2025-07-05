package state

import (
	"sync/atomic"
	"time"

	"github.com/zishang520/socket.io/v2/socket"
	"mgarnier11.fr/go/dashboard/config"
	"mgarnier11.fr/go/libs/logger"
)

var clientCount int32
var stopRoutine chan struct{}
var routineRunning int32

func ClientConnect(
	logger *logger.Logger,
	io *socket.Server,
	client *socket.Socket,
	dashboardConfing *config.DashboardConfig,
) {
	if atomic.AddInt32(&clientCount, 1) == 1 {
		if atomic.CompareAndSwapInt32(&routineRunning, 0, 1) {
			stopRoutine = make(chan struct{})
			go watchState()
		}
	}
}

func ClientDisconnect() {
	if atomic.AddInt32(&clientCount, -1) == 0 && stopRoutine != nil {
		close(stopRoutine)
	}
}

func watchState() {
	logger.Infof("Background routine started")
	defer func() {
		logger.Infof("Background routine stopped")
		atomic.StoreInt32(&routineRunning, 0)
	}()
	for {
		select {
		case <-stopRoutine:
			return
		default:
			// Your background task here
			logger.Infof("Background routine is running, current client count: %d", atomic.LoadInt32(&clientCount))
			time.Sleep(60 * time.Second)
		}
	}
}
