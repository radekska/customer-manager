package repositories

import (
	"customer-manager/database"
	"errors"
	"fmt"
	"strings"

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

func (d *DBCustomerRepository) ListBy(
	customerFirstName string,
	customerLastName string,
	limit int,
	offset int,
) (error, []database.Customer, int) {
	var customers []database.Customer
	var result *gorm.DB
	firstName := strings.ToLower(customerFirstName)
	lastName := strings.ToLower(customerLastName)
	if firstName == "" && lastName == "" {
		result = d.DB.Offset(offset).Limit(limit).Order("first_name asc").Find(&customers)
	} else {
		firstNameQuery := fmt.Sprintf("%%%s%%", firstName)
		lastNameQuery := fmt.Sprintf("%%%s%%", lastName)

		result = d.DB.Where("LOWER(first_name) LIKE ? AND LOWER(last_name) LIKE ?",
			firstNameQuery, lastNameQuery).Offset(offset).Limit(limit).Order("first_name asc").Find(&customers)
	}

	var total int64
	d.DB.Model(&database.Customer{}).Count(&total)
	return result.Error, customers, int(total)
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
