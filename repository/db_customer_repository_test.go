package repository

import (
	"customer-manager/database"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"testing"
)

func TestCreateCustomer(t *testing.T) {
	customer := &Customer{FirstName: "John", LastName: "Doe", Age: 32}
	db := database.GetDatabase("../test.db", &gorm.Config{})
	storage := DBCustomerRepository{db}

	storage.Create(customer)

	var dbCustomer database.Customer
	db.Where("first_name = ? AND last_name = ? AND age = ?", customer.FirstName, customer.LastName, customer.Age).First(&dbCustomer)

	assert.Equal(t, customer.FirstName, dbCustomer.FirstName)
	assert.Equal(t, customer.LastName, dbCustomer.LastName)
	assert.Equal(t, customer.Age, dbCustomer.Age)
}
