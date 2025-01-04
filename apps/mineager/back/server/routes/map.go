package routes

import (
	"encoding/json"
	"mgarnier11/go/logger"
	"mgarnier11/mineager/server/controllers"
	"net/http"
	"path/filepath"
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

func serializeAndSendResponse(w http.ResponseWriter, response interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("Error serializing response: %v", err)
		http.Error(w, "Error serializing response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
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

type MapRequest struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
}

func postMap(w http.ResponseWriter, r *http.Request) {
	const maxUploadSize = 1 << 30 // 1 GB

	err := r.ParseMultipartForm(maxUploadSize)
	if err != nil {
		http.Error(w, "Failed to parse form data. File may be too large.", http.StatusBadRequest)
		return
	}

	mapRequest := MapRequest{
		Name:        strings.ToLower(r.FormValue("name")),
		Version:     r.FormValue("version"),
		Description: r.FormValue("description"),
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if filepath.Ext(fileHeader.Filename) != ".zip" {
		http.Error(w, "Only .zip files are allowed", http.StatusBadRequest)
		return
	}

	newMap, err := controllers.PostMap(mapRequest.Name, mapRequest.Version, mapRequest.Description, file)

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
