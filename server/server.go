package server

import (
	"customer-manager/database"
	"github.com/gofiber/fiber/v2"
)
import "customer-manager/repositories"

type CustomerManagerServer struct {
	App        *fiber.App
	repository repositories.CustomerRepository
}

func NewCustomerManagerServer(app *fiber.App, repository repositories.CustomerRepository) *CustomerManagerServer {
	server := &CustomerManagerServer{
		App:        app,
		repository: repository,
	}
	server.App.Get("/api/customers", func(ctx *fiber.Ctx) error {
		err, customers := server.repository.GetAll()
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return ctx.JSON(customers)
	})

	server.App.Post("/api/customers", func(ctx *fiber.Ctx) error {
		c := new(database.Customer)
		err := ctx.BodyParser(c)
		if err == fiber.ErrUnprocessableEntity {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": err.Error(),
			})
		}

		errors := fiber.Map{}
		if c.FirstName == "" {
			errors["first_name"] = "this field is required"
		}
		if c.LastName == "" {
			errors["last_name"] = "this field is required"
		}
		if c.TelephoneNumber == "" {
			errors["telephone_number"] = "this field is required"
		}
		if len(errors) != 0 {
			return ctx.Status(fiber.StatusBadRequest).JSON(errors)
		}

		err, _ = server.repository.Create(c)
		if err != nil {
			return err
		}
		ctx.Status(fiber.StatusCreated)
		return nil
	})

	return server
}
