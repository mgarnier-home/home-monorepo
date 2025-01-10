package routers

import (
	"mgarnier11/go/logger"
	"mgarnier11/mineager/server/controllers"
	"mgarnier11/mineager/server/objects/dto"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/mux"
)

type HostsRouter struct {
	utils           RouterUtils
	hostsController *controllers.HostsController
}

func NewHostsRouter(router *mux.Router, serverLogger *logger.Logger) *HostsRouter {
	hostsRouter := &HostsRouter{
		utils: RouterUtils{
			logger: logger.NewLogger(
				"[HOSTS]",
				"%-10s ",
				lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
				serverLogger,
			),
		},
		hostsController: controllers.NewHostsController(),
	}
	router.HandleFunc("/hosts", hostsRouter.getHosts).Methods("GET")

	router.HandleFunc("/hosts/{hostName}", hostsRouter.getHost).Methods("GET")

	return hostsRouter
}

func (router *HostsRouter) getHosts(w http.ResponseWriter, r *http.Request) {
	hosts := router.hostsController.GetHosts()

	router.utils.serializeAndSendResponse(w, r, dto.MapHostsBoHostsDto(hosts), http.StatusOK)
}

func (router *HostsRouter) getHost(w http.ResponseWriter, r *http.Request) {
	hostName := mux.Vars(r)["hostName"]

	host, err := router.hostsController.GetHost(hostName)

	if err != nil {
		router.utils.logger.Errorf("Error getting host: %v", err)
		router.utils.sendErrorResponse(w, r, "Error getting host", http.StatusInternalServerError)
	} else {
		router.utils.serializeAndSendResponse(w, r, dto.MapHostBoToHostDto(host), http.StatusOK)
	}
}
