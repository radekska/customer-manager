package server

import (
	"customer-manager/database"
	"customer-manager/repositories"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/gookit/validate"
)

func getCustomersHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		err, customers := server.repository.GetAll()
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return ctx.Status(fiber.StatusOK).JSON(customers)
	}
}

func createCustomerHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		newCustomer := new(createCustomerRequest)
		err := ctx.BodyParser(newCustomer)
		if err == fiber.ErrUnprocessableEntity {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": err.Error(),
			})
		}

		validator := getValidator(newCustomer)
		if !validator.Validate() {
			return ctx.Status(fiber.StatusBadRequest).JSON(validator.Errors)
		}

		err, _ = server.repository.Create(
			&database.Customer{
				FirstName:       newCustomer.FirstName,
				LastName:        newCustomer.LastName,
				TelephoneNumber: newCustomer.TelephoneNumber,
			},
		)
		if err != nil {
			// TODO - fix DB error message 'UNIQUE constraint failed: customers.telephone_number'
			return err
		}
		ctx.Status(fiber.StatusCreated)
		return nil
	}
}

func getCustomerByIDHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
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
	}
}

func editCustomerByIDHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
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
		// TODO during update the "created_at": "0001-01-01T00:00:00Z" is zeroed
		return ctx.Status(fiber.StatusOK).JSON(customer)
	}
}

func deleteCustomerByIDHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		customerID := ctx.Params("customerID")
		_, err := uuid.Parse(customerID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", customerID),
			})
		}
		if err := server.repository.DeleteByID(customerID); err != nil {
			target := &repositories.CustomerNotFoundError{}
			if errors.As(err, &target) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"detail": fmt.Sprintf("customer with given id '%s' does not exists", customerID),
				})
			}
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"detail": err,
			})
		}
		ctx.Status(fiber.StatusNoContent)
		return nil
	}
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
