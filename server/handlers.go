package server

import (
	"customer-manager/database"
	"customer-manager/repositories"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

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

func genericListHandler[V []database.Purchase | []database.Repair](getAll func(customerID string) (error, V)) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		customerID := ctx.Params("customerID")
		_, err := uuid.Parse(customerID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", customerID),
			})
		}
		err, items := getAll(customerID)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		return ctx.Status(fiber.StatusOK).JSON(items)
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
	return genericListHandler[[]database.Purchase](server.purchasesRepository.GetAll)
}

// createPurchaseHandler godoc
//
//	@Summary		Create a purchase for a customer
//	@Description	Creates a new purchase for a customer by ID
//	@Tags			create-customer-purchase
//	@Accept			json
//	@Produce		json
//	@Success		200	{object} database.Purchase
//	@Failure		404	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Failure		400	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Param			customerID	path	string	true "Customer ID"
//	@Param			purchaseDetails	body	server.CreatePurchaseRequest	true "Purchase details"
//	@Router			/api/customers/{customerID}/purchases [post]
func createPurchaseHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		newPurchase := new(CreatePurchaseRequest)
		err := ctx.BodyParser(newPurchase)
		if err == fiber.ErrUnprocessableEntity {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": err.Error(),
			})
		}

		validator := getValidator(newPurchase)
		if !validator.Validate() {
			return ctx.Status(fiber.StatusBadRequest).JSON(validator.Errors)
		}

		customerID := ctx.Params("customerID")
		_, err = uuid.Parse(customerID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", customerID),
			})
		}
		err, customer := server.customerRepository.GetByID(customerID)
		if errors.Is(err, &repositories.CustomerNotFoundError{}) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"detail": fmt.Sprintf("customer with given id '%s' does not exists", customerID),
			})
		}

		err, purchase := server.purchasesRepository.Create(customer, &database.Purchase{
			FrameModel:   newPurchase.FrameModel,
			LensType:     newPurchase.LensType,
			LensPower:    newPurchase.LensPower,
			PD:           newPurchase.PD,
			PurchaseType: newPurchase.PurchaseType,
			PurchasedAt:  time.Time(newPurchase.PurchasedAt),
		})
		if err != nil {
			// TODO: handle error
			return err
		}
		return ctx.Status(fiber.StatusCreated).JSON(purchase)
	}
}

// deletePurchaseByIDHandler godoc
//
//	@Summary		Delete a purchase
//	@Description	Deletes a purchase by ID
//	@Tags			delete-customer-purchase
//	@Success		204
//	@Failure		404	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Failure		400	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Param			customerID	path	string	true "Customer ID"
//	@Param			purchaseID	path	string	true "Purchase ID"
//	@Router			/api/customers/{customerID}/purchases/{purchaseID} [delete]
func deletePurchaseByIDHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		customerID := ctx.Params("customerID")
		_, err := uuid.Parse(customerID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", customerID),
			})
		}
		purchaseID := ctx.Params("purchaseID")
		_, err = uuid.Parse(purchaseID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given purchase id '%s' is not a valid UUID", purchaseID),
			})
		}

		if err := server.purchasesRepository.DeleteByID(purchaseID); err != nil {
			target := &repositories.PurchaseNotFoundError{}
			if errors.As(err, &target) {
				return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
					"detail": fmt.Sprintf("purchase with given id '%s' does not exists", purchaseID),
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

// editPurchaseByIDHandler godoc
//
//	@Summary		Update a purchase
//	@Description	Updates a purchase for a customer by ID
//	@Tags			update-customer-purchase
//	@Success		200	{object} database.Purchase
//	@Failure		404	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Failure		400	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Param			customerID	path	string	true "Customer ID"
//	@Param			purchaseID	path	string	true "Purchase ID"
//	@Param			customerDetails	body	server.EditPurchaseRequest	true "New purchase details"
//	@Router			/api/customers/{customerID}/purchases/{purchaseID} [put]
func editPurchaseByIDHandler(server *CustomerManagerServer) fiber.Handler {
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

		purchaseID := ctx.Params("purchaseID")
		_, err = uuid.Parse(purchaseID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given purchase id '%s' is not a valid UUID", purchaseID),
			})
		}

		newPurchaseDetails := new(EditPurchaseRequest)
		err = ctx.BodyParser(newPurchaseDetails)
		if err == fiber.ErrUnprocessableEntity {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": err.Error(),
			})
		}
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": err,
			})
		}

		validator := getValidator(newPurchaseDetails)
		if !validator.Validate() {
			fmt.Println(validator.Errors)
			return ctx.Status(fiber.StatusBadRequest).JSON(validator.Errors)
		}

		_, purchase := server.purchasesRepository.Update(
			&database.Purchase{
				ID:           purchaseID,
				FrameModel:   newPurchaseDetails.FrameModel,
				LensType:     newPurchaseDetails.LensType,
				LensPower:    newPurchaseDetails.LensPower,
				PD:           newPurchaseDetails.PD,
				CustomerID:   customerID,
				PurchaseType: newPurchaseDetails.PurchaseType,
				PurchasedAt:  time.Time(newPurchaseDetails.PurchasedAt),
			},
		)
		if purchase == nil {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"detail": fmt.Sprintf("purchase with given id '%s' does not exists", customerID),
			})
		}
		// TODO during update the "created_at": "0001-01-01T00:00:00Z" is zeroed
		return ctx.Status(fiber.StatusOK).JSON(purchase)
	}
}

// getRepairsHandler godoc
//
//	@Summary		Get list of repairs
//	@Description	Returns full list of repairs for a specific customer by ID
//	@Tags			get-customer-repairs
//	@Produce		json
//	@Success		200	{array} database.Repair
//	@Param			customerID	path	string	true "Customer ID"
//	@Router			/api/customers/{customerID}/repairs [get]
func getRepairsHandler(server *CustomerManagerServer) fiber.Handler {
	return genericListHandler[[]database.Repair](server.repairsRepository.GetAll)
}

// createRepairHandler godoc
//
//	@Summary		  Create a repair for a customer
//	@Description	Creates a new repair for a customer by ID
//	@Tags			    create-customer-repair
//	@Accept			  json
//	@Produce		  json
//	@Success		  200	{object} database.Repair
//	@Failure		  404	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Failure		  400	{string} string "IMPLEMENTED BUT DOCS TODO"
//	@Param			  customerID	path	string	true "Customer ID"
//	@Param			  repairDetails	body	server.CreateRepairRequest	true "Repair details"
//	@Router			  /api/customers/{customerID}/repairs [post]
func createRepairHandler(server *CustomerManagerServer) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		newRepair := new(CreateRepairRequest)
		err := ctx.BodyParser(newRepair)
		if err == fiber.ErrUnprocessableEntity {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors": err.Error(),
			})
		}

		validator := getValidator(newRepair)
		if !validator.Validate() {
			return ctx.Status(fiber.StatusBadRequest).JSON(validator.Errors)
		}

		customerID := ctx.Params("customerID")
		_, err = uuid.Parse(customerID)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", customerID),
			})
		}
		err, customer := server.customerRepository.GetByID(customerID)
		if errors.Is(err, &repositories.CustomerNotFoundError{}) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"detail": fmt.Sprintf("customer with given id '%s' does not exists", customerID),
			})
		}

		err, repair := server.repairsRepository.Create(customer, &database.Repair{
			Description: newRepair.Description,
			Cost:        newRepair.Cost,
			ReportedAt:  time.Time(newRepair.ReportedAt),
		})
		if err != nil {
			// TODO: handle error
			return err
		}
		return ctx.Status(fiber.StatusCreated).JSON(repair)
	}
}
