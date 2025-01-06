package routes

import (
	"mgarnier11/go/logger"
	"mgarnier11/mineager/server/controllers"
	"mgarnier11/mineager/server/routes/validation"
	"net/http"

	"github.com/gorilla/mux"
)

func ServerRoutes(router *mux.Router) {
	serverRouter := router.PathPrefix("/server").Subrouter()

	serverRouter.HandleFunc("/", getServers).Methods("GET")
	serverRouter.HandleFunc("/{name}", getServer).Methods("GET")
	serverRouter.HandleFunc("/{name}/start", startServer).Methods("POST")
	serverRouter.HandleFunc("/{name}/stop", stopServer).Methods("POST")
	serverRouter.HandleFunc("/", postServer).Methods("POST")
	serverRouter.HandleFunc("/{name}", deleteServer).Methods("DELETE")
}

func getServers(w http.ResponseWriter, r *http.Request) {
	servers, err := controllers.GetServers("", "")

	if err != nil {
		logger.Errorf("Error getting servers: %v", err)
		http.Error(w, "Error getting servers", http.StatusInternalServerError)
		return
	}

	serializeAndSendResponse(w, servers)
}

func getServer(w http.ResponseWriter, r *http.Request) {
	serverName := mux.Vars(r)["name"]

	servers, err := controllers.GetServers("", serverName)

	if err != nil {
		logger.Errorf("Error getting server: %v", err)
		http.Error(w, "Error getting server", http.StatusInternalServerError)
		return
	}

	serializeAndSendResponse(w, servers[0])
}

func postServer(w http.ResponseWriter, r *http.Request) {
	requestValidated, err := validation.ValidateServerPostRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	requestedMap, err := controllers.GetMap(requestValidated.MapName)

	if err != nil {
		logger.Errorf("Map %s not found: %v", requestValidated.MapName, err)
		http.Error(w, "mapName not found", http.StatusBadRequest)
		return
	}

	if requestValidated.Version == "" {
		requestValidated.Version = requestedMap.Version
	}

	newServer, err := controllers.CreateServer(
		requestValidated.HostName,
		requestValidated.Name,
		requestValidated.Version,
		requestedMap.Name,
		requestValidated.Memory,
		requestValidated.Url,
	)

	if err != nil {
		logger.Errorf("Error creating server: %v", err)
		http.Error(w, "Error creating server", http.StatusInternalServerError)
		return
	}

	serializeAndSendResponse(w, newServer)
}

func deleteServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete server"))
}

func startServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Start server"))
}

func stopServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Stop server"))
}
