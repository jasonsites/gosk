package exampletest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jasonsites/gosk/internal/core/models"
	"github.com/jasonsites/gosk/internal/resolver"
	fx "github.com/jasonsites/gosk/test/fixtures"
	utils "github.com/jasonsites/gosk/test/testutils"
)

type UpdateSetup struct {
	Name        string
	Description string
	Expected    utils.Expected
	Model       *models.ExampleInputData
	Request     *utils.RequestData
	Resolver    *resolver.Resolver
}

func setupUpdateSuite(tb testing.TB) func(tb testing.TB) {
	// setup for test table

	return func(tb testing.TB) {
		// teardown for test table
	}
}

func setupUpdateTest(tb testing.TB, r *resolver.Resolver) func(tb testing.TB) {
	// setup for each test

	return func(tb testing.TB) {
		utils.Cleanup(r)
	}
}

func Test_Example_Update(t *testing.T) {
	teardownUpdateSuite := setupUpdateSuite(t)
	defer teardownUpdateSuite(t)

	conf := &resolver.Config{}
	resolver, err := utils.InitializeResolver(conf, "")
	if err != nil {
		t.Fatalf("app initialization error: %+v\n", err)
	}

	handler := resolver.HTTPServer().Server.Handler

	tests := []UpdateSetup{
		{
			Name:        "success",
			Description: "succeeds (200) with valid payload",
			Expected:    utils.Expected{Code: http.StatusOK},
			Model:       fx.ExampleModel(nil),
			Request:     &utils.RequestData{},
			Resolver:    resolver,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			teardownUpdateTest := setupUpdateTest(t, resolver)
			defer teardownUpdateTest(t)

			db := resolver.PostgreSQLClient()
			record, err := insertExampleRecord(db)
			if err != nil {
				t.Fatalf("db insert error: %+v\n", err)
			}

			rd := &utils.RequestData{
				Body:   fx.ComposeJSONBody(fx.ExampleRequest(tc.Model)),
				Method: http.MethodPut,
				Route:  fmt.Sprintf("%s/%s", routePrefix, record.ID.String()),
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
