package routes

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	ctrl "github.com/jasonsites/gosk-api/internal/httpapi/controllers"
	"github.com/jasonsites/gosk-api/internal/types"
)

// BaseRouter only exists to easily verify a working app and should normally be removed
func BaseRouter(r *chi.Mux, c *ctrl.Controller, ns string) {
	prefix := fmt.Sprintf("/%s", ns)

	get := func(w http.ResponseWriter, r *http.Request) {
		headers := r.Header
		host := r.Host
		path := r.URL.Path
		remoteAddress := r.RemoteAddr

		c.WriteJSON(w, http.StatusOK, types.Map{
			"data": "base router is working...",
			"request": types.Map{
				"headers":       headers,
				"host":          host,
				"path":          path,
				"remoteAddress": remoteAddress,
			},
		})
	}

	r.Route(prefix, func(r chi.Router) {
		r.Get("/", get)
	})
}
