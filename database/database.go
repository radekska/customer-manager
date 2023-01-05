package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Customer struct {
	gorm.Model
	FirstName string
	LastName  string
	Age       int
}

func GetDatabase(dsn string, config *gorm.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), config)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func RunMigration(db *gorm.DB, schemas ...interface{}) {
	err := db.AutoMigrate(schemas...)
	if err != nil {
		panic("failed during migration")
	}
}
