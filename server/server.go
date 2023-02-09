package server

import (
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
	customersPath := "/api/customers"
	server.App.Use(jsonContentValidator())
	server.App.Get(customersPath, getCustomersHandler(server))
	server.App.Post(customersPath, createCustomerHandler(server))
	server.App.Get(customersPath+"/:customerID", getCustomerByIDHandler(server))
	server.App.Put(customersPath+"/:customerID", editCustomerByIDHandler(server))
	server.App.Delete(customersPath+"/:customerID", deleteCustomerByIDHandler(server))

	return server
}
