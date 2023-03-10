package main

import (
	"customer-manager/database"
	"customer-manager/repositories"
	"customer-manager/server"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	app := fiber.New()
	db := database.GetDatabase(
		"/home/rskalbania/GolandProjects/customer-manager/test.db",
		&gorm.Config{Logger: database.GetLogger(logger.Info)},
	)
	customerManagerServer := server.NewCustomerManagerServer(
		app,
		&repositories.DBCustomerRepository{DB: db},
		&repositories.DBPurchaseRepository{DB: db},
	)

	panic(customerManagerServer.App.Listen(":8080"))

}
