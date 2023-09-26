package exampletest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	fx "github.com/jasonsites/gosk/test/fixtures"
	utils "github.com/jasonsites/gosk/test/testutils"
)

type DeleteSetup struct {
	Name        string
	Description string
	Expected    utils.Expected
}

func Test_Example_Delete(t *testing.T) {
	s := Suite{}
	teardownSuite := s.SetupSuite(t)
	defer teardownSuite(t)

	tests := []DeleteSetup{
		{
			Name:        "success",
			Description: "succeeds (204)",
			Expected:    utils.Expected{Code: http.StatusNoContent},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			teardownTest := s.SetupTest(t)
			defer teardownTest(t)

			entity := fx.ExampleEntityRecord(nil, nil)
			record, err := insertExampleRecord(entity, s.DB)
			if err != nil {
				t.Fatalf("db insert error: %+v\n", err)
			}

			rd := &utils.RequestData{
				Method: http.MethodDelete,
				Route:  fmt.Sprintf("%s/%s", s.RoutePrefix, record.ID.String()),
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
