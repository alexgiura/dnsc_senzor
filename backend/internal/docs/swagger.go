package docs

import (
	_ "embed"
	"net/http"

	"github.com/gorilla/mux"
)

//go:embed openapi.yaml
var openAPISpec []byte

//go:embed swagger.html
var swaggerHTML []byte

// Register mounts OpenAPI spec and Swagger UI at /openapi.yaml and /swagger.
func Register(router *mux.Router) {
	router.HandleFunc("/openapi.yaml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/yaml; charset=utf-8")
		_, _ = w.Write(openAPISpec)
	}).Methods(http.MethodGet)

	router.HandleFunc("/swagger", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(swaggerHTML)
	}).Methods(http.MethodGet)
}
