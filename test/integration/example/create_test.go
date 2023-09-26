package exampletest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jasonsites/gosk/internal/core/models"
	"github.com/jasonsites/gosk/internal/resolver"
	fx "github.com/jasonsites/gosk/test/fixtures"
	utils "github.com/jasonsites/gosk/test/testutils"
)

type CreateSetup struct {
	Name        string
	Description string
	Expected    utils.Expected
	Model       *models.ExampleInputData
	Resolver    *resolver.Resolver
}

func setupCreateSuite(tb testing.TB) func(tb testing.TB) {
	// setup for test table

	return func(tb testing.TB) {
		// teardown for test table
	}
}

func setupCreateTest(tb testing.TB, r *resolver.Resolver) func(tb testing.TB) {
	// setup for each test

	return func(tb testing.TB) {
		utils.Cleanup(r)
	}
}

func Test_Example_Create(t *testing.T) {
	teardownCreateSuite := setupCreateSuite(t)
	defer teardownCreateSuite(t)

	conf := &resolver.Config{}
	resolver, err := utils.InitializeResolver(conf, "")
	if err != nil {
		t.Fatalf("app initialization error: %+v\n", err)
	}

	handler := resolver.HTTPServer().Server.Handler

	tests := []CreateSetup{
		{
			Name:        "success",
			Description: "succeeds (201) with valid payload",
			Expected:    utils.Expected{Code: http.StatusCreated},
			Model:       fx.ExampleModel(nil),
			Resolver:    resolver,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			teardownCreateTest := setupCreateTest(t, resolver)
			defer teardownCreateTest(t)

			rd := &utils.RequestData{
				Body:   fx.ComposeJSONBody(fx.ExampleRequest(tc.Model)),
				Method: http.MethodPost,
				Route:  routePrefix,
			}

			req, err := rd.SetRequestData(nil)
			if err != nil {
				t.Fatalf("http request error: %+v\n", err)
			}

			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			res := rec.Result()
			if res.StatusCode != tc.Expected.Code {
				t.Errorf("expected '%d', actual '%d'", tc.Expected.Code, res.StatusCode)
			}
		})

	}
}
