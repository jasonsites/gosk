package routes

// import (
// 	"net/http/httptest"
// 	"testing"

// 	"github.com/gofiber/fiber/v2"
// 	"github.com/stretchr/testify/assert" // add Testify package
// )

// func TestEpisodeCreate(t *testing.T) {
// 	tests := []struct {
// 		description  string
// 		route        string
// 		expectedCode int
// 	}{
// 		{
// 			description:  "Create Episode",
// 			route:        "/",
// 			expectedCode: 201,
// 		},
// 	}

// 	app := fiber.New()

// 	// Iterate through test single test cases
// 	for _, test := range tests {
// 		// Create a new http request with the route from the test case
// 		req := httptest.NewRequest("POST", test.route, nil)

// 		// Perform the request plain with the app,
// 		// the second argument is a request latency
// 		// (set to -1 for no latency)
// 		resp, _ := app.Test(req, 1)

// 		// Verify, if the status code is as expected
// 		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
// 	}
// }
