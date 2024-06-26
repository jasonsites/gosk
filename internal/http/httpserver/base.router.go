package httpserver

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jasonsites/gosk/internal/http/jsonapi"
	"github.com/jasonsites/gosk/internal/http/jsonio"
)

// BaseRouter only exists to easily verify a working app and should normally be removed
func BaseRouter(r *chi.Mux, ns string) {
	prefix := fmt.Sprintf("/%s", ns)

	get := func(w http.ResponseWriter, r *http.Request) {
		headers := r.Header
		host := r.Host
		path := r.URL.Path
		remoteAddress := r.RemoteAddr

		data := jsonapi.Envelope{
			"data": "base router is working...",
			"request": jsonapi.Envelope{
				"headers":       headers,
				"host":          host,
				"path":          path,
				"remoteAddress": remoteAddress,
			},
		}

		jsonio.EncodeResponse(w, r, http.StatusOK, data)
	}

	r.Route(prefix, func(r chi.Router) {
		r.Get("/", get)
	})
}
