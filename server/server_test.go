package server

import (
	"customer-manager/database"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

type StubCustomerRepository struct {
}

func (s *StubCustomerRepository) Create(customer *database.Customer) (error, *database.Customer) {
	return nil, nil
}

func (s *StubCustomerRepository) DeleteByID(customerID string) error {
	return nil
}

func (s *StubCustomerRepository) GetAll() (error, []database.Customer) {
	return nil, []database.Customer{{FirstName: "John", LastName: "Doe", TelephoneNumber: "123-456-789"}}
}

func TestCustomerManagerServer(t *testing.T) {
	t.Run("test get all customers", func(t *testing.T) {
		expectedCustomers := []database.Customer{{FirstName: "John", LastName: "Doe", TelephoneNumber: "123-456-789"}}
		server := NewCustomerManagerServer(fiber.New(), &StubCustomerRepository{})
		req, _ := http.NewRequest(http.MethodGet, "/api/customers", nil)

		resp, _ := server.App.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var currentCustomers []database.Customer
		json.NewDecoder(resp.Body).Decode(&currentCustomers)
		assert.ElementsMatch(t, expectedCustomers, currentCustomers)
	})
}
