package routes

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	"github.com/jasonsites/gosk/internal/core/models"
	ctrl "github.com/jasonsites/gosk/internal/http/controllers"
)

// ExampleRouter implements a router group for an Example resource
func ExampleRouter(r *chi.Mux, ns string, c *ctrl.Controller) {
	prefix := fmt.Sprintf("/%s/examples", ns)

	// resource provides a RequestBody with data binding for the Example model
	// for use with Create/Update Controller methods
	resource := func() *ctrl.RequestBody {
		return &ctrl.RequestBody{
			Data: &ctrl.RequestResource{
				Attributes: &models.ExampleInputData{},
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
