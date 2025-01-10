package routers

import (
	"context"
	"mgarnier11/go/logger"
	mineagerutils "mgarnier11/mineager/mineager-utils"
	"mgarnier11/mineager/server/controllers"
	"mgarnier11/mineager/server/objects/bo"
	"mgarnier11/mineager/server/objects/dto"
	"mgarnier11/mineager/server/routers/validation"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/mux"
)

type ServersRouter struct {
	utils             RouterUtils
	hostsController   *controllers.HostsController
	serversController *controllers.ServersController
}

func NewServersRouter(router *mux.Router, serverLogger *logger.Logger) *ServersRouter {
	serverRouter := &ServersRouter{
		utils: RouterUtils{
			logger: logger.NewLogger(
				"[SERVERS]",
				"%-10s ",
				lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
				serverLogger,
			),
		},
		hostsController:   controllers.NewHostsController(),
		serversController: controllers.NewServersController(),
	}

	hostsRouter := router.PathPrefix("/hosts/{hostName}").Subrouter()
	hostsRouter.Use(serverRouter.getServersControllerMiddleware)
	hostsRouter.HandleFunc("/servers", serverRouter.getServers).Methods("GET")
	hostsRouter.HandleFunc("/servers", serverRouter.postServer).Methods("POST")

	serversRouter := hostsRouter.PathPrefix("/servers/{serverName}").Subrouter()
	serversRouter.Use(serverRouter.getServerMiddleware)
	serversRouter.HandleFunc("", serverRouter.getServer).Methods("GET")
	serversRouter.HandleFunc("", serverRouter.deleteServer).Methods("DELETE")
	serversRouter.HandleFunc("/start", serverRouter.startServer).Methods("POST")
	serversRouter.HandleFunc("/stop", serverRouter.stopServer).Methods("POST")

	return serverRouter
}

func getControllerFromContext(ctx context.Context) *controllers.ServersController {
	return ctx.Value(mineagerutils.ServerControllerKey).(*controllers.ServersController)
}

func getServerFromContext(ctx context.Context) *bo.ServerBo {
	return ctx.Value(mineagerutils.ServerKey).(*bo.ServerBo)
}

func (router *ServersRouter) getServersControllerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		host, err := router.hostsController.GetHost(mux.Vars(r)["hostName"])

		if err != nil {
			router.utils.logger.Errorf("Error getting host: %v", err)
			router.utils.sendErrorResponse(w, r, "Error getting host", http.StatusInternalServerError)
			return
		}

		if !host.Ping {
			router.utils.sendErrorResponse(w, r, "Host is not reachable", http.StatusInternalServerError)
			return
		}

		serversController, err := router.serversController.WithHost(host)

		if err != nil {
			router.utils.sendErrorResponse(w, r, "Error getting servers controller", http.StatusInternalServerError)
		}

		defer serversController.Dispose()

		ctx := context.WithValue(r.Context(), mineagerutils.ServerControllerKey, serversController)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (router *ServersRouter) getServerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller := getControllerFromContext(r.Context())

		server, err := controller.GetServer(mux.Vars(r)["serverName"])

		if err != nil {
			router.utils.logger.Errorf("Error getting server: %v", err)
			router.utils.sendErrorResponse(w, r, "Error getting server", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), mineagerutils.ServerKey, server)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (router *ServersRouter) getServers(w http.ResponseWriter, r *http.Request) {
	controller := getControllerFromContext(r.Context())

	servers, err := controller.GetServers()

	if err != nil {
		router.utils.logger.Errorf("Error getting servers: %v", err)
		router.utils.sendErrorResponse(w, r, "Error getting servers", http.StatusInternalServerError)
	} else {
		router.utils.serializeAndSendResponse(w, r, dto.MapServersBoToServersDto(servers), http.StatusOK)
	}
}

func (router *ServersRouter) getServer(w http.ResponseWriter, r *http.Request) {
	server := getServerFromContext(r.Context())

	router.utils.serializeAndSendResponse(w, r, dto.MapServerBoToServerDto(server), http.StatusOK)
}

func (router *ServersRouter) postServer(w http.ResponseWriter, r *http.Request) {
	controller := getControllerFromContext(r.Context())

	createServerDto, err := validation.ValidateServerPostRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newServer, err := controller.CreateServer(createServerDto)

	if err != nil {
		router.utils.logger.Errorf("Error creating server: %v", err)
		router.utils.sendErrorResponse(w, r, "Error creating server", http.StatusInternalServerError)
	} else {
		router.utils.serializeAndSendResponse(w, r, dto.MapServerBoToServerDto(newServer), http.StatusOK)
	}
}

func (router *ServersRouter) deleteServer(w http.ResponseWriter, r *http.Request) {
	controller := getControllerFromContext(r.Context())
	server := getServerFromContext(r.Context())

	deleteServerDto, err := validation.ValidateServerDeleteRequest(r)

	if err != nil {
		router.utils.sendErrorResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.DeleteServer(server, deleteServerDto.Full)

	if err != nil {
		router.utils.logger.Errorf("Error deleting server: %v", err)
		router.utils.sendErrorResponse(w, r, "Error deleting server", http.StatusInternalServerError)
	} else {
		router.utils.sendOKResponse(w, r, "Server deleted")
	}
}

func (router *ServersRouter) startServer(w http.ResponseWriter, r *http.Request) {
	controller := getControllerFromContext(r.Context())
	server := getServerFromContext(r.Context())

	err := controller.StartServer(server)

	if err != nil {
		router.utils.logger.Errorf("Error starting server: %v", err)
		router.utils.sendErrorResponse(w, r, "Error starting server", http.StatusInternalServerError)
	} else {
		router.utils.sendOKResponse(w, r, "Server started")
	}
}

func (router *ServersRouter) stopServer(w http.ResponseWriter, r *http.Request) {
	controller := getControllerFromContext(r.Context())
	server := getServerFromContext(r.Context())

	err := controller.StopServer(server)

	if err != nil {
		router.utils.logger.Errorf("Error stopping server: %v", err)
		router.utils.sendErrorResponse(w, r, "Error stopping server", http.StatusInternalServerError)
	} else {
		router.utils.sendOKResponse(w, r, "Server stopped")
	}
}
