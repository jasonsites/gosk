package example

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jasonsites/gosk/internal/http/jsonapi"
)

// ExampleController
type ExampleController interface {
	Create(func() *jsonapi.RequestBody) http.HandlerFunc
	Delete() http.HandlerFunc
	Detail() http.HandlerFunc
	List() http.HandlerFunc
	Update(func() *jsonapi.RequestBody) http.HandlerFunc
}

// ExampleRouter implements a router group for an Example resource
func ExampleRouter(r *chi.Mux, ns string, c ExampleController) {
	prefix := fmt.Sprintf("/%s/examples", ns)

	// resource provides a RequestBody with data binding for the Example model
	// for use with Create/Update Controller methods
	resource := func() *jsonapi.RequestBody {
		return &jsonapi.RequestBody{
			Data: &jsonapi.RequestResource{
				Attributes: &ExampleDTORequest{},
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
