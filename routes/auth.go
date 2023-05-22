package routes

import (
	"time"

	"github.com/JustGritt/go-grpc/database"
	"github.com/JustGritt/go-grpc/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func Login(c *fiber.Ctx) error {
	var user models.User

	if err := c.BodyParser(&user); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if user.Email == "" || user.Password == "" {
		return c.Status(400).JSON("Please provide an email and password")
	}

	var existingUser models.User
	database.Database.Db.Where("email = ?", user.Email).First(&existingUser)
	if existingUser.ID == 0 {
		return c.Status(400).JSON("Invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(user.Password)); err != nil {
		return c.Status(400).JSON("Invalid credentials")
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"name":    existingUser.FirstName,
		"user_id": existingUser.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})

}
