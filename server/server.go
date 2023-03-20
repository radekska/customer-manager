package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)
import "customer-manager/repositories"

type CustomerManagerServer struct {
	App                 *fiber.App
	customerRepository  repositories.CustomerRepository
	purchasesRepository repositories.PurchaseRepository
}

func mountMiddlewares(server *CustomerManagerServer) {
	server.App.Use(jsonContentValidator())
	server.App.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
	}))
	server.App.Use(logger.New()) // TODO - do not log requests during tests
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

	mountMiddlewares(server)

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
