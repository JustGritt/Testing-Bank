package routes

import (
	"github.com/JustGritt/go-grpc/database"
	"github.com/JustGritt/go-grpc/models"
	"github.com/gofiber/fiber/v2"
)

type Payment struct {
	ID        uint `json:"id"`
	ProductId uint `json:"product_id"`
	Amount    uint `json:"amount"`
}

func CreateResponsePayment(payment models.Payment) Payment {
	return Payment{
		ID:        payment.ID,
		ProductId: payment.ProductID,
		Amount:    payment.Price,
	}
}

func CreateResponsePayments(payments []models.Payment) []Payment {
	var response []Payment
	for _, payment := range payments {
		response = append(response, CreateResponsePayment(payment))
	}
	return response
}

// Get all payments godoc
// @Summary Get all payments
// @Description Get all payments
// @Tags payments
// @Accept  json
// @Produce  json
// @Success 200 {object} []Payment
// @Failure 400 {string} string "Payment not found"
// @Router /api/payments [get]
func GetPayments(c *fiber.Ctx) error {
	var payments []models.Payment
	database.Database.Db.Find(&payments)

	if len(payments) == 0 {
		return c.Status(404).JSON("No payments found")
	}

	return c.Status(200).JSON(CreateResponsePayments(payments))
}

// Get payment by id godoc
// @Summary Get payment by id
// @Description Get payment by id
// @Tags payments
// @Accept  json
// @Produce  json
// @Param id path int true "Payment ID"
// @Success 200 {object} Payment
// @Failure 400 {string} string "Payment not found"
// @Router /api/payments/{id} [get]
func GetPayment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var payment models.Payment
	database.Database.Db.Where("id = ?", id).Find(&payment)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if payment.ID == 0 {
		return c.Status(404).JSON("Payment not found")
	}

	return c.Status(200).JSON(CreateResponsePayment(payment))
}

// Create a new payment godoc
// @Summary Create a new payment
// @Description Create a new payment
// @Tags payments
// @Accept  json
// @Produce  json
// @Param payment body Payment true "Payment"
// @Success 200 {object} Payment
// @Failure 400 {string} string "Something went wrong"
// @Router /api/payments [post]
func CreatePayment(c *fiber.Ctx) error {
	var payment models.Payment

	if err := c.BodyParser(&payment); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	// Check if product exists
	var product models.Product
	database.Database.Db.Where("id = ?", payment.ProductID).Find(&product)

	if product.ID == 0 {
		return c.Status(404).JSON("Product not found")
	}

	// Check if the price of the product is the same as the price of the payment
	if product.Price != payment.Price {
		// add the correct price to the payment
		payment.Price = product.Price
		//return c.Status(400).JSON("Price of the product is not the same as the price of the payment")
	}

	database.Database.Db.Create(&payment)
	b.Publish(payment)
	return c.Status(200).JSON(CreateResponsePayment(payment))
}

// Delete the payment with the given id godoc
// @Summary Delete the payment with the given id
// @Description Delete the payment with the given id
// @Tags payments
// @Accept  json
// @Produce  json
// @Param id path int true "Payment ID"
// @Success 200 {object} Payment
// @Failure 400 {string} string "Payment not found"
// @Router /api/payments/{id} [delete]
func DeletePayment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var payment models.Payment
	database.Database.Db.Where("id = ?", id).Find(&payment)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if payment.ID == 0 {
		return c.Status(404).JSON("Payment not found")
	}

	database.Database.Db.Delete(&payment)
	return c.Status(200).JSON(CreateResponsePayment(payment))
}

// Update the payment with the given id godoc
// @Summary Update the payment with the given id
// @Description Update the payment with the given id
// @Tags payments
// @Accept  json
// @Produce  json
// @Param id path int true "Payment ID"
// @Param payment body Payment true "Payment"
// @Success 200 {object} Payment
// @Failure 400 {string} string "Payment not found"
// @Router /api/payments/{id} [put]
func UpdatePayment(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var payment models.Payment
	database.Database.Db.Where("id = ?", id).Find(&payment)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if payment.ID == 0 {
		return c.Status(404).JSON("Payment not found")
	}

	if err := c.BodyParser(&payment); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Save(&payment)
	return c.Status(200).JSON(CreateResponsePayment(payment))
}

func GetPaymentByProductId(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var payments []models.Payment
	database.Database.Db.Where("product_id = ?", id).Find(&payments)

	if err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if len(payments) == 0 {
		return c.Status(404).JSON("Payment not found")
	}

	return c.Status(200).JSON(CreateResponsePayments(payments))
}

func GetAllPaymentsByProductId(c *fiber.Ctx) error {
	var payments []models.Payment
	database.Database.Db.Find(&payments)

	if len(payments) == 0 {
		return c.Status(404).JSON("No payments found")
	}

	return c.Status(200).JSON(CreateResponsePayments(payments))
}
