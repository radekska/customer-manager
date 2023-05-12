package repositories

import (
	"customer-manager/database"
	"fmt"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"testing"
	"time"
)

func clearRecords(t *testing.T, db *gorm.DB) {
	t.Helper()
	tables := []string{"purchases", "repairs", "customers"}
	for _, name := range tables {
		tx := db.Exec(fmt.Sprintf("DELETE FROM %s", name))
		if tx.Error != nil {
			t.Fatal(tx.Error)
		}
	}
}

var db = database.GetDatabase(&gorm.Config{Logger: database.GetLogger(logger.Silent)})

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

func assertPurchase(t *testing.T, expected *database.Purchase, actual *database.Purchase) {
	t.Helper()

	assert.Equal(t, expected.FrameModel, actual.FrameModel)
	assert.Equal(t, expected.LensPower, actual.LensPower)
	assert.Equal(t, expected.LensType, actual.LensType)
}

func getCustomerFixture(t *testing.T) *database.Customer {
	t.Helper()
	return &database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123456789"}
}

func getPurchaseFixture(t *testing.T) *database.Purchase {
	t.Helper()
	return &database.Purchase{FrameModel: "Model1", LensType: "LensType1",
		LensPower: "LensPower", PD: "CustomPD", PurchaseType: "CustomPurchaseType", PurchasedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}
}

func TestDBCustomerRepository(t *testing.T) {
	customerRepository := DBCustomerRepository{db}
	repairRepository := DBRepairRepository{db}
	purchaseRepository := DBPurchaseRepository{db}

	customer := &database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123456789"}
	purchase := &database.Purchase{FrameModel: "Model1", LensType: "LensType1",
		LensPower: "LensPower", PD: "CustomPD", PurchasedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}
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

	t.Run("test delete not existing customer", func(t *testing.T) {
		invalidID := uuid.NewString()
		err := customerRepository.DeleteByID(invalidID)
		assert.Equal(t, err, &CustomerNotFoundError{CustomerID: invalidID})
		assert.Equal(t, 0, len(getAllCustomers(t, db)))
	})

	t.Run("test get all customers", func(t *testing.T) {
		var customers []database.Customer
		customersData := []database.Customer{
			{FirstName: "John", LastName: "Doe", TelephoneNumber: "123"},
			{FirstName: "Jane", LastName: "Doe", TelephoneNumber: "321"},
			{FirstName: "Bob", LastName: "Smith", TelephoneNumber: "893"},
		}
		for _, customerData := range customersData {
			err, customer := customerRepository.Create(
				&database.Customer{
					ID:              "customerID",
					FirstName:       customerData.FirstName,
					LastName:        customerData.LastName,
					TelephoneNumber: customerData.TelephoneNumber,
				},
			)
			assert.NoError(t, err)
			customers = append(customers, *customer)
		}

		err, dbCustomers := customerRepository.ListBy("J", "Do")

		assert.NoError(t, err)
		assertCustomer(t, &customers[0], &dbCustomers[0])
		assertCustomer(t, &customers[1], &dbCustomers[1])
		assert.Len(t, dbCustomers, 2)

		clearRecords(t, db)
	})

	t.Run("test get all customers when no records ", func(t *testing.T) {
		err, dbCustomers := customerRepository.ListBy("", "")

		assert.NoError(t, err)
		assert.Len(t, dbCustomers, 0)
	})

	t.Run("test get customer by its id", func(t *testing.T) {
		err, expectedCustomer := customerRepository.Create(
			&database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123456789"},
		)
		assert.NoError(t, err)
		err, _ = customerRepository.Create(
			&database.Customer{FirstName: "Jane", LastName: "Doe", TelephoneNumber: "987456123"},
		)
		assert.NoError(t, err)

		err, currentCustomer := customerRepository.GetByID(expectedCustomer.ID)

		assert.NoError(t, err)
		assert.Equal(t, expectedCustomer.ID, currentCustomer.ID)
		assert.Equal(t, expectedCustomer.FirstName, currentCustomer.FirstName)
		assert.Equal(t, expectedCustomer.LastName, currentCustomer.LastName)
		assert.Equal(t, expectedCustomer.TelephoneNumber, currentCustomer.TelephoneNumber)

		clearRecords(t, db)
	})

	t.Run("test get customer by its id but not found", func(t *testing.T) {
		invalidID := "4a923682-b051-47c1-b37a-666544d71419"
		err, _ := customerRepository.Create(
			&database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123456789"},
		)
		assert.NoError(t, err)

		err, currentCustomer := customerRepository.GetByID(invalidID)

		assert.Equal(t, err, &CustomerNotFoundError{CustomerID: invalidID})
		assert.Nil(t, currentCustomer)
		clearRecords(t, db)
	})

	t.Run("test edit customer details", func(t *testing.T) {
		err, existingCustomer := customerRepository.Create(
			&database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123456789"},
		)
		assert.NoError(t, err)
		updatedCustomer := &database.Customer{
			ID:              existingCustomer.ID,
			FirstName:       "Bob",
			LastName:        "Toe",
			TelephoneNumber: "897564321",
		}

		err, returnedCustomer := customerRepository.Update(updatedCustomer)

		_, dbCustomer := customerRepository.GetByID(returnedCustomer.ID)
		assertCustomer(t, updatedCustomer, dbCustomer)
		assertCustomer(t, updatedCustomer, returnedCustomer)
		clearRecords(t, db)
	})

	t.Run("test edit customer details but not found", func(t *testing.T) {
		// TODO
	})

}

func TestDBPurchaseRepository(t *testing.T) {
	customerRepository := DBCustomerRepository{db}
	purchaseRepository := DBPurchaseRepository{db}

	clearRecords(t, db)

	t.Run("test get all purchases for a customer", func(t *testing.T) {
		customer := getCustomerFixture(t)
		purchase1 := database.Purchase{FrameModel: "Model1", LensType: "LensType1",
			LensPower: "LensPower1", PD: "CustomPD1", PurchasedAt: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)}
		purchase2 := database.Purchase{FrameModel: "Model2", LensType: "LensType2",
			LensPower: "LensPower2", PD: "CustomPD2", PurchasedAt: time.Date(2019, 1, 1, 0, 0, 0, 0, time.UTC)}
		customer.Purchases = []database.Purchase{purchase1, purchase2}
		err, customer := customerRepository.Create(customer)
		assert.NoError(t, err)

		err, purchases := purchaseRepository.GetAll(customer.ID)

		assert.NoError(t, err)
    assert.Len(t, purchases, 2)
		assertPurchase(t, &purchase1, &purchases[0])
		assertPurchase(t, &purchase2, &purchases[1])

		clearRecords(t, db)
	})

	t.Run("test add purchase to a customer", func(t *testing.T) {
		err, dbCustomer := customerRepository.Create(getCustomerFixture(t))
		assert.NoError(t, err)

		err, dbPurchase := purchaseRepository.Create(dbCustomer, getPurchaseFixture(t))

		assert.NoError(t, err)

		assert.Equal(t, "Model1", dbPurchase.FrameModel)
		assert.Equal(t, "LensType1", dbPurchase.LensType)
		assert.Equal(t, "LensPower", dbPurchase.LensPower)
		assert.Equal(t, "CustomPD", dbPurchase.PD)
		assert.Equal(t, "CustomPurchaseType", dbPurchase.PurchaseType)
		assert.Equal(t, dbCustomer.ID, dbPurchase.CustomerID)
		assert.Equal(t, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC), dbPurchase.PurchasedAt)

		clearRecords(t, db)
	})

	t.Run("test update purchase details", func(t *testing.T) {
		err, dbCustomer := customerRepository.Create(getCustomerFixture(t))
		assert.NoError(t, err)
		err, dbPurchase := purchaseRepository.Create(dbCustomer, getPurchaseFixture(t))
		assert.NoError(t, err)
		updatedPurchase := &database.Purchase{
			ID:           dbPurchase.ID,
			CustomerID:   dbCustomer.ID,
			FrameModel:   "UpdatedModel",
			LensType:     "UpdatedLensType",
			LensPower:    "UpdatedLensPower",
			PD:           "UpdatedPD",
			PurchaseType: "UpdatedPurchaseType",
			PurchasedAt:  time.Date(2000, 10, 20, 15, 0, 0, 0, time.UTC),
		}

		err, updatedDbPurchase := purchaseRepository.Update(updatedPurchase)

		err, dbPurchases := purchaseRepository.GetAll(dbCustomer.ID)
		assert.NoError(t, err)
		assert.Len(t, dbPurchases, 1)
		assertPurchase(t, &dbPurchases[0], updatedDbPurchase)
		clearRecords(t, db)
	})

	t.Run("test remove purchase by ID", func(t *testing.T) {
		err, dbCustomer := customerRepository.Create(getCustomerFixture(t))
		assert.NoError(t, err)
		err, dbPurchase := purchaseRepository.Create(dbCustomer, getPurchaseFixture(t))
		assert.NoError(t, err)

		err = purchaseRepository.DeleteByID(dbPurchase.ID)

		assert.NoError(t, err)
		assert.Equal(t, 1, len(getAllCustomers(t, db)))
		assert.Nil(t, getPurchaseByID(dbPurchase.ID, t, db))
		clearRecords(t, db)
	})

	t.Run("test remove purchase by ID but not found", func(t *testing.T) {
		err := purchaseRepository.DeleteByID("4a923682-1234-47c1-b37a-666544d71419")

		assert.Equal(t, err, &PurchaseNotFoundError{PurchaseID: "4a923682-1234-47c1-b37a-666544d71419"})
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
