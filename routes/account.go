package routes

import (
	"github.com/JustGritt/go-grpc/database"
	"github.com/JustGritt/go-grpc/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

type Account struct {
	ID      uint   `json:"id"`
	IBAN    string `json:"iban"`
	Balance uint   `json:"balance"`
}

func CreateResponseAccount(account models.Account) Account {
	return Account{
		ID:      account.ID,
		IBAN:    account.IBAN,
		Balance: account.Balance,
	}
}

// Create account godoc
// @Summary Create account
// @Description Create account
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param account body Account true "Account"
// @Success 201 {object} Account
// @Failure 400 {string} string "Invalid request"
// @Failure 500 {string} string "Error encrypting password"
// @Router /api/accounts [post]
func CreateAccount(c *fiber.Ctx) error {
	var account models.Account

	// Get user from JWT
	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	userID := claims["id"].(float64)
	account.UserID = uint(userID)

	if err := c.BodyParser(&account); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if account.UserID == "" {
		return c.Status(400).JSON("Account name is required")
	}

	var existingAccount models.Account
	database.Database.Db.Where("name = ?", account.IBAN).First(&existingAccount)
	if existingAccount.ID != 0 {
		return c.Status(400).JSON("Account name already taken")
	}

	database.Database.Db.Create(&account)
	responseAccount := CreateResponseAccount(account)
	return c.Status(200).JSON(responseAccount)
}

func GetAccountId(id int, account *models.Account) error {
	database.Database.Db.Find(&account, id)
	if account.ID == 0 {
		return fiber.NewError(404, "Account not found")
	}

	return nil
}

// Get accounts godoc
// @Summary Get accounts
// @Description Get accounts
// @Tags accounts
// @Accept  json
// @Produce  json
// @Success 200 {object} Account
// @Failure 404 {string} string "Account not found"
// @Router /api/accounts [get]
func GetAccounts(c *fiber.Ctx) error {
	var accounts []models.Account
	database.Database.Db.Find(&accounts)

	var responseAccounts []Account
	for _, account := range accounts {
		responseAccounts = append(responseAccounts, CreateResponseAccount(account))
	}

	return c.Status(200).JSON(responseAccounts)
}

// Get account godoc
// @Summary Get account
// @Description Get account
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} Account
// @Failure 404 {string} string "Account not found"
// @Router /api/accounts/{id} [get]
func GetAccount(c *fiber.Ctx) error {
	id := c.Params("id")
	var account models.Account
	database.Database.Db.Find(&account, id)

	if account.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Account not found",
		})
	}

	return c.Status(200).JSON(CreateResponseAccount(account))
}

// Update account godoc
// @Summary Update account
// @Description Update account
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Param account body Account true "Account"
// @Success 200 {object} Account
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Account not found"
// @Failure 500 {string} string "Error encrypting password"
// @Router /api/accounts/{id} [put]
func UpdateAccount(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var account models.Account

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = GetAccountId(id, &account)

	token := c.Locals("account").(*jwt.Token)

	if !validToken(token, id) {
		return c.Status(500).JSON(fiber.Map{"status": "error", "message": "Invalid token id", "data": nil})
	}

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	type UpdateAccount struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Password  string `json:"password"`
	}

	var updateData UpdateAccount

	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	// encrypt password with bcrypt
	pass, err := bcrypt.GenerateFromPassword([]byte(updateData.Password), 14)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"message": "Error encrypting password",
		})
	}

	account.FirstName = updateData.FirstName
	account.LastName = updateData.LastName
	account.Password = string(pass)

	database.Database.Db.Save(&account)

	responseAccount := CreateResponseAccount(account)

	return c.Status(200).JSON(responseAccount)
}

// Delete account godoc
// @Summary Delete account
// @Description Delete account
// @Tags accounts
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {string} string "Account deleted"
// @Failure 400 {string} string "Invalid request"
// @Failure 404 {string} string "Account not found"
// @Router /api/accounts/{id} [delete]
func DeleteAccount(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var account models.Account

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	err = GetAccountId(id, &account)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Delete(&account)

	return c.Status(200).JSON(fiber.Map{
		"message": "Account deleted",
	})
}
