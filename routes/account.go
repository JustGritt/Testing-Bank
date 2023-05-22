package routes

import (
	"strconv"

	"github.com/JustGritt/go-grpc/database"
	"github.com/JustGritt/go-grpc/models"
	"github.com/gofiber/fiber/v2"
)

type Account struct {
	ID             uint   `json:"id"`
	IBAN           string `json:"iban"`
	Balance        uint   `json:"balance"`
	NumberAccounts int    `json:"number_accounts"`
}

type Transaction struct {
	IBAN     string `json:"iban"`
	Balance  uint   `json:"balance"`
	rejected uint
}

func CreateResponseAccount(account models.Account, user models.User) Account {
	return Account{
		ID:             account.ID,
		IBAN:           account.IBAN,
		Balance:        account.Balance,
		NumberAccounts: user.NumberAccounts,
	}
}

func CreateAccount(c *fiber.Ctx) error {
	id := c.Params("id")
	var user models.User
	database.Database.Db.First(&user, id)

	var account models.Account
	account.UserID = uint(user.ID)
	account.Balance = 0
	account.IBAN = "NL" + strconv.Itoa(int(user.ID)) + "ING"

	user.NumberAccounts = int(user.NumberAccounts + 1)

	if user.NumberAccounts > 5 {
		return c.Status(400).JSON(fiber.Map{
			"message": "You can't have more than 5 accounts",
		})

	}

	database.Database.Db.Create(&account)

	database.Database.Db.Save(&user)

	return c.Status(201).JSON(CreateResponseAccount(account, user))

}

func GetAccountsByUser(c *fiber.Ctx) error {
	id, _ := strconv.Atoi(c.Params("id"))

	var user models.User
	database.Database.Db.Preload("Accounts").First(&user, id)
	var accounts []Account

	for _, account := range user.Accounts {
		accounts = append(accounts, CreateResponseAccount(account, user))
	}

	return c.JSON(accounts)

}

func DeleteAccount(c *fiber.Ctx) error {
	iban := c.Params("iban")

	var account models.Account
	database.Database.Db.Where("IBAN = ?", iban).Delete(&account)
	if account.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Account not found",
		})
	}

	return c.SendString("Account deleted")
}

func DeleteAllAccounts(c *fiber.Ctx) error {
	//delete all accounts for the user with id
	id := c.Params("id")

	var user models.User
	database.Database.Db.First(&user, id)

	var accounts []models.Account
	database.Database.Db.Where("user_id = ?", id).Delete(&accounts)

	user.NumberAccounts = 0
	database.Database.Db.Save(&user)

	return c.SendString("All accounts deleted")

}

func Debit(c *fiber.Ctx) error {

	type TransactionInfo struct {
		IBAN  string `json:"iban"`
		Debit int    `json:"debit"`
	}

	var transactionInfo TransactionInfo

	if err := c.BodyParser(&transactionInfo); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	debitInt := transactionInfo.Debit
	if debitInt < 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Debit can't be negative",
		})
	}

	var account models.Account
	database.Database.Db.Where("iban = ?", transactionInfo.IBAN).First(&account)

	if account.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Account not found",
		})
	}

	//check account balance and debit if the balance is 0 or less than the debit return an error
	if account.Balance < uint(debitInt) || account.Balance == 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Not enough balance",
		})
	}

	account.Balance = account.Balance - uint(debitInt)
	database.Database.Db.Save(&account)

	return c.Status(200).JSON(account)
}

func Credit(c *fiber.Ctx) error {
	type TransactionInfo struct {
		IBAN   string `json:"iban"`
		Credit int    `json:"debit"`
	}

	var transactionInfo TransactionInfo

	if err := c.BodyParser(&transactionInfo); err != nil {
		return c.Status(500).JSON(err.Error())
	}

	creditInt := transactionInfo.Credit

	if creditInt < 0 {
		return c.Status(400).JSON(fiber.Map{
			"message": "Credit can't be negative",
		})
	}

	var account models.Account
	database.Database.Db.Where("iban = ?", transactionInfo.IBAN).First(&account)

	if account.ID == 0 {
		return c.Status(404).JSON(fiber.Map{
			"message": "Account not found",
		})
	}

	account.Balance = account.Balance + uint(creditInt)

	var rejected = 0
	if account.Balance > 1000 {
		tmp := account.Balance
		rejected = int(account.Balance) - 1000
		account.Balance = tmp - uint(rejected)
	}

	database.Database.Db.Save(&account)

	responseAccount := Transaction{
		IBAN:     account.IBAN,
		Balance:  account.Balance,
		rejected: uint(rejected),
	}
	//return somethin like that Account {iban: "fadf", balance: 1000, rejected: 100}
	return c.Status(200).JSON(responseAccount)
}
