package repositories

import (
	"customer-manager/database"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
)

func clearRecords(t *testing.T, db *gorm.DB) {
	t.Helper()
	tables := []string{"customers", "purchases", "repairs"}
	for _, name := range tables {
		db.Exec(fmt.Sprintf("DELETE FROM %s", name))
	}
}

var db = database.GetDatabase("../test.db", &gorm.Config{Logger: database.GetLogger(logger.Silent)})

func getAllCustomers(t *testing.T, db *gorm.DB) []database.Customer {
	t.Helper()
	var customers []database.Customer
	db.Find(&customers)
	return customers
}

func getPurchaseByID(purchaseID string, t *testing.T, db *gorm.DB) *database.Purchase {
	t.Helper()
	var purchase database.Purchase
	result := db.Where(&database.Purchase{ID: purchaseID}).Find(&purchase)
	if result.RowsAffected == 0 {
		return nil
	}
	return &purchase
}

func getRepairByID(repairID string, t *testing.T, db *gorm.DB) *database.Repair {
	t.Helper()
	var repair database.Repair
	result := db.Where(&database.Repair{ID: repairID}).Find(&repair)
	if result.RowsAffected == 0 {
		return nil
	}
	return &repair
}

func getAllRepairs(t *testing.T, db *gorm.DB) []database.Repair {
	t.Helper()
	var repairs []database.Repair
	db.Find(&repairs)
	return repairs
}

func getAllPurchases(t *testing.T, db *gorm.DB) []database.Purchase {
	t.Helper()
	var purchase []database.Purchase
	db.Find(&purchase)
	return purchase
}

func assertCustomer(t *testing.T, expected *database.Customer, actual *database.Customer) {
	t.Helper()
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.FirstName, actual.FirstName)
	assert.Equal(t, expected.LastName, actual.LastName)
	assert.Equal(t, expected.TelephoneNumber, actual.TelephoneNumber)
	assert.Equal(t, expected.Purchases, actual.Purchases)
	assert.Equal(t, expected.Repairs, actual.Repairs)
}

func TestDBCustomerRepository(t *testing.T) {
	customerRepository := DBCustomerRepository{db}
	repairRepository := DBRepairRepository{db}
	purchaseRepository := DBPurchaseRepository{db}

	customer := &database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123456789"}
	purchase := &database.Purchase{FrameModel: "Model1", LensType: "LensType1",
		LensPower: "LensPower", PD: "CustomPD"}
	repair := &database.Repair{Description: "some issue with the thing", Cost: 12.32}

	clearRecords(t, db)
	t.Run("test create customer", func(t *testing.T) {
		err, dbCustomer := customerRepository.Create(customer)

		assert.NoError(t, err)
		assert.Equal(t, "John", dbCustomer.FirstName)
		assert.Equal(t, "Doe", dbCustomer.LastName)
		assert.Equal(t, "123456789", dbCustomer.TelephoneNumber)
		clearRecords(t, db)
	})

	t.Run("test cannot create customer with the same name and telephone number", func(t *testing.T) {
		err, _ := customerRepository.Create(customer)
		assert.NoError(t, err)

		err, _ = customerRepository.Create(customer)

		dbCustomers := getAllCustomers(t, db)

		assert.Error(t, err, "unique constraint failed error must be present") // compare errors here
		assert.Equal(t, 1, len(dbCustomers))

		clearRecords(t, db)
	})

	t.Run("test delete customer", func(t *testing.T) {
		err, dbCustomer := customerRepository.Create(customer)
		assert.NoError(t, err)
		err, _ = purchaseRepository.Create(customer, purchase)
		assert.NoError(t, err)
		err, _ = repairRepository.Create(customer, repair)
		assert.NoError(t, err)

		err = customerRepository.DeleteByID(dbCustomer.ID)
		assert.NoError(t, err)

		assert.Equal(t, 0, len(getAllCustomers(t, db)))
		assert.Equal(t, 0, len(getAllPurchases(t, db)), "no purchases should be left")
		assert.Equal(t, 0, len(getAllRepairs(t, db)), "no repairs should be left")

		clearRecords(t, db)
	})

	t.Run("test delete customer indempotently", func(t *testing.T) {
		err := customerRepository.DeleteByID(uuid.NewString())
		assert.NoError(t, err)
		assert.Equal(t, 0, len(getAllCustomers(t, db)))
	})

	t.Run("test get all customers", func(t *testing.T) {
		var customers []database.Customer
		customersData := []database.Customer{
			{FirstName: "John", LastName: "Doe", TelephoneNumber: "123"},
			{FirstName: "Jane", LastName: "Doe", TelephoneNumber: "321"},
			{FirstName: "Bob", LastName: "Doe", TelephoneNumber: "893"},
		}
		for _, customerData := range customersData {
			err, customer := customerRepository.Create(
				&database.Customer{FirstName: customerData.FirstName, LastName: customerData.LastName})
			assert.NoError(t, err)
			customers = append(customers, *customer)
		}

		err, dbCustomers := customerRepository.GetAll()

		assert.NoError(t, err)
		for i := 0; i < len(customers); i++ {
			assertCustomer(t, &customers[i], &dbCustomers[i])
		}

		clearRecords(t, db)
	})

	t.Run("test get all customers when no records ", func(t *testing.T) {
		err, dbCustomers := customerRepository.GetAll()

		assert.NoError(t, err)
		assert.Len(t, dbCustomers, 0)
	})
}

func TestDBPurchaseRepository(t *testing.T) {
	customerRepository := DBCustomerRepository{db}
	purchaseRepository := DBPurchaseRepository{db}
	customer := &database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123456789"}
	purchase := &database.Purchase{FrameModel: "Model1", LensType: "LensType1",
		LensPower: "LensPower", PD: "CustomPD"}

	clearRecords(t, db)

	t.Run("test add purchase to a customer", func(t *testing.T) {
		err, dbCustomer := customerRepository.Create(customer)
		assert.NoError(t, err)

		err, dbPurchase := purchaseRepository.Create(dbCustomer, purchase)

		assert.NoError(t, err)

		assert.Equal(t, "Model1", dbPurchase.FrameModel)
		assert.Equal(t, "LensType1", dbPurchase.LensType)
		assert.Equal(t, "LensPower", dbPurchase.LensPower)
		assert.Equal(t, "CustomPD", dbPurchase.PD)
		assert.Equal(t, dbCustomer.ID, dbPurchase.CustomerID)

		clearRecords(t, db)
	})

	t.Run("test remove purchase by ID", func(t *testing.T) {
		err, dbCustomer := customerRepository.Create(customer)
		assert.NoError(t, err)

		err, dbPurchase := purchaseRepository.Create(dbCustomer, purchase)
		assert.NoError(t, err)

		err = purchaseRepository.DeleteByID(dbPurchase.ID)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(getAllCustomers(t, db)))
		assert.Nil(t, getPurchaseByID(dbPurchase.ID, t, db))

		clearRecords(t, db)
	})
}

func TestDBRepairRepository(t *testing.T) {
	customerRepository := DBCustomerRepository{db}
	repairRepository := DBRepairRepository{db: db}
	customer := &database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123456789"}
	repair := &database.Repair{Description: "some issue with the thing", Cost: 12.32}

	t.Run("test add repair to a customer", func(t *testing.T) {
		err, dbCustomer := customerRepository.Create(customer)
		assert.NoError(t, err)

		err, dbRepair := repairRepository.Create(dbCustomer, repair)

		assert.NoError(t, err)

		assert.Equal(t, "some issue with the thing", dbRepair.Description)
		assert.Equal(t, 12.32, dbRepair.Cost)
		assert.Equal(t, dbCustomer.ID, dbRepair.CustomerID)

		clearRecords(t, db)
	})

	t.Run("test remove repair by ID", func(t *testing.T) {
		err, dbCustomer := customerRepository.Create(customer)
		assert.NoError(t, err)

		err, dbRepair := repairRepository.Create(dbCustomer, repair)
		assert.NoError(t, err)

		err = repairRepository.DeleteByID(dbRepair.ID)
		assert.NoError(t, err)

		assert.Equal(t, 1, len(getAllCustomers(t, db)))
		assert.Nil(t, getRepairByID(dbRepair.ID, t, db))

		clearRecords(t, db)
	})
}
