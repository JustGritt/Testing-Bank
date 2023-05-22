package routes_test

import (
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

	err = database.Database.Db.Migrator().DropTable("accounts")
	if err != nil {
		t.Fatalf("failed to delete test database: %v", err)
	}

	err = database.Database.Db.Migrator().DropTable("users")
	if err != nil {
		t.Fatalf("failed to delete test database: %v", err)
	}

}

func TestCreateMaxAccountHandler(t *testing.T) {

	// Create a new Fiber app
	app := fiber.New()

	database.ConnectTest()

	// Define the test route
	app.Post("/api/users/:id/account", routes.CreateAccount)

	//create 6 user
	for i := 0; i < 4; i++ {
		req := httptest.NewRequest(http.MethodPost, "/api/users/1/account", nil)
		// Create a response recorder to capture the response
		resp, err := app.Test(req, -1)
		if err != nil {
			t.Fatalf("failed to perform request: %v", err)
		}
	}

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}
