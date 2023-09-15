package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	ctrl "github.com/jasonsites/gosk/internal/http/controllers"
)

// HealthRouter implements a router for healthcheck
func HealthRouter(r *chi.Mux, ns string) {
	prefix := fmt.Sprintf("/%s/health", ns)

	status := func(w http.ResponseWriter, r *http.Request) {
		data := ctrl.Envelope{"meta": ctrl.Envelope{"status": "healthy"}}
		ctrl.EncodeResponse(w, r, http.StatusOK, data)
	}

	r.Get(prefix, status)
}
