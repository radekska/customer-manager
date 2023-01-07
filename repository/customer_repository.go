package repository

import "customer-manager/database"

type CustomerRepository interface {
	Create(customer *database.Customer)
	AddPurchase(customer *database.Customer, purchase *database.Purchase)
	AddRepair(customer *database.Customer, repair *database.Repair)
}
