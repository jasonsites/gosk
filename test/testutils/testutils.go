package testutils

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/jasonsites/gosk-api/internal/resolver"
)

// Cleanup deletes all rows on all database tables
func Cleanup(r *resolver.Resolver) error {
	db, err := r.PostgreSQLClient()
	if err != nil {
		return err
	}

	tables := []string{"entity"}

	for _, t := range tables {
		sql := fmt.Sprintf("DELETE from %s", t)
		_, err = db.Exec(context.TODO(), sql)
		if err != nil {
			return err
		}
	}

	return nil
}

// InitializeApp creates a new Resolver from the given config and returns a reference to the HTTP Server's App (Fiber) instance and the Resolver itself
func InitializeApp(conf *resolver.Config) (*fiber.App, *resolver.Resolver, error) {
	resolver := resolver.NewResolver(context.Background(), conf)
	if err := resolver.Initialize(); err != nil {
		return nil, nil, err
	}

	server, err := resolver.HTTPServer()
	if err != nil {
		return nil, resolver, err
	}

	return server.App, resolver, nil
}

// SetRequestData creates a new HTTP Request instance from the give data
func SetRequestData(method, route string, body io.Reader, headers map[string]string) *http.Request {
	req := httptest.NewRequest(method, route, body)
	if headers != nil {
		req = SetRequestHeaders(req, headers)
	}
	return req
}

// SetRequestHeaders set all headers on the given request
func SetRequestHeaders(req *http.Request, headers map[string]string) *http.Request {
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	return req
}
