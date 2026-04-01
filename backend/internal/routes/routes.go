package routes

import (
	"senzor/internal/docs"
	"senzor/internal/handlers"
	"senzor/internal/services"

	"github.com/gorilla/mux"
)

func RegisterRoutes(app *services.AppServices) *mux.Router {
	router := mux.NewRouter()

	docs.Register(router)
	handlers.RegisterSystemRoutes(router)
	if app != nil {
		handlers.RegisterNetworkAlertRoutes(router, app.NetworkAlert)
	}

	return router
}
