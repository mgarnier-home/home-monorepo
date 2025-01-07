package routes

import (
	"context"
	"mgarnier11/go/logger"
	"mgarnier11/mineager/server/controllers"
	"mgarnier11/mineager/server/objects/dto"
	"mgarnier11/mineager/server/routes/validation"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

const mapsContextKey contextKey = "maps"

func getMapControllerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller := controllers.NewMapController()

		ctx := context.WithValue(r.Context(), mapsContextKey, controller)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MapRoutes(router *mux.Router) {
	mapRouter := router.PathPrefix("/maps").Subrouter()
	mapRouter.Use(getMapControllerMiddleware)

	mapRouter.HandleFunc("", getMaps).Methods("GET")
	mapRouter.HandleFunc("/{name}", getMap).Methods("GET")
	mapRouter.HandleFunc("", postMap).Methods("POST")
	mapRouter.HandleFunc("/{name}", deleteMap).Methods("DELETE")
}

func getMaps(w http.ResponseWriter, r *http.Request) {
	controller := r.Context().Value(mapsContextKey).(*controllers.MapController)

	maps, err := controller.GetMaps()

	if err != nil {
		logger.Errorf("Error getting maps: %v", err)
		sendErrorResponse(w, "Error getting maps", http.StatusInternalServerError)
	} else {
		serializeAndSendResponse(w, dto.MapsBoToMapsDto(maps), http.StatusOK)
	}
}

func getMap(w http.ResponseWriter, r *http.Request) {
	controller := r.Context().Value(mapsContextKey).(*controllers.MapController)

	mapName := strings.ToLower(mux.Vars(r)["name"])

	mapBo, err := controller.GetMap(mapName)

	if err != nil {
		logger.Errorf("Error getting map: %v", err)
		sendErrorResponse(w, "Error getting map", http.StatusInternalServerError)
	} else {
		serializeAndSendResponse(w, dto.MapBoToMapDto(mapBo), http.StatusOK)
	}
}

func postMap(w http.ResponseWriter, r *http.Request) {
	controller := r.Context().Value(mapsContextKey).(*controllers.MapController)

	requestValidated, err := validation.ValidateMapPostRequest(r)

	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMap, err := controller.PostMap(
		requestValidated.Name,
		requestValidated.Version,
		requestValidated.Description,
		requestValidated.File,
	)

	if err != nil {
		logger.Errorf("Error creating map: %v", err)
		sendErrorResponse(w, "Error creating map", http.StatusInternalServerError)
	} else {
		serializeAndSendResponse(w, dto.MapBoToMapDto(newMap), http.StatusOK)
	}
}

func deleteMap(w http.ResponseWriter, r *http.Request) {
	controller := r.Context().Value(mapsContextKey).(*controllers.MapController)

	mapName := strings.ToLower(mux.Vars(r)["name"])

	err := controller.DeleteMap(mapName)

	if err != nil {
		logger.Errorf("Error deleting map: %v", err)
		sendErrorResponse(w, "Error deleting map", http.StatusInternalServerError)
	} else {
		sendOKResponse(w, "Map deleted")
	}
}
