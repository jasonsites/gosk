package exampletest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jasonsites/gosk/internal/core/models"
	fx "github.com/jasonsites/gosk/test/fixtures"
	utils "github.com/jasonsites/gosk/test/testutils"
)

type CreateSetup struct {
	Name        string
	Description string
	Expected    utils.Expected
	Model       *models.ExampleDTO
}

func Test_Example_Create(t *testing.T) {
	s := Suite{}
	teardownSuite := s.SetupSuite(t)
	defer teardownSuite(t)

	tests := []CreateSetup{
		{
			Name:        "success",
			Description: "succeeds (201) with valid payload",
			Expected:    utils.Expected{Code: http.StatusCreated},
			Model:       fx.ExampleModel(nil),
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			t.Parallel()

			teardownTest := s.SetupTest(t)
			defer teardownTest(t)

			rd := &utils.RequestData{
				Body:   fx.ComposeJSONBody(fx.ExampleRequest(tc.Model)),
				Method: http.MethodPost,
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
