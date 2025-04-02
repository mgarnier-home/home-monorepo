package routers

import (
	"net/http"
	"strings"

	"mgarnier11.fr/go/mineager/server/controllers"
	"mgarnier11.fr/go/mineager/server/objects/dto"
	"mgarnier11.fr/go/mineager/server/routers/validation"

	"mgarnier11.fr/go/libs/logger"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/mux"
)

type MapsRouter struct {
	utils          RouterUtils
	mapsController *controllers.MapsController
}

func NewMapsRouter(router *mux.Router, serverLogger *logger.Logger) *MapsRouter {
	mapsRouter := &MapsRouter{
		utils: RouterUtils{
			logger: logger.NewLogger(
				"[MAPS]",
				"%-10s ",
				lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
				serverLogger,
			),
		},
		mapsController: controllers.NewMapsController(),
	}

	mapRouter := router.PathPrefix("/maps").Subrouter()
	mapRouter.HandleFunc("", mapsRouter.getMaps).Methods("GET")
	mapRouter.HandleFunc("/{name}", mapsRouter.getMap).Methods("GET")
	mapRouter.HandleFunc("", mapsRouter.postMap).Methods("POST")
	mapRouter.HandleFunc("/{name}", mapsRouter.deleteMap).Methods("DELETE")

	return mapsRouter
}

func (router *MapsRouter) getMaps(w http.ResponseWriter, r *http.Request) {
	maps, err := router.mapsController.GetMaps()

	if err != nil {
		router.utils.logger.Errorf("Error getting maps: %v", err)
		router.utils.sendErrorResponse(w, r, "Error getting maps", http.StatusInternalServerError)
	} else {
		router.utils.serializeAndSendResponse(w, r, dto.MapMapsBoToMapsDto(maps), http.StatusOK)
	}
}

func (router *MapsRouter) getMap(w http.ResponseWriter, r *http.Request) {
	mapName := strings.ToLower(mux.Vars(r)["name"])

	mapBo, err := router.mapsController.GetMap(mapName)

	if err != nil {
		router.utils.logger.Errorf("Error getting map: %v", err)
		router.utils.sendErrorResponse(w, r, "Error getting map", http.StatusInternalServerError)
	} else {
		router.utils.serializeAndSendResponse(w, r, dto.MapMapBoToMapDto(mapBo), http.StatusOK)
	}
}

func (router *MapsRouter) postMap(w http.ResponseWriter, r *http.Request) {
	requestValidated, err := validation.ValidateMapPostRequest(r)

	if err != nil {
		router.utils.sendErrorResponse(w, r, err.Error(), http.StatusBadRequest)
		return
	}

	newMap, err := router.mapsController.PostMap(
		requestValidated.Name,
		requestValidated.Version,
		requestValidated.Description,
		requestValidated.File,
	)

	if err != nil {
		router.utils.logger.Errorf("Error creating map: %v", err)
		router.utils.sendErrorResponse(w, r, "Error creating map", http.StatusInternalServerError)
	} else {
		router.utils.serializeAndSendResponse(w, r, dto.MapMapBoToMapDto(newMap), http.StatusOK)
	}
}

func (router *MapsRouter) deleteMap(w http.ResponseWriter, r *http.Request) {
	mapName := strings.ToLower(mux.Vars(r)["name"])

	err := router.mapsController.DeleteMap(mapName)

	if err != nil {
		router.utils.logger.Errorf("Error deleting map: %v", err)
		router.utils.sendErrorResponse(w, r, "Error deleting map", http.StatusInternalServerError)
	} else {
		router.utils.sendOKResponse(w, r, "Map deleted")
	}
}
