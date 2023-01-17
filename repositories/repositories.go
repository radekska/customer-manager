package repositories

import "customer-manager/database"

type CustomerRepository interface {
	Create(customer *database.Customer) (error, *database.Customer)
	DeleteByID(customerID string) error
	GetAll() (error, []database.Customer)
	GetByID(customerID string) (error, *database.Customer)
}

type PurchaseRepository interface {
	Create(customer *database.Customer, purchase *database.Purchase) (error, *database.Purchase)
	DeleteByID(purchaseID string) error
}

type RepairRepository interface {
	Create(customer *database.Customer, repair *database.Repair) (error, *database.Repair)
	DeleteByID(repairID string) error
}
