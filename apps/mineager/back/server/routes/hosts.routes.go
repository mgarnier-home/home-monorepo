package routes

import (
	"context"
	"mgarnier11/mineager/server/controllers"
	"mgarnier11/mineager/server/objects/dto"
	"net/http"

	"github.com/gorilla/mux"
)

const hostsContextKey contextKey = "hosts"

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

func HostsRoutes(router *mux.Router) {
	hostRouter := router.PathPrefix("/hosts").Subrouter()
	hostRouter.Use(getHostsControllerMiddleware)

	hostRouter.HandleFunc("", getHosts).Methods("GET")
}

func getHosts(w http.ResponseWriter, r *http.Request) {
	controller := getHostsControllerFromContext(r.Context())

	hosts := controller.GetHosts()

	serializeAndSendResponse(w, dto.MapHostsBoHostsDto(hosts), http.StatusOK)
}
