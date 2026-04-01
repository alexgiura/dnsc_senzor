package routes

import (
	"senzor/internal/handlers"
	"senzor/internal/services"

	"github.com/gorilla/mux"
)

func RegisterRoutes(app *services.AppServices) *mux.Router {
	router := mux.NewRouter()

	handlers.RegisterSystemRoutes(router)
	if app != nil {
		handlers.RegisterNetworkAlertRoutes(router, app.NetworkAlert)
	}

	return router
}
