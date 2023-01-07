package repositories

import (
	"customer-manager/database"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"testing"
	"time"
)

var newLogger = logger.New(
	log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
	logger.Config{
		SlowThreshold:             time.Second,   // Slow SQL threshold
		LogLevel:                  logger.Silent, // Log level
		IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
		Colorful:                  true,          // Disable color
	},
)

func clearRecords(t *testing.T, db *gorm.DB) {
	t.Helper()
	tables := []string{"customers", "purchases", "repairs"}
	for _, name := range tables {
		db.Exec(fmt.Sprintf("DELETE FROM %s", name))
	}
}

var db = database.GetDatabase("../test.db", &gorm.Config{Logger: newLogger})

func TestDBCustomerRepository(t *testing.T) {
	customerRepository := DBCustomerRepository{db}
	customer := &database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123456789"}

	clearRecords(t, db)
	t.Run("test create customer", func(t *testing.T) {
		err := customerRepository.Create(customer)

		var dbCustomer database.Customer
		db.Where("first_name = ? AND last_name = ?", "John", "Doe").First(&dbCustomer)

		assert.NoError(t, err)
		assert.Equal(t, "John", dbCustomer.FirstName)
		assert.Equal(t, "Doe", dbCustomer.LastName)
		assert.Equal(t, "123456789", dbCustomer.TelephoneNumber)
		clearRecords(t, db)
	})

	t.Run("test cannot create customer with the same name and telephone number", func(t *testing.T) {
		err := customerRepository.Create(customer)
		assert.NoError(t, err)

		err = customerRepository.Create(customer)

		var dbCustomers []database.Customer
		db.Where("first_name = ? AND last_name = ?", "John", "Doe").Find(&dbCustomers)

		assert.Error(t, err, "unique constraint failed error must be present")
		assert.Equal(t, 1, len(dbCustomers))

		clearRecords(t, db)
	})

	t.Run("test delete customer", func(t *testing.T) {

	}) // TODO

	t.Run("test delete customer indempotently", func(t *testing.T) {

	}) // TODO
}
func TestDBPurchaseRepository(t *testing.T) {
	customerRepository := DBCustomerRepository{db}
	purchaseRepository := DBPurchaseRepository{db: db}
	customer := &database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123456789"}

	t.Run("test add purchase to a customer", func(t *testing.T) {
		err := customerRepository.Create(customer)
		assert.NoError(t, err)

		var dbCustomer database.Customer
		db.Where("first_name = ? AND last_name = ?", "John", "Doe").First(&dbCustomer)

		err = purchaseRepository.Create(&dbCustomer, &database.Purchase{FrameModel: "Model1", LensType: "LensType1",
			LensPower: "LensPower", PD: "CustomPD"})

		assert.NoError(t, err)
		var dbPurchase database.Purchase
		db.Where("customer_id = ?", dbCustomer.ID).First(&dbPurchase)

		assert.Equal(t, "Model1", dbPurchase.FrameModel)
		assert.Equal(t, "LensType1", dbPurchase.LensType)
		assert.Equal(t, "LensPower", dbPurchase.LensPower)
		assert.Equal(t, "CustomPD", dbPurchase.PD)
		assert.Equal(t, dbCustomer.ID, dbPurchase.CustomerID)

		clearRecords(t, db)
	})

	t.Run("test remove purchase", func(t *testing.T) {

	}) // TODO
}

func TestDBRepairRepository(t *testing.T) {
	customerRepository := DBCustomerRepository{db}
	repairRepository := DBRepairRepository{db: db}
	customer := &database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123456789"}

	t.Run("test add repair to a customer", func(t *testing.T) {
		err := customerRepository.Create(customer)
		assert.NoError(t, err)

		var dbCustomer database.Customer
		db.Where("first_name = ? AND last_name = ?", "John", "Doe").First(&dbCustomer)

		err = repairRepository.Create(&dbCustomer, &database.Repair{Description: "some issue with the thing", Cost: 12.32})

		assert.NoError(t, err)

		var dbRepair database.Repair
		db.Where("customer_id = ?", dbCustomer.ID).First(&dbRepair)

		assert.Equal(t, "some issue with the thing", dbRepair.Description)
		assert.Equal(t, 12.32, dbRepair.Cost)
		assert.Equal(t, dbCustomer.ID, dbRepair.CustomerID)

		clearRecords(t, db)
	})

	t.Run("test remove repair", func(t *testing.T) {

	}) // TODO
}
