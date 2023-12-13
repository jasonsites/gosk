package mock

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Controller struct{}

func (c *Controller) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
func (c *Controller) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
func (c *Controller) Detail() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
func (c *Controller) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}
func (c *Controller) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {}
}

func Mux() http.Handler {
	return chi.NewRouter()
}

func RouteExists(testRoute, testMethod string, chiRoutes chi.Routes) bool {
	exists := false

	_ = chi.Walk(chiRoutes, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if strings.EqualFold(method, testMethod) && strings.EqualFold(route, testRoute) {
			exists = true
		}
		return nil
	})

	return exists
}
