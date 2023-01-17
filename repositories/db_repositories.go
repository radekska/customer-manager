package repositories

import (
	"customer-manager/database"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBCustomerRepository struct {
	DB *gorm.DB
}

func (d *DBCustomerRepository) Create(customer *database.Customer) (error, *database.Customer) {
	return d.DB.Create(&customer).Error, customer
}

func (d *DBCustomerRepository) DeleteByID(customerID string) error {
	return d.DB.Select(clause.Associations).Delete(&database.Customer{ID: customerID}).Error
}

func (d *DBCustomerRepository) GetAll() (error, []database.Customer) {
	var customers []database.Customer
	result := d.DB.Find(&customers)
	return result.Error, customers
}

func (d *DBCustomerRepository) GetByID(customerID string) (error, *database.Customer) {
	var customer database.Customer
	result := d.DB.Where("id = ?", customerID).First(&customer)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return result.Error, nil
	}
	return result.Error, &customer
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

type DBRepairRepository struct {
	db *gorm.DB
}

func (d *DBRepairRepository) Create(customer *database.Customer, repair *database.Repair) (error, *database.Repair) {
	return d.db.Model(customer).Association("Repairs").Append(repair), repair
}

func (d *DBRepairRepository) DeleteByID(repairID string) error {
	return d.db.Delete(&database.Repair{ID: repairID}).Error
}
