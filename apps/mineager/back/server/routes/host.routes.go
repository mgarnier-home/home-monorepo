package routes

import (
	"mgarnier11/mineager/config"
	"mgarnier11/mineager/server/objects/dto"
	"net/http"

	"github.com/gorilla/mux"
)

func HostRoutes(router *mux.Router) {
	hostRouter := router.PathPrefix("/hosts").Subrouter()

	hostRouter.HandleFunc("", getHosts).Methods("GET")
}

func getHosts(w http.ResponseWriter, r *http.Request) {
	serializeAndSendResponse(w, dto.DockerHostsToHostsDto(config.Config.AppConfig.DockerHosts), http.StatusOK)
}
