package routers

import (
	"context"
	"mgarnier11/go/logger"
	"mgarnier11/mineager/server/controllers"
	"mgarnier11/mineager/server/objects/dto"
	"net/http"

	"github.com/charmbracelet/lipgloss"
	"github.com/gorilla/mux"
)

const hostsContextKey contextKey = "hosts"

type HostsRouter struct {
	contextKey             contextKey
	logger                 *logger.Logger
	hostsController        *controllers.HostsController
	hostsSubrouter         *mux.Router
	hostsHostnameSubrouter *mux.Router
}

func NewHostsRouter(router *mux.Router, serverLogger *logger.Logger) *HostsRouter {
	return &HostsRouter{
		contextKey: hostsContextKey,
		logger: logger.NewLogger(
			"[HOSTS]",
			"%-10s ",
			lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")),
			serverLogger,
		),
		hostsController: controllers.NewHostsController(),
	}
}

func getHostsControllerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		controller := controllers.NewHostsController()

		ctx := context.WithValue(r.Context(), hostsContextKey, controller)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getHostsControllerFromContext(ctx context.Context) *controllers.HostsController {
	return ctx.Value(hostsContextKey).(*controllers.HostsController)
}

func HostsRoutes(router *mux.Router) *mux.Router {
	hostRouter := router.PathPrefix("/hosts").Subrouter()
	hostRouter.Use(getHostsControllerMiddleware)

	hostRouter.HandleFunc("", getHosts).Methods("GET")
	hostNameRouter := hostRouter.PathPrefix("/{hostName}").Subrouter()

	hostNameRouter.HandleFunc("", getHost).Methods("GET")

	return hostNameRouter
}

func getHosts(w http.ResponseWriter, r *http.Request) {
	controller := getHostsControllerFromContext(r.Context())

	hosts := controller.GetHosts()

	serializeAndSendResponse(w, dto.MapHostsBoHostsDto(hosts), http.StatusOK)
}

func getHost(w http.ResponseWriter, r *http.Request) {
	controller := getHostsControllerFromContext(r.Context())

	hostName := mux.Vars(r)["hostName"]

	host, err := controller.GetHost(hostName)

	if err != nil {
		sendErrorResponse(w, "Error getting host", http.StatusInternalServerError)
	} else {
		serializeAndSendResponse(w, dto.MapHostBoToHostDto(host), http.StatusOK)
	}
}
