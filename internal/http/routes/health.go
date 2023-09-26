package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jasonsites/gosk/internal/core/jsonapi"
	ctrl "github.com/jasonsites/gosk/internal/http/controllers"
)

// HealthRouter implements a router for healthcheck
func HealthRouter(r *chi.Mux, ns string) {
	prefix := fmt.Sprintf("/%s/health", ns)

	status := func(w http.ResponseWriter, r *http.Request) {
		data := jsonapi.Envelope{"meta": jsonapi.Envelope{"status": "healthy"}}
		ctrl.EncodeResponse(w, r, http.StatusOK, data)
	}

	r.Get(prefix, status)
}
