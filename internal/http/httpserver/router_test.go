package httpserver

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jasonsites/gosk/test/mock"
)

func Test_RegisterRoutes(t *testing.T) {
	ns := "test"
	var registered = []struct {
		route  string
		method string
	}{
		{fmt.Sprintf("/%s/", ns), http.MethodGet},
		{fmt.Sprintf("/%s/health", ns), http.MethodGet},
		{fmt.Sprintf("/%s/examples/", ns), http.MethodGet},
		{fmt.Sprintf("/%s/examples/{id}", ns), http.MethodGet},
		{fmt.Sprintf("/%s/examples/", ns), http.MethodPost},
		{fmt.Sprintf("/%s/examples/{id}", ns), http.MethodPut},
		{fmt.Sprintf("/%s/examples/{id}", ns), http.MethodDelete},
	}

	routerConf := &RouterConfig{Namespace: ns}
	controllers := &controllerRegistry{
		ExampleController: &mock.Controller{},
	}

	router := mock.Mux()
	registerRoutes(routerConf, router.(*chi.Mux), controllers)
	routes := router.(chi.Routes)

	for _, r := range registered {
		if !mock.RouteExists(r.route, r.method, routes) {
			t.Errorf("unregistered route '%s'", r.route)
		}
	}
}
