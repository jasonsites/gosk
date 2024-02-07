package routes

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/jasonsites/gosk/internal/core/interfaces"
	"github.com/jasonsites/gosk/internal/core/jsonapi"
	"github.com/jasonsites/gosk/internal/core/models"
)

// ExampleRouter implements a router group for an Example resource
func ExampleRouter(r *chi.Mux, ns string, c interfaces.ExampleController) {
	prefix := fmt.Sprintf("/%s/examples", ns)

	// resource provides a RequestBody with data binding for the Example model
	// for use with Create/Update Controller methods
	resource := func() *jsonapi.RequestBody {
		return &jsonapi.RequestBody{
			Data: &jsonapi.RequestResource{
				Attributes: &models.ExampleDTO{},
			},
		}
	}

	r.Route(prefix, func(r chi.Router) {
		r.Get("/", c.List())
		r.Get("/{id}", c.Detail())
		r.Post("/", c.Create(resource))
		r.Put("/{id}", c.Update(resource))
		r.Delete("/{id}", c.Delete())
	})
}
