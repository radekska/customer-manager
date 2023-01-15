package server

import (
	"customer-manager/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/validate"
)
import "customer-manager/repositories"

type CustomerManagerServer struct {
	App        *fiber.App
	repository repositories.CustomerRepository
}

func getValidator(s interface{}) *validate.Validation {
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})
	v := validate.New(s)
	v.AddMessages(map[string]string{
		"required": "The '{field}' is required",
	})
	return v
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

		validator := getValidator(c)
		if !validator.Validate() {
			return ctx.Status(fiber.StatusBadRequest).JSON(validator.Errors)
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
