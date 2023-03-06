package resourcetest

import (
	"bytes"
	"testing"

	utils "github.com/jasonsites/gosk-api/test/testutils"
	"github.com/stretchr/testify/assert"
)

func TestResourceCreate(t *testing.T) {
	var (
		routePrefix = "/domain/resources"
		method      = "POST"
	)

	tests := []utils.Setup{
		{
			Description: "create resource succeeds (201) with valid payload",
			Route:       routePrefix,
			Request: utils.Request{
				Body: bytes.NewBuffer([]byte(`{
					"data": {
						"type": "resource",
						"properties": {
							"title": "Resource Title",
							"description": "Resource Description",
							"enabled": true,
							"status": 1
						}
					}
				}`)),
				Headers: map[string]string{
					"Content-Type": "application/json",
				},
			},

			Expected: utils.Expected{Code: 201},
		},
	}

	app, _, err := utils.InitializeApp(nil)
	if err != nil {
		t.Log(err)
	}

	for _, test := range tests {
		req := utils.SetRequestData(method, test.Route, test.Request.Body, test.Request.Headers)

		res, err := app.Test(req, 1000)
		if err != nil {
			t.Log(err)
		}

		assert.Equalf(t, test.Expected.Code, res.StatusCode, test.Description)
	}
}
