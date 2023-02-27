package server

import (
	"github.com/gofiber/fiber/v2"
)
import "customer-manager/repositories"

type CustomerManagerServer struct {
	App                 *fiber.App
	customerRepository  repositories.CustomerRepository
	purchasesRepository repositories.PurchaseRepository
}

func NewCustomerManagerServer(
	app *fiber.App,
	customerRepository repositories.CustomerRepository,
	purchasesRepository repositories.PurchaseRepository,
) *CustomerManagerServer {
	server := &CustomerManagerServer{
		App:                 app,
		customerRepository:  customerRepository,
		purchasesRepository: purchasesRepository,
	}

	server.App.Use(jsonContentValidator())

	customersPath := "/api/customers"

	server.App.Get(customersPath, getCustomersHandler(server))
	server.App.Post(customersPath, createCustomerHandler(server))
	server.App.Get(customersPath+"/:customerID", getCustomerByIDHandler(server))
	server.App.Put(customersPath+"/:customerID", editCustomerByIDHandler(server))
	server.App.Delete(customersPath+"/:customerID", deleteCustomerByIDHandler(server))

	purchasesPath := customersPath + "/:customerID" + "/purchases"

	server.App.Get(purchasesPath, getPurchasesHandler(server))

	return server
}
