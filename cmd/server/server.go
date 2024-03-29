package main

import (
	"customer-manager/database"
	"customer-manager/repositories"
	"customer-manager/server"
	"os"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		return ":8080"
	}
	return ":" + port
}

func main() {
	app := fiber.New(fiber.Config{Network: fiber.NetworkTCP})
	db := database.GetDatabase(&gorm.Config{Logger: database.GetLogger(logger.Info)})
	customerManagerServer := server.NewCustomerManagerServer(
		app,
		&repositories.DBCustomerRepository{DB: db},
		&repositories.DBPurchaseRepository{DB: db},
		&repositories.DBRepairRepository{DB: db},
	)

	panic(customerManagerServer.App.Listen(getServerPort()))
}
