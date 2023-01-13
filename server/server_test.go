package server

import (
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
	return nil, nil
}

func (s *StubCustomerRepository) DeleteByID(customerID string) error {
	return nil
}

func (s *StubCustomerRepository) GetAll() (error, []database.Customer) {
	return nil, s.customers
}

func decodeResponse(t *testing.T, body io.Reader, toStructure interface{}) {
	t.Helper()
	err := json.NewDecoder(body).Decode(&toStructure)
	if err != nil {
		t.Fatal(err)
	}
}

func TestCustomerManagerServer(t *testing.T) {
	t.Run("test get all customers", func(t *testing.T) {
		expectedCustomers := []database.Customer{{FirstName: "John", LastName: "Doe", TelephoneNumber: "123-456-789"}}
		server := NewCustomerManagerServer(fiber.New(), &StubCustomerRepository{customers: expectedCustomers})
		req, _ := http.NewRequest(http.MethodGet, "/api/customers", nil)

		resp, _ := server.App.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var currentCustomers []database.Customer
		decodeResponse(t, resp.Body, currentCustomers)
		assert.ElementsMatch(t, expectedCustomers, currentCustomers)
	})

	t.Run("test get no customers as empty repository", func(t *testing.T) {
		var expectedCustomers []database.Customer
		server := NewCustomerManagerServer(fiber.New(), &StubCustomerRepository{customers: expectedCustomers})
		req, _ := http.NewRequest(http.MethodGet, "/api/customers", nil)

		resp, _ := server.App.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var currentCustomers []database.Customer
		decodeResponse(t, resp.Body, currentCustomers)
		assert.ElementsMatch(t, expectedCustomers, currentCustomers)
	})
}
