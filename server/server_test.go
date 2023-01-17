package server

import (
	"bytes"
	"customer-manager/database"
	"encoding/json"
	"errors"
	"fmt"
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

func (s *StubCustomerRepository) GetByID(customerID string) (error, *database.Customer) {
	for _, customer := range s.customers {
		if customer.ID == customerID {
			return nil, &customer
		}
	}
	return errors.New("customer not found"), nil
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
		server := NewCustomerManagerServer(fiber.New(), &StubCustomerRepository{customers: []database.Customer{
			{ID: "7dd4ace2-d792-4532-bda2-c986a9a04363", FirstName: "Jane", LastName: "Doe", TelephoneNumber: "123567848"},
			{ID: "8a5cae65-222c-4164-a08b-9983af7e366c", FirstName: "Bob", LastName: "Toe", TelephoneNumber: "367654567"},
		}})
		req, _ := http.NewRequest(http.MethodGet, "/api/customers", nil)

		resp, _ := server.App.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var actualCustomers []map[string]string
		err := json.NewDecoder(resp.Body).Decode(&actualCustomers)
		assert.NoError(t, err)
		assert.Equal(t, []map[string]string{{
			"id":               "7dd4ace2-d792-4532-bda2-c986a9a04363",
			"first_name":       "Jane",
			"last_name":        "Doe",
			"telephone_number": "123567848",
			"created_at":       "0001-01-01T00:00:00Z",
			"updated_at":       "0001-01-01T00:00:00Z",
		}, {
			"id":               "8a5cae65-222c-4164-a08b-9983af7e366c",
			"first_name":       "Bob",
			"last_name":        "Toe",
			"telephone_number": "367654567",
			"created_at":       "0001-01-01T00:00:00Z",
			"updated_at":       "0001-01-01T00:00:00Z",
		},
		}, actualCustomers)
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

	t.Run("test create new customer invalid payload", func(t *testing.T) {
		repository := &StubCustomerRepository{}
		server := NewCustomerManagerServer(fiber.New(), repository)
		body, _ := json.Marshal(map[string]string{"invalid": "invalid"})
		req, _ := http.NewRequest(http.MethodPost, "/api/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := server.App.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")

		actualErrorMessage := make(map[string]map[string]string)
		err := json.NewDecoder(resp.Body).Decode(&actualErrorMessage)
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]map[string]string{
				"first_name":       {"required": "The 'first_name' is required"},
				"last_name":        {"required": "The 'last_name' is required"},
				"telephone_number": {"required": "The 'telephone_number' is required"},
			},
			actualErrorMessage,
		)

		_, currentCustomers := repository.GetAll()
		assert.ElementsMatch(t, []database.Customer{}, currentCustomers)
	})

	t.Run("test get customer by its id", func(t *testing.T) {
		repository := &StubCustomerRepository{
			customers: []database.Customer{
				{ID: "7dd4ace2-d792-4532-bda2-c986a9a04363", FirstName: "Jane", LastName: "Doe", TelephoneNumber: "123567848"},
				{ID: "8a5cae65-222c-4164-a08b-9983af7e366c", FirstName: "Bob", LastName: "Toe", TelephoneNumber: "367654567"},
			},
		}
		server := NewCustomerManagerServer(fiber.New(), repository)
		req, _ := http.NewRequest(http.MethodGet, "/api/customers/8a5cae65-222c-4164-a08b-9983af7e366c", nil)

		resp, _ := server.App.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")

		actualCustomer := make(map[string]string)
		err := json.NewDecoder(resp.Body).Decode(&actualCustomer)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{
			"id":               "8a5cae65-222c-4164-a08b-9983af7e366c",
			"first_name":       "Bob",
			"last_name":        "Toe",
			"telephone_number": "367654567",
			"created_at":       "0001-01-01T00:00:00Z",
			"updated_at":       "0001-01-01T00:00:00Z",
		}, actualCustomer)
	})
	t.Run("test get customer by its id but not found", func(t *testing.T) {
		repository := &StubCustomerRepository{
			customers: []database.Customer{
				{ID: "7dd4ace2-d792-4532-bda2-c986a9a04363", FirstName: "Jane", LastName: "Doe", TelephoneNumber: "123567848"},
			},
		}
		invalidCustomerID := "8a5cae65-222c-4164-a08b-9983af7e366c"
		server := NewCustomerManagerServer(fiber.New(), repository)
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/customers/%s", invalidCustomerID), nil)

		resp, _ := server.App.Test(req)

		assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")

		errorMessage := make(map[string]string)
		err := json.NewDecoder(resp.Body).Decode(&errorMessage)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{
			"detail": fmt.Sprintf("customer with given id '%s' does not exists", invalidCustomerID),
		}, errorMessage)
	})

	t.Run("test get customer by invalid id", func(t *testing.T) {
		repository := &StubCustomerRepository{}
		invalidCustomerID := "im-not-uuid"
		server := NewCustomerManagerServer(fiber.New(), repository)
		req, _ := http.NewRequest(http.MethodGet, fmt.Sprintf("/api/customers/%s", invalidCustomerID), nil)

		resp, _ := server.App.Test(req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assert.Equal(t, resp.Header.Get("Content-Type"), "application/json")

		errorMessage := make(map[string]string)
		err := json.NewDecoder(resp.Body).Decode(&errorMessage)
		assert.NoError(t, err)
		assert.Equal(t, map[string]string{
			"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", invalidCustomerID),
		}, errorMessage)
	})

}
