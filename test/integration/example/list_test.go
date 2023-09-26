package exampletest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jasonsites/gosk/internal/resolver"
	utils "github.com/jasonsites/gosk/test/testutils"
)

type ListSetup struct {
	Name        string
	Description string
	Expected    utils.Expected
	Resolver    *resolver.Resolver
}

func setupListSuite(tb testing.TB) func(tb testing.TB) {
	// setup for test table

	return func(tb testing.TB) {
		// teardown for test table
	}
}

func setupListTest(tb testing.TB, r *resolver.Resolver) func(tb testing.TB) {
	// setup for each test

	return func(tb testing.TB) {
		utils.Cleanup(r)
	}
}

func Test_Example_List(t *testing.T) {
	teardownListSuite := setupListSuite(t)
	defer teardownListSuite(t)

	conf := &resolver.Config{}
	resolver, err := utils.InitializeResolver(conf, "")
	if err != nil {
		t.Fatalf("app initialization error: %+v\n", err)
	}

	handler := resolver.HTTPServer().Server.Handler

	tests := []ListSetup{
		{
			Name:        "success",
			Description: "succeeds (200)",
			Expected:    utils.Expected{Code: http.StatusOK},
			Resolver:    resolver,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			teardownListTest := setupListTest(t, resolver)
			defer teardownListTest(t)

			// db := resolver.PostgreSQLClient()
			// record, err := insertRecord(db)
			// if err != nil {
			// 	t.Fatalf("db insert error: %+v\n", err)
			// }

			rd := &utils.RequestData{
				Method: http.MethodGet,
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
