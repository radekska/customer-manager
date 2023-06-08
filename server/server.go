package server

import (
	_ "customer-manager/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/swagger"
)
import "customer-manager/repositories"

type CustomerManagerServer struct {
	App                 *fiber.App
	customerRepository  repositories.CustomerRepository
	purchasesRepository repositories.PurchaseRepository
	repairsRepository   repositories.RepairRepository
}

func mountMiddlewares(server *CustomerManagerServer) {
	server.App.Use(jsonContentValidator())
	server.App.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:3000",
	}))
	server.App.Use(logger.New()) // TODO - do not log requests during tests
}

func mountSwaggerDocs(server *CustomerManagerServer) {
	server.App.Get("/swagger/*", swagger.HandlerDefault) // default
}

func NewCustomerManagerServer(
	app *fiber.App,
	customerRepository repositories.CustomerRepository,
	purchasesRepository repositories.PurchaseRepository,
	repairsRepository repositories.RepairRepository,
) *CustomerManagerServer {
	server := &CustomerManagerServer{
		App:                 app,
		customerRepository:  customerRepository,
		purchasesRepository: purchasesRepository,
		repairsRepository:   repairsRepository,
	}

	mountMiddlewares(server)
	mountSwaggerDocs(server)

	customersPath := "/api/customers"

	server.App.Get(customersPath, getCustomersHandler(server))
	server.App.Post(customersPath, createCustomerHandler(server))
	server.App.Get(customersPath+"/:customerID", getCustomerByIDHandler(server))
	server.App.Put(customersPath+"/:customerID", editCustomerByIDHandler(server))
	server.App.Delete(customersPath+"/:customerID", deleteCustomerByIDHandler(server))

	purchasesPath := customersPath + "/:customerID" + "/purchases"
	server.App.Get(purchasesPath, getPurchasesHandler(server))
	server.App.Post(purchasesPath, createPurchaseHandler(server))
	server.App.Delete(purchasesPath+"/:purchaseID", deletePurchaseByIDHandler(server))
	server.App.Put(purchasesPath+"/:purchaseID", editPurchaseByIDHandler(server))

 	repairsPath := customersPath + "/:customerID" + "/repairs"
  server.App.Get(repairsPath, getRepairsHandler(server))

	return server
}
