package core

import (
	"net/http/httptest"
	"testing"

	"github.com/dozto/dozto-aster-auth/pkg/test"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var app *fiber.App

func TestMain(m *testing.M) {
	{ // Init test environment
		test.Init()
		app = fiber.New()

		InitCoreRoutes(app)
	}

	m.Run()
}

func TestHealthCheck(t *testing.T) {

	tests := []struct {
		description  string // description of the test case
		route        string // route path to test
		expectedCode int    // expected HTTP status code
	}{
		{description: "healt check should return 200",
			route:        "/health",
			expectedCode: 200},
	}

	for _, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest("GET", test.route, nil)

		log.Info().Any("req", req)

		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Test(req, -1)

		log.Info().Any("res", resp)

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
