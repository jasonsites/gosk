package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	ctrl "github.com/jasonsites/gosk-api/internal/httpapi/controllers"
	"github.com/jasonsites/gosk-api/internal/types"
)

// HealthRouter implements an router group for healthcheck
func HealthRouter(r *chi.Mux, c *ctrl.Controller, ns string) {
	prefix := fmt.Sprintf("/%s/health", ns)

	status := func(w http.ResponseWriter, r *http.Request) {
		c.WriteJSON(w, http.StatusOK, types.Map{"meta": types.Map{"status": "healthy"}})
	}

	r.Get(prefix, status)
}
