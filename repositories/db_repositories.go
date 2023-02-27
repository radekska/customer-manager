package repositories

import (
	"customer-manager/database"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DBCustomerRepository struct {
	DB *gorm.DB
}

type CustomerNotFoundError struct {
	CustomerID string
}

func (c *CustomerNotFoundError) Error() string {
	return fmt.Sprintf("customer with ID '%s' does not exist", c.CustomerID)
}

func (d *DBCustomerRepository) Create(customer *database.Customer) (error, *database.Customer) {
	return d.DB.Create(&customer).Error, customer
}

func (d *DBCustomerRepository) DeleteByID(customerID string) error {
	result := d.DB.Select(clause.Associations).Delete(&database.Customer{ID: customerID})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return &CustomerNotFoundError{CustomerID: customerID}
	}
	return nil
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
		return &CustomerNotFoundError{CustomerID: customerID}, nil
	}
	return result.Error, &customer
}

func (d *DBCustomerRepository) Update(customer *database.Customer) (error, *database.Customer) {
	result := d.DB.Model(customer).Select("FirstName", "LastName", "TelephoneNumber").Updates(customer)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return &CustomerNotFoundError{CustomerID: customer.ID}, nil
	}
	return result.Error, customer
}

type DBPurchaseRepository struct {
	DB *gorm.DB
}

func (d *DBPurchaseRepository) Create(
	customer *database.Customer,
	purchase *database.Purchase,
) (error, *database.Purchase) {
	return d.DB.Model(customer).Association("Purchases").Append(purchase), purchase
}

func (d *DBPurchaseRepository) GetAll(customerID string) (error, []database.Purchase) {
	var purchases []database.Purchase
	result := d.DB.Where("customer_id = ?", customerID).Find(&purchases)
	return result.Error, purchases
}

func (d *DBPurchaseRepository) DeleteByID(purchaseID string) error {
	return d.DB.Delete(&database.Purchase{ID: purchaseID}).Error
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
