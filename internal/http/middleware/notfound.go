package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jasonsites/gosk/internal/cerror"
	"github.com/jasonsites/gosk/internal/http/jsonio"
)

// NotFound
func NotFound(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		tctx := chi.NewRouteContext()
		if !rctx.Routes.Match(tctx, r.Method, r.URL.Path) {
			err := cerror.NewNotFoundError(nil, "path not found")
			jsonio.EncodeError(w, r, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}
