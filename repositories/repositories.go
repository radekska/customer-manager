package repositories

import "customer-manager/database"

type CustomerRepository interface {
	Create(customer *database.Customer) (error, *database.Customer)
	DeleteByID(customerID string) error
}

type PurchaseRepository interface {
	Create(customer *database.Customer, purchase *database.Purchase) error
	DeleteByID(purchaseID string) error
}

type RepairRepository interface {
	Create(customer *database.Customer, repair *database.Repair) error
}
