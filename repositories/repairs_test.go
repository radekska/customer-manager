package repositories

import (
	"customer-manager/database"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBRepairRepository(t *testing.T) {
	customerRepository := DBCustomerRepository{db}
	repairRepository := DBRepairRepository{db}
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
