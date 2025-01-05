package routes

import (
	"mgarnier11/go/logger"
	"mgarnier11/mineager/server/controllers"
	"mgarnier11/mineager/server/routes/validation"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func MapRoutes(router *mux.Router) {
	mapRouter := router.PathPrefix("/map").Subrouter()

	mapRouter.HandleFunc("/", getMaps).Methods("GET")
	mapRouter.HandleFunc("/{name}", getMap).Methods("GET")
	mapRouter.HandleFunc("/", postMap).Methods("POST")
	mapRouter.HandleFunc("/{name}", deleteMap).Methods("DELETE")
}

func getMaps(w http.ResponseWriter, r *http.Request) {
	maps, err := controllers.GetMaps()

	if err != nil {
		logger.Errorf("Error getting maps: %v", err)
		http.Error(w, "Error getting maps", http.StatusInternalServerError)
		return
	}

	serializeAndSendResponse(w, maps)
}

func getMap(w http.ResponseWriter, r *http.Request) {
	mapName := strings.ToLower(mux.Vars(r)["name"])

	mapBo, err := controllers.GetMap(mapName)

	if err != nil {
		logger.Errorf("Error getting map: %v", err)
		http.Error(w, "Error getting map", http.StatusInternalServerError)
		return
	}

	serializeAndSendResponse(w, mapBo)
}

func postMap(w http.ResponseWriter, r *http.Request) {
	requestValidated, err := validation.ValidateMapPostRequest(r)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMap, err := controllers.PostMap(
		requestValidated.Name,
		requestValidated.Version,
		requestValidated.Description,
		requestValidated.File,
	)

	if err != nil {
		logger.Errorf("Error creating map: %v", err)
		http.Error(w, "Error creating map", http.StatusInternalServerError)
		return
	}

	serializeAndSendResponse(w, newMap)
}

func deleteMap(w http.ResponseWriter, r *http.Request) {
	mapName := strings.ToLower(mux.Vars(r)["name"])

	err := controllers.DeleteMap(mapName)

	if err != nil {
		logger.Errorf("Error deleting map: %v", err)
		http.Error(w, "Error deleting map", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
