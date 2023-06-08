package repositories

import (
	"customer-manager/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

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
			{FirstName: "Alice", LastName: "Doe", TelephoneNumber: "123"},
			{FirstName: "Bob", LastName: "Doe", TelephoneNumber: "321"},
			{FirstName: "Xin", LastName: "Smith", TelephoneNumber: "893"},
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

		err, dbCustomers := customerRepository.ListBy("", "Do")

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
