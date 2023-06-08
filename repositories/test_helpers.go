package repositories

import (
	"customer-manager/database"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
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
	return &database.Purchase{
		FrameModel:   "Model1",
		LensType:     "LensType1",
		LensPower:    "LensPower",
		PD:           "CustomPD",
		PurchaseType: "CustomPurchaseType",
		PurchasedAt:  time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
	}
}
