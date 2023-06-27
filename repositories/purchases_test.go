package repositories

import (
	"customer-manager/database"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
