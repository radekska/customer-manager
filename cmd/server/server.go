package main

import (
	"customer-manager/database"
	"customer-manager/repositories"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	app := fiber.New()
	db := database.GetDatabase("../test.db", &gorm.Config{Logger: database.GetLogger(logger.Silent)})
	repository := repositories.DBCustomerRepository{db}
	app.Get("/api/customers", func(c *fiber.Ctx) error {
		_, customer := repository.Create(&database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123456789"})

		return c.JSON(customer)
	})

	// TODO - https://github.com/radekska/customer-manager/issues/4

	app.Listen(":3000")
}
