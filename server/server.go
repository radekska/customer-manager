package server

import (
	"customer-manager/database"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
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

		return ctx.Status(fiber.StatusOK).JSON(customers)
	})

	server.App.Post("/api/customers", func(ctx *fiber.Ctx) error {
		// TODO check for proper content-type header and inform when invalid.
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

	server.App.Get("/api/customers/:customerID", func(ctx *fiber.Ctx) error {
		customerID := ctx.Params("customerID")
		_, err := uuid.Parse(customerID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", customerID),
			})
		}
		_, customer := server.repository.GetByID(customerID)
		if customer == nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"detail": fmt.Sprintf("customer with given id '%s' does not exists", customerID),
			})
		}
		return ctx.Status(fiber.StatusOK).JSON(customer)
	})

	server.App.Put("/api/customers/:customerID", func(ctx *fiber.Ctx) error {
		// TODO check for proper content-type header and inform when invalid.
		customerID := ctx.Params("customerID")
		_, err := uuid.Parse(customerID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", customerID),
			})
		}
		newCustomerDetails := new(editCustomerDetailsRequest)
		err = ctx.BodyParser(newCustomerDetails)
		if err == fiber.ErrUnprocessableEntity {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": err.Error(),
			})
		}
		validator := getValidator(newCustomerDetails)
		if !validator.Validate() {
			fmt.Println(validator.Errors)
			return ctx.Status(fiber.StatusBadRequest).JSON(validator.Errors)
		}

		_, customer := server.repository.Update(
			&database.Customer{
				ID:              customerID,
				FirstName:       newCustomerDetails.FirstName,
				LastName:        newCustomerDetails.LastName,
				TelephoneNumber: newCustomerDetails.TelephoneNumber,
			},
		)
		if customer == nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"detail": fmt.Sprintf("customer with given id '%s' does not exists", customerID),
			})
		}

		return ctx.Status(fiber.StatusOK).JSON(customer)
	})

	return server
}
