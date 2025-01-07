package routes

import (
	"context"
	"mgarnier11/go/logger"
	"mgarnier11/mineager/server/controllers"
	"mgarnier11/mineager/server/objects/bo"
	"mgarnier11/mineager/server/objects/dto"
	"mgarnier11/mineager/server/routes/validation"
	"net/http"

	"github.com/gorilla/mux"
)

const serversContextKey contextKey = "servers"
const serverContextKey contextKey = "server"

func getServerControllerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller, err := controllers.NewServerController(mux.Vars(r)["hostName"])

		if err != nil {
			logger.Errorf("Error getting server controller: %v", err)
			sendErrorResponse(w, "Error getting server controller", http.StatusInternalServerError)
			return
		}

		defer controller.Dispose()

		ctx := context.WithValue(r.Context(), serversContextKey, controller)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getServerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller := r.Context().Value(serversContextKey).(*controllers.ServerController)
		server, err := controller.GetServer(mux.Vars(r)["name"])

		if err != nil {
			logger.Errorf("Error getting server: %v", err)
			sendErrorResponse(w, "Error getting server", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), serverContextKey, server)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ServerRoutes(router *mux.Router) {
	serverRouter := router.PathPrefix("/{hostName}/servers").Subrouter()
	serverRouter.Use(getServerControllerMiddleware)

	serverRouter.HandleFunc("", getServers).Methods("GET")
	serverRouter.HandleFunc("", postServer).Methods("POST")

	serverRouter = serverRouter.PathPrefix("/{name}").Subrouter()
	serverRouter.Use(getServerMiddleware)

	serverRouter.HandleFunc("", getServer).Methods("GET")
	serverRouter.HandleFunc("/start", startServer).Methods("POST")
	serverRouter.HandleFunc("/stop", stopServer).Methods("POST")
	serverRouter.HandleFunc("", deleteServer).Methods("DELETE")
}

func getServers(w http.ResponseWriter, r *http.Request) {
	controller := r.Context().Value(serversContextKey).(*controllers.ServerController)

	servers, err := controller.GetServers()

	if err != nil {
		logger.Errorf("Error getting servers: %v", err)
		http.Error(w, "Error getting servers", http.StatusInternalServerError)
	} else {
		serializeAndSendResponse(w, dto.ServersBoToServersDto(servers), http.StatusOK)
	}
}

func getServer(w http.ResponseWriter, r *http.Request) {
	server := r.Context().Value(serverContextKey).(*bo.ServerBo)

	serializeAndSendResponse(w, dto.ServerBoToServerDto(server), http.StatusOK)
}

func postServer(w http.ResponseWriter, r *http.Request) {
	controller := r.Context().Value(serversContextKey).(*controllers.ServerController)

	createServerDto, err := validation.ValidateServerPostRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newServer, err := controller.CreateServer(createServerDto)

	if err != nil {
		logger.Errorf("Error creating server: %v", err)
		sendErrorResponse(w, "Error creating server", http.StatusInternalServerError)
	} else {
		serializeAndSendResponse(w, dto.ServerBoToServerDto(newServer), http.StatusOK)
	}
}

func deleteServer(w http.ResponseWriter, r *http.Request) {
	controller := r.Context().Value(serversContextKey).(*controllers.ServerController)
	server := r.Context().Value(serverContextKey).(*bo.ServerBo)

	deleteServerDto, err := validation.ValidateServerDeleteRequest(r)

	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = controller.DeleteServer(server, deleteServerDto.Full)

	if err != nil {
		logger.Errorf("Error deleting server: %v", err)
		sendErrorResponse(w, "Error deleting server", http.StatusInternalServerError)
	} else {
		sendOKResponse(w, "Server deleted")
	}
}

func startServer(w http.ResponseWriter, r *http.Request) {
	controller := r.Context().Value(serversContextKey).(*controllers.ServerController)
	server := r.Context().Value(serverContextKey).(*bo.ServerBo)

	err := controller.StartServer(server)

	if err != nil {
		logger.Errorf("Error starting server: %v", err)
		sendErrorResponse(w, "Error starting server", http.StatusInternalServerError)
	} else {
		sendOKResponse(w, "Server started")
	}
}

func stopServer(w http.ResponseWriter, r *http.Request) {
	controller := r.Context().Value(serversContextKey).(*controllers.ServerController)
	server := r.Context().Value(serverContextKey).(*bo.ServerBo)

	err := controller.StopServer(server)

	if err != nil {
		logger.Errorf("Error stopping server: %v", err)
		sendErrorResponse(w, "Error stopping server", http.StatusInternalServerError)
	} else {
		sendOKResponse(w, "Server stopped")
	}
}
