package repositories

import (
	"customer-manager/database"
	"gorm.io/gorm"
)

type DBCustomerRepository struct {
	db *gorm.DB
}

type DBRepairRepository struct {
	db *gorm.DB
}
type DBPurchaseRepository struct {
	db *gorm.DB
}

func (d *DBCustomerRepository) Create(customer *database.Customer) error {
	return d.db.Create(&customer).Error
}

func (d *DBPurchaseRepository) Create(customer *database.Customer, purchase *database.Purchase) error {
	return d.db.Model(customer).Association("Purchases").Append(purchase)
}

func (d *DBRepairRepository) Create(customer *database.Customer, repair *database.Repair) error {
	return d.db.Model(customer).Association("Repairs").Append(repair)
}
