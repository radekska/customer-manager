package repository

type Customer struct {
	FirstName string
	LastName  string
	Age       int
}

type CustomerRepository interface {
	Create(customer *Customer)
}
