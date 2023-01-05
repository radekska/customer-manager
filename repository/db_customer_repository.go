package repository

import (
	"customer-manager/database"
	"gorm.io/gorm"
)

type DBCustomerRepository struct {
	db *gorm.DB
}

func (d *DBCustomerRepository) Create(customer *Customer) {
	d.db.Create(&database.Customer{FirstName: customer.FirstName, LastName: customer.LastName, Age: customer.Age})
}
