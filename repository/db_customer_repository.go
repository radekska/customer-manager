package repository

import (
	"customer-manager/database"
	"gorm.io/gorm"
)

type DBCustomerRepository struct {
	db *gorm.DB
}

func (d *DBCustomerRepository) Create(customer *database.Customer) {
	d.db.Create(&customer)
}

func (d *DBCustomerRepository) AddPurchase(customer *database.Customer, purchase *database.Purchase) {
	err := d.db.Model(customer).Association("Purchases").Append(purchase)
	if err != nil {
		panic(err)
	}
}

func (d *DBCustomerRepository) AddRepair(customer *database.Customer, repair *database.Repair) {
	err := d.db.Model(customer).Association("Repairs").Append(repair)
	if err != nil {
		panic(err)
	}
}
