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

// getCustomersHandler godoc
//
//	@Summary		Get list of customers
//	@Description	Returns full list of existing customers
//	@Tags			list-customers
//	@Produce		json
//	@Success		200	{array} database.Customer
//	@Param        	firstName    query     string  false  "first name search"
//	@Param        	lastName    query     string  false  "last name search"
//	@Router			/api/customers [get]
func getCustomersHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		firstName := ctx.Query("firstName")
		lastName := ctx.Query("lastName")
		err, customers := server.customerRepository.ListBy(firstName, lastName)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return ctx.Status(fiber.StatusOK).JSON(customers)
	}
}

// createCustomerHandler godoc
//
//	@Summary		Create customer
//	@Description	Create customer object
//	@Tags			create-customer
//	@Accept			json
//	@Produce		json
//	@Success		201	{object} database.Customer
//	@Failure		400	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Param			customerDetails	body	server.CreateCustomerRequest	true "Customer details"
//	@Router			/api/customers [post]
func createCustomerHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		newCustomer := new(CreateCustomerRequest)
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

		err, customer := server.customerRepository.Create(
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
		return ctx.Status(fiber.StatusCreated).JSON(customer)
	}
}

// getCustomerByIDHandler godoc
//
//	@Summary		Get customer
//	@Description	Returns customer details by ID
//	@Tags			get-customer
//	@Produce		json
//	@Success		200	{object} database.Customer
//	@Failure		400	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Failure		404	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Param			customerID	path	string	true "Customer ID"
//	@Router			/api/customers/{customerID} [get]
func getCustomerByIDHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		customerID := ctx.Params("customerID")
		_, err := uuid.Parse(customerID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", customerID),
			})
		}
		_, customer := server.customerRepository.GetByID(customerID)
		if customer == nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"detail": fmt.Sprintf("customer with given id '%s' does not exists", customerID),
			})
		}
		return ctx.Status(fiber.StatusOK).JSON(customer)
	}
}

// editCustomerByIDHandler godoc
//
//	@Summary		Edit customer
//	@Description	Edit customer details by ID
//	@Tags			edit-customer
//	@Produce		json
//	@Success		200	{object} database.Customer
//	@Failure		400	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Failure		404	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Param			customerID	path	string	true "Customer ID"
//	@Param			customerDetails	body	server.EditCustomerDetailsRequest	true "New customer details"
//	@Router			/api/customers/{customerID} [put]
func editCustomerByIDHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		customerID := ctx.Params("customerID")
		_, err := uuid.Parse(customerID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", customerID),
			})
		}
		newCustomerDetails := new(EditCustomerDetailsRequest)
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

		_, customer := server.customerRepository.Update(
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

// deleteCustomerByIDHandler godoc
//
//	@Summary		Delete customer
//	@Description	Delete customer details and it's relations by ID
//	@Tags			delete-customer
//	@Success		204
//	@Failure		400	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Failure		404	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Param			customerID	path	string	true "Customer ID"
//	@Router			/api/customers/{customerID} [delete]
func deleteCustomerByIDHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		customerID := ctx.Params("customerID")
		_, err := uuid.Parse(customerID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", customerID),
			})
		}
		if err := server.customerRepository.DeleteByID(customerID); err != nil {
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

// getPurchasesHandler godoc
//
//	@Summary		Get list of purchases
//	@Description	Returns full list of purchases for a specific customer by ID
//	@Tags			get-customer-purchases
//	@Produce		json
//	@Success		200	{array} database.Purchase
//	@Param			customerID	path	string	true "Customer ID"
//	@Router			/api/customers/{customerID}/purchases [get]
func getPurchasesHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		customerID := ctx.Params("customerID")
		_, err := uuid.Parse(customerID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", customerID),
			})
		}
		err, purchases := server.purchasesRepository.GetAll(customerID)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return ctx.Status(fiber.StatusOK).JSON(purchases)
	}
}
