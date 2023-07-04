package database

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func getDatabaseURL() string {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		panic("DATABASE_URL environment variable is not set")
	}
	return dsn
}

func GetDatabase(config *gorm.Config) *gorm.DB {
	dialector := mysql.Open(getDatabaseURL())
	db, err := gorm.Open(dialector, config)
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
