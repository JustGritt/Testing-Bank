package routes

import (
	"github.com/JustGritt/go-grpc/database"
	"github.com/JustGritt/go-grpc/models"
	"github.com/gofiber/fiber/v2"
)

type Product struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Price uint   `json:"price"`
}

func CreateResponseProduct(product models.Product) Product {
	return Product{
		ID:    product.ID,
		Name:  product.Name,
		Price: product.Price,
	}
}

// Create product godoc
// @Summary Create product
// @Description Create product
// @Tags products
// @Accept  json
// @Produce  json
// @Param product body Product true "Product"
// @Success 200 {object} Product
// @Failure 400 {string} string "Product name is required"
// @Failure 400 {string} string "Product name already taken"
// @Router /api/products [post]
func CreateProduct(c *fiber.Ctx) error {
	var product models.Product

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if product.Name == "" {
		return c.Status(400).JSON("Product name is required")
	}

	var existingProduct models.Product
	database.Database.Db.Where("name = ?", product.Name).First(&existingProduct)
	if existingProduct.ID != 0 {
		return c.Status(400).JSON("Product name already taken")
	}

	database.Database.Db.Create(&product)
	responseProduct := CreateResponseProduct(product)
	return c.Status(200).JSON(responseProduct)
}

func GetProductId(id int, product *models.Product) error {
	database.Database.Db.Find(&product, id)
	if product.ID == 0 {
		return fiber.NewError(404, "Product not found")
	}

	return nil
}

// Get all products godoc
// @Summary Get all products
// @Description Get all products
// @Tags products
// @Accept  json
// @Produce  json
// @Success 200 {object} []Product
// @Failure 400 {string} string "Product not found"
// @Router /api/products [get]
func GetProducts(c *fiber.Ctx) error {
	var products []models.Product
	database.Database.Db.Find(&products)

	var responseProducts []Product
	for _, product := range products {
		responseProducts = append(responseProducts, CreateResponseProduct(product))
	}

	return c.Status(200).JSON(responseProducts)
}

// Get product by id godoc
// @Summary Get product by id
// @Description Get product by id
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {object} Product
// @Failure 400 {string} string "Please ensure that :id is an integer"
// @Failure 400 {string} string "Product not found"
// @Router /api/products/{id} [get]
func GetProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var product models.Product

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := GetProductId(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	responseProduct := CreateResponseProduct(product)

	return c.Status(200).JSON(responseProduct)
}

// Update product godoc
// @Summary Update product
// @Description Update product
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Param product body Product true "Product"
// @Success 200 {object} Product
// @Failure 400 {string} string "Please ensure that :id is an integer"
// @Failure 400 {string} string "Product not found"
// @Failure 400 {string} string "Product name already taken"
// @Router /api/products/{id} [put]
func UpdateProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var product models.Product

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := GetProductId(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	if err := c.BodyParser(&product); err != nil {
		return c.Status(400).JSON("Invalid request")
	}

	database.Database.Db.Save(&product)

	responseProduct := CreateResponseProduct(product)

	return c.Status(200).JSON(responseProduct)
}

// Delete product godoc
// @Summary Delete product
// @Description Delete product
// @Tags products
// @Accept  json
// @Produce  json
// @Param id path int true "Product ID"
// @Success 200 {string} string "Product deleted"
// @Failure 400 {string} string "Please ensure that :id is an integer"
// @Failure 400 {string} string "Product not found"
// @Router /api/products/{id} [delete]
func DeleteProduct(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")

	var product models.Product

	if err != nil {
		return c.Status(400).JSON("Please ensure that :id is an integer")
	}

	if err := GetProductId(id, &product); err != nil {
		return c.Status(400).JSON(err.Error())
	}

	database.Database.Db.Delete(&product)

	return c.Status(200).JSON("Product deleted")
}
