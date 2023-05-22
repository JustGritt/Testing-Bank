package routes_test

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/JustGritt/go-grpc/database"
	"github.com/JustGritt/go-grpc/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateAccountHandler(t *testing.T) {

	// Create a new Fiber app
	app := fiber.New()

	database.ConnectTest()

	// Define the test route
	app.Post("/api/users/:id/account", routes.CreateAccount)

	// Create a test request for the specific user's account
	req := httptest.NewRequest(http.MethodPost, "/api/users/1/account", nil)

	// Create a response recorder to capture the response
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	// Assert the response status code, body, or any other expected behavior
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

}

func TestCreateMaxAccountHandler(t *testing.T) {

	// Create a new Fiber app
	app := fiber.New()

	database.ConnectTest()

	// Define the test route
	app.Post("/api/users/:id/account", routes.CreateAccount)

	//create 6 user
	for i := 0; i < 6; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/users/1/account", nil)
		// Create a response recorder to capture the response
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("failed to perform request: %v", err)
		}
		if i == 5 {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}
			assert.Equal(t, "{\"message\":\"You can't have more than 5 accounts\"}", string(body))
		}

	}

}

func TestDebitHandler(t *testing.T) {

	// Create a new Fiber app
	app := fiber.New()

	database.ConnectTest()

	// Define the test route
	app.Post("/api/users/:id/account/debit", routes.Debit)

	//create a body io.Reader
	body := []byte(`{"iban":"NL1ING","debit":1000}`)
	bodyReader := bytes.NewReader(body)

	// Create a test request for the specific user's account with 1000 debit and iban NL1ING
	req := httptest.NewRequest(http.MethodPost, "/api/users/1/account/debit", bodyReader)

	// Create a response recorder to capture the response
	resp, err := app.Test(req, -1)
	if err != nil {
		t.Fatalf("failed to perform request: %v", err)
	}

	// Assert the response status code, body, or any other expected behavior
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

}
