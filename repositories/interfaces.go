package repositories

import "customer-manager/database"

type CustomerRepository interface {
	Create(customer *database.Customer) (error, *database.Customer)
	DeleteByID(customerID string) error
	ListBy(customerFirstName string, customerLastName string) (error, []database.Customer)
	GetByID(customerID string) (error, *database.Customer)
	Update(customer *database.Customer) (error, *database.Customer)
}

type PurchaseRepository interface {
	Create(customer *database.Customer, purchase *database.Purchase) (error, *database.Purchase)
	GetAll(customerID string) (error, []database.Purchase)
	DeleteByID(purchaseID string) error
	Update(customer *database.Purchase) (error, *database.Purchase)
}

type RepairRepository interface {
	Create(customer *database.Customer, repair *database.Repair) (error, *database.Repair)
  GetAll(customerID string) (error, []database.Repair)
	DeleteByID(repairID string) error
}
