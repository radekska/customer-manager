package main

import (
	"customer-manager/database"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	db := database.GetDatabase("test.db", &gorm.Config{Logger: database.GetLogger(logger.Info)})
	if err := database.RunMigration(db, &database.Purchase{}, &database.Repair{}); err != nil {
		panic("failed during performing migrations")
	}
}
