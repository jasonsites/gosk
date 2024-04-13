package health

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jasonsites/gosk/internal/http/jsonapi"
	"github.com/jasonsites/gosk/internal/http/jsonio"
)

// HealthRouter implements a router for healthcheck
func HealthRouter(r *chi.Mux, ns string) {
	prefix := fmt.Sprintf("/%s/health", ns)

	status := func(w http.ResponseWriter, r *http.Request) {
		data := jsonapi.Envelope{"meta": jsonapi.Envelope{"status": "healthy"}}
		jsonio.EncodeResponse(w, r, http.StatusOK, data)
	}

	r.Get(prefix, status)
}
