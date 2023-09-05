package routes

import (
	"fmt"

	"github.com/go-chi/chi/v5"
	ctrl "github.com/jasonsites/gosk-api/internal/httpapi/controllers"
	"github.com/jasonsites/gosk-api/internal/types"
)

// ExampleRouter implements a router group for an Example resource
func ExampleRouter(r *chi.Mux, c *ctrl.Controller, ns string) {
	prefix := fmt.Sprintf("/%s/examples", ns)

	// createResource provides a JSONRequestBody with data binding for the Example model
	// for use with Create/Update Controller methods
	createResource := func() *types.JSONRequestBody {
		return &types.JSONRequestBody{
			Data: &types.RequestResource{
				Attributes: &types.ExampleRequestData{},
			},
		}
	}

	r.Route(prefix, func(r chi.Router) {
		// r.With(httpin.NewInput(ListUserReposInput{})).Get("/", c.List())
		r.Get("/", c.List())
		r.Get("/{id}", c.Detail())
		r.Post("/", c.Create(createResource))
		r.Put("/{id}", c.Update(createResource))
		r.Delete("/{id}", c.Delete())
	})
}
