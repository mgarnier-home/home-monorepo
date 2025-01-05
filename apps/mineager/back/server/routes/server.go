package routes

import (
	"mgarnier11/mineager/server/models"
	"net/http"

	"github.com/gorilla/mux"
)

func ServerRoutes(router *mux.Router) {
	serverRouter := router.PathPrefix("/server").Subrouter()

	serverRouter.HandleFunc("/", getServers).Methods("GET")
	serverRouter.HandleFunc("/{name}", getServer).Methods("GET")
	serverRouter.HandleFunc("/", postServer).Methods("POST")
	serverRouter.HandleFunc("/{name}", deleteServer).Methods("DELETE")
}

func getServers(w http.ResponseWriter, r *http.Request) {

	models.GetServers()
	w.Write([]byte("Get all servers"))
}

func getServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Get server"))
}

func postServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Post server"))
}

func deleteServer(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Delete server"))
}
