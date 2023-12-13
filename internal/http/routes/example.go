package routes

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/jasonsites/gosk/internal/core/interfaces"
)

// ExampleRouter implements a router group for an Example resource
func ExampleRouter(r *chi.Mux, ns string, c interfaces.Controller) {
	prefix := fmt.Sprintf("/%s/examples", ns)

	r.Route(prefix, func(r chi.Router) {
		r.Get("/", c.List())
		r.Get("/{id}", c.Detail())
		r.Post("/", c.Create())
		r.Put("/{id}", c.Update())
		r.Delete("/{id}", c.Delete())
	})
}
