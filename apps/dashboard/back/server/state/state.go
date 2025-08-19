package state

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/zishang520/socket.io/v2/socket"
	"mgarnier11.fr/go/dashboard/config"
	"mgarnier11.fr/go/libs/logger"
)

type HealthCheckState struct {
	Code    int    `yaml:"code" json:"code"`
	Message string `yaml:"message" json:"message"`
}

type HostState struct {
	Name string `yaml:"name" json:"name"`
	Ping int    `yaml:"ping" json:"ping"`
	Up   bool   `yaml:"up" json:"up"`
}

type ServiceState struct {
	Name         string                       `yaml:"name" json:"name"`
	DockerHealth string                       `yaml:"dockerHealth" json:"dockerHealth"`
	Up           bool                         `yaml:"up" json:"up"`
	HealthChecks map[string]*HealthCheckState `yaml:"healthChecks" json:"healthChecks"`
}

var clientCount int32
var stopRoutine chan struct{}
var routineRunning int32

func ClientConnect(
	logger *logger.Logger,
	io *socket.Server,
	client *socket.Socket,
) {
	if atomic.AddInt32(&clientCount, 1) == 1 {
		if atomic.CompareAndSwapInt32(&routineRunning, 0, 1) {
			stopRoutine = make(chan struct{})
			// First client connected, start the background routine
			go watchState(io)
		}
	}
}

func ClientDisconnect() {
	if atomic.AddInt32(&clientCount, -1) == 0 && stopRoutine != nil {
		close(stopRoutine)
	}
}

func watchState(io *socket.Server) {
	logger.Infof("Background routine started")
	defer func() {
		logger.Infof("Background routine stopped")
		atomic.StoreInt32(&routineRunning, 0)
	}()

	dashboardConfig, err := config.Config.GetDashboardConfig()
	if err != nil {
		logger.Errorf("Error loading dashboardConfig: %v", err)
		return
	}

	for _, host := range dashboardConfig.Hosts {
		go watchHostState(host, io)

		for _, service := range host.Services {
			go watchServiceState(service, io)
		}
	}

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

func watchHostState(host *config.Host, io *socket.Server) {
	logger.Infof("Watching host: %s", host.Name)

	hostState := &HostState{
		Name: host.Name,
		Ping: 0,     // Initialize ping value
		Up:   false, // Initialize host state
	}

	io.Emit("hostState", hostState)

	for {
		select {
		case <-stopRoutine:
			return
		default:
			// Here you can implement the logic to check the host state
			logger.Infof("Host %s is up", host.Name)

			time.Sleep(30 * time.Second) // Adjust the interval as needed
		}
	}
}

func watchServiceState(service *config.Service, io *socket.Server) {
	logger.Infof("Watching service: %s on host: %s", service.Name, service.DockerName)

	healthChecks := make(map[string]*HealthCheckState)
	for _, healthCheck := range service.HealthChecks {
		healthCheckName := fmt.Sprintf("%s.%s", service.Name, healthCheck.Name)
		healthChecks[healthCheckName] = &HealthCheckState{
			Code:    0,    // Initialize health check code
			Message: "OK", // Initialize health check message
		}
	}

	serviceState := &ServiceState{
		Name:         service.Name,
		DockerHealth: "unknown", // Initialize Docker health state
		Up:           false,     // Initialize service state
		HealthChecks: healthChecks,
	}

	io.Emit("serviceState", serviceState)

	for {
		select {
		case <-stopRoutine:
			return
		default:
			// Here you can implement the logic to check the service state
			logger.Infof("Service %s is up", service.Name)

			time.Sleep(30 * time.Second) // Adjust the interval as needed
		}
	}
}

// func getApplicationState(
// 	dashboardConfig *config.DashboardConfig,
// ) *ApplicationState {
// 	health := make(map[string]*HealthCheckState)

// 	for _, host := range dashboardConfig.Hosts {
// 		for _, service := range host.Services {
// 			for _, healthCheck := range service.HealthChecks {
// 				healthCheckName := fmt.Sprintf("%s.%s.%s", host.Name, service.Name, healthCheck.Name)

// 				health[healthCheckName] = &HealthCheckState{
// 					Code:    0,
// 					Message: "OK",
// 				}
// 			}
// 		}
// 	}

// 	return &ApplicationState{
// 		Config: dashboardConfig,
// 		Health: health,
// 	}
// }
