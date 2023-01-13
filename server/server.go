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
	server.App.Get("/api/customers", func(ctx *fiber.Ctx) error {
		err, customers := server.repository.GetAll()
		if err != nil {
			return fiber.ErrInternalServerError
		}

		return ctx.JSON(customers)
	})

	return server
}
