package user

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dozto/dozto-aster-auth/pkg/test"
	"github.com/gofiber/fiber/v3"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

var app *fiber.App

func TestMain(m *testing.M) {
	var userModel *UserModel

	{ // Init test environment
		db := test.Init()
		app = fiber.New()

		userModel = NewUserModel(db, "users")
		InitUserRoutes(app, userModel)
	}

	m.Run()

	userModel.users.Drop(context.Background())
}

func TestUser(t *testing.T) {

	tests := []struct {
		description  string // description of the test case
		method       string // HTTP method to test
		route        string // route path to test
		body         string
		expectedCode int // expected HTTP status code
	}{
		{description: "Create User",
			method:       "POST",
			route:        "/users",
			body:         `{"email":"ole@gmail.com","phone":"+8613302123021","Password":"123321","Meta":{"emailVerify":"fds"}}`,
			expectedCode: 201},
	}

	for _, test := range tests {
		// Create a new http request with the route from the test case
		req := httptest.NewRequest(test.method, test.route, strings.NewReader(test.body))
		req.Header.Set("Content-type", "application/json")

		// Perform the request plain with the app,
		// the second argument is a request latency
		// (set to -1 for no latency)
		resp, _ := app.Test(req, -1)

		log.Info().Any("res", resp)

		// Verify, if the status code is as expected
		assert.Equalf(t, test.expectedCode, resp.StatusCode, test.description)
	}
}
