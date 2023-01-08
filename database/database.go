package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

func GetDatabase(dsn string, config *gorm.Config) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(dsn), config)
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func GetLogger(logLevel logger.LogLevel) logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Second,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)
}

func RunMigration(db *gorm.DB, schemas ...interface{}) error {
	return db.AutoMigrate(schemas...)
}
