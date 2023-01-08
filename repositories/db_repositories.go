package repositories

import (
	"customer-manager/database"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBCustomerRepository struct {
	db *gorm.DB
}

func (d *DBCustomerRepository) Create(customer *database.Customer) (error, *database.Customer) {
	return d.db.Create(&customer).Error, customer
}

func (d *DBCustomerRepository) DeleteByID(customerID string) error {
	return d.db.Select(clause.Associations).Delete(&database.Customer{ID: customerID}).Error
}

type DBRepairRepository struct {
	db *gorm.DB
}
type DBPurchaseRepository struct {
	db *gorm.DB
}

func (d *DBPurchaseRepository) Create(customer *database.Customer, purchase *database.Purchase) (error, *database.Purchase) {
	return d.db.Model(customer).Association("Purchases").Append(purchase), purchase
}

func (d *DBPurchaseRepository) DeleteByID(purchaseID string) error {
	return d.db.Delete(&database.Purchase{ID: purchaseID}).Error
}

func (d *DBRepairRepository) Create(customer *database.Customer, repair *database.Repair) error {
	return d.db.Model(customer).Association("Repairs").Append(repair)
}
