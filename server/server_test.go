package server

import (
	"bytes"
	"customer-manager/database"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

type StubCustomerRepository struct {
	customers []database.Customer
}

func (s *StubCustomerRepository) Create(customer *database.Customer) (error, *database.Customer) {
	s.customers = append(s.customers, *customer)
	return nil, customer
}

func (s *StubCustomerRepository) DeleteByID(customerID string) error {
	return nil
}

func (s *StubCustomerRepository) GetAll() (error, []database.Customer) {
	return nil, s.customers
}

func decodeCustomers(t *testing.T, body io.Reader) []database.Customer {
	t.Helper()
	var currentCustomers []database.Customer
	err := json.NewDecoder(body).Decode(&currentCustomers)
	if err != nil {
		t.Fatal(err)
	}
	return currentCustomers
}

func TestCustomerManagerServer(t *testing.T) {
	customer := database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123-456-789"}

	t.Run("test get all customers", func(t *testing.T) {
		expectedCustomers := []database.Customer{customer}
		server := NewCustomerManagerServer(fiber.New(), &StubCustomerRepository{customers: expectedCustomers})
		req, _ := http.NewRequest(http.MethodGet, "/api/customers", nil)

		resp, _ := server.App.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.ElementsMatch(t, expectedCustomers, decodeCustomers(t, resp.Body))
	})

	t.Run("test get no customers as empty repository", func(t *testing.T) {
		var expectedCustomers []database.Customer
		server := NewCustomerManagerServer(fiber.New(), &StubCustomerRepository{customers: expectedCustomers})
		req, _ := http.NewRequest(http.MethodGet, "/api/customers", nil)

		resp, _ := server.App.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.ElementsMatch(t, expectedCustomers, decodeCustomers(t, resp.Body))
	})

	t.Run("test create new customer", func(t *testing.T) {
		repository := &StubCustomerRepository{}
		server := NewCustomerManagerServer(fiber.New(), repository)
		body, _ := json.Marshal(&customer)
		req, _ := http.NewRequest(http.MethodPost, "/api/customers", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := server.App.Test(req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		_, currentCustomers := repository.GetAll()
		assert.ElementsMatch(t, []database.Customer{customer}, currentCustomers)
	})
}
