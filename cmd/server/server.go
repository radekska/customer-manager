package main

import (
	"customer-manager/database"
	"customer-manager/repositories"
	"customer-manager/server"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
)

func getServerPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = ":8080"
	}
	return ":" + port
}
func getDatabaseURL() string {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_URL environment variable is not set")
	}
	return dsn
}

func main() {
	app := fiber.New()
	db := database.GetDatabase(
		getDatabaseURL(),
		&gorm.Config{Logger: database.GetLogger(logger.Info)},
	)
	customerManagerServer := server.NewCustomerManagerServer(
		app,
		&repositories.DBCustomerRepository{DB: db},
		&repositories.DBPurchaseRepository{DB: db},
	)

	panic(customerManagerServer.App.Listen(getServerPort()))
}
