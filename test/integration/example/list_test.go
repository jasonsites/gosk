package exampletest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	utils "github.com/jasonsites/gosk/test/testutils"
)

type ListSetup struct {
	Name        string
	Description string
	Expected    utils.Expected
}

func Test_Example_List(t *testing.T) {
	s := Suite{}
	teardownSuite := s.SetupSuite(t)
	defer teardownSuite(t)

	tests := []ListSetup{
		{
			Name:        "success",
			Description: "succeeds (200)",
			Expected:    utils.Expected{Code: http.StatusOK},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			teardownTest := s.SetupTest(t)
			defer teardownTest(t)

			// db := resolver.PostgreSQLClient()
			// record, err := insertRecord(db)
			// if err != nil {
			// 	t.Fatalf("db insert error: %+v\n", err)
			// }

			rd := &utils.RequestData{
				Method: http.MethodGet,
				Route:  s.RoutePrefix,
			}

			req, err := rd.SetRequestData(nil)
			if err != nil {
				t.Fatalf("http request error: %+v\n", err)
			}

			rec := httptest.NewRecorder()
			s.Handler.ServeHTTP(rec, req)

			res := rec.Result()
			if res.StatusCode != tc.Expected.Code {
				t.Errorf("expected '%d', actual '%d'", tc.Expected.Code, res.StatusCode)
			}
		})
	}
}
