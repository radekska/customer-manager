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

func (s *StubCustomerRepository) Update(customerDetails *database.Customer) (error, *database.Customer) {
	err, customer := s.GetByID(customerDetails.ID)
	if err != nil {
		return err, nil
	}
	customer.FirstName = customerDetails.FirstName
	customer.LastName = customerDetails.LastName
	customer.TelephoneNumber = customerDetails.TelephoneNumber
	return nil, customer
}

func createRequest(t *testing.T, method string, path string, body io.Reader) *http.Request {
	t.Helper()
	req, err := http.NewRequest(method, path, body)
	if err != nil {
		t.Fatal("failed during request creation.", err)
	}
	if method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch {
		req.Header.Set("Content-Type", "application/json")
	}
	return req
}

func getResponse(t *testing.T, server *CustomerManagerServer, request *http.Request) *http.Response {
	resp, err := server.App.Test(request)
	if err != nil {
		t.Fatal("failed to get a response.", err)
	}
	return resp
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

func assertCustomerDetailsResponse(t *testing.T, resp *http.Response, expectedCustomerDetails map[string]string) {
	t.Helper()
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	actualCustomerDetails := make(map[string]string)
	err := json.NewDecoder(resp.Body).Decode(&actualCustomerDetails)
	assert.NoError(t, err)
	assert.Equal(t, expectedCustomerDetails, actualCustomerDetails)
}

func assertResponse(t *testing.T, resp *http.Response, expectedDetails map[string]string) {
	t.Helper()
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
	actualDetails := make(map[string]string)
	err := json.NewDecoder(resp.Body).Decode(&actualDetails)
	assert.NoError(t, err)
	assert.Equal(t, expectedDetails, actualDetails)
}

func assertBadRequestResponse(t *testing.T, resp *http.Response, expectedDetails map[string]string) {
	t.Helper()
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	assertResponse(t, resp, expectedDetails)
}

func assertNotFoundResponse(t *testing.T, resp *http.Response, expectedDetails map[string]string) {
	t.Helper()
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	assertResponse(t, resp, expectedDetails)
}

func TestCustomerManagerServer(t *testing.T) {
	customer := database.Customer{FirstName: "John", LastName: "Doe", TelephoneNumber: "123-456-789"}
	server := NewCustomerManagerServer(fiber.New(), &StubCustomerRepository{})

	t.Run("test get all customers", func(t *testing.T) {
		server.repository = &StubCustomerRepository{customers: []database.Customer{
			{ID: "7dd4ace2-d792-4532-bda2-c986a9a04363", FirstName: "Jane", LastName: "Doe", TelephoneNumber: "123567848"},
			{ID: "8a5cae65-222c-4164-a08b-9983af7e366c", FirstName: "Bob", LastName: "Toe", TelephoneNumber: "367654567"},
		}}
		req := createRequest(t, http.MethodGet, "/api/customers", nil)

		resp := getResponse(t, server, req)

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
		server.repository = &StubCustomerRepository{customers: expectedCustomers}
		req := createRequest(t, http.MethodGet, "/api/customers", nil)

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.ElementsMatch(t, expectedCustomers, decodeCustomers(t, resp.Body))
	})

	t.Run("test create new customer", func(t *testing.T) {
		server.repository = &StubCustomerRepository{}
		body, _ := json.Marshal(&customer)
		req := createRequest(t, http.MethodPost, "/api/customers", bytes.NewBuffer(body))

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		_, currentCustomers := server.repository.GetAll()
		assert.ElementsMatch(t, []database.Customer{customer}, currentCustomers)
	})

	t.Run("test create new customer invalid payload", func(t *testing.T) {
		server.repository = &StubCustomerRepository{}
		body, _ := json.Marshal(map[string]string{"invalid": "invalid"}) // TODO test unique tel. number
		req := createRequest(t, http.MethodPost, "/api/customers", bytes.NewBuffer(body))

		resp := getResponse(t, server, req)

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
		_, currentCustomers := server.repository.GetAll()
		assert.ElementsMatch(t, []database.Customer{}, currentCustomers)
	})

	t.Run("test create new customer invalid content-type header", func(t *testing.T) {
		invalidContentType := "text/html"
		server.repository = &StubCustomerRepository{}
		req := createRequest(t, http.MethodPost, "/api/customers", nil)
		req.Header.Set("Content-Type", invalidContentType)

		resp := getResponse(t, server, req)

		assertBadRequestResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf(
				"invalid content-type header specified: '%s', allowed: 'application/json'",
				invalidContentType,
			),
		})
	})

	t.Run("test get customer by its id", func(t *testing.T) {
		server.repository = &StubCustomerRepository{
			customers: []database.Customer{
				{ID: "7dd4ace2-d792-4532-bda2-c986a9a04363", FirstName: "Jane", LastName: "Doe", TelephoneNumber: "123567848"},
				{ID: "8a5cae65-222c-4164-a08b-9983af7e366c", FirstName: "Bob", LastName: "Toe", TelephoneNumber: "367654567"},
			},
		}
		req := createRequest(t, http.MethodGet, "/api/customers/8a5cae65-222c-4164-a08b-9983af7e366c", nil)

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assertCustomerDetailsResponse(t, resp, map[string]string{
			"id":               "8a5cae65-222c-4164-a08b-9983af7e366c",
			"first_name":       "Bob",
			"last_name":        "Toe",
			"telephone_number": "367654567",
			"created_at":       "0001-01-01T00:00:00Z",
			"updated_at":       "0001-01-01T00:00:00Z",
		})
	})

	t.Run("test get customer by its id but not found", func(t *testing.T) {
		server.repository = &StubCustomerRepository{
			customers: []database.Customer{
				{ID: "7dd4ace2-d792-4532-bda2-c986a9a04363", FirstName: "Jane", LastName: "Doe", TelephoneNumber: "123567848"},
			},
		}
		invalidCustomerID := "8a5cae65-222c-4164-a08b-9983af7e366c"
		req := createRequest(t, http.MethodGet, fmt.Sprintf("/api/customers/%s", invalidCustomerID), nil)

		resp := getResponse(t, server, req)

		assertNotFoundResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf("customer with given id '%s' does not exists", invalidCustomerID),
		})
	})

	t.Run("test get customer by invalid id", func(t *testing.T) {
		server.repository = &StubCustomerRepository{}
		invalidCustomerID := "im-not-uuid"
		req := createRequest(t, http.MethodGet, fmt.Sprintf("/api/customers/%s", invalidCustomerID), nil)

		resp := getResponse(t, server, req)

		assertBadRequestResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", invalidCustomerID),
		})
	})

	t.Run("test edit customer details", func(t *testing.T) {
		server.repository = &StubCustomerRepository{customers: []database.Customer{
			{
				ID:              "8a5cae65-222c-4164-a08b-9983af7e366c",
				FirstName:       "Bob",
				LastName:        "Toe",
				TelephoneNumber: "367654567",
			},
		}}

		body, _ := json.Marshal(
			map[string]string{
				"first_name":       "John",
				"last_name":        "Doe",
				"telephone_number": "123456891",
			},
		)
		req := createRequest(
			t,
			http.MethodPut,
			fmt.Sprintf("/api/customers/%s", "8a5cae65-222c-4164-a08b-9983af7e366c"),
			bytes.NewBuffer(body),
		)

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assertCustomerDetailsResponse(t, resp, map[string]string{
			"id":               "8a5cae65-222c-4164-a08b-9983af7e366c",
			"first_name":       "John",
			"last_name":        "Doe",
			"telephone_number": "123456891",
			"created_at":       "0001-01-01T00:00:00Z",
			"updated_at":       "0001-01-01T00:00:00Z",
		})
	})

	t.Run("test edit customer details invalid id", func(t *testing.T) {
		invalidCustomerID := "im-not-uuid"
		server.repository = &StubCustomerRepository{}
		req := createRequest(
			t,
			http.MethodPut,
			fmt.Sprintf("/api/customers/%s", invalidCustomerID),
			bytes.NewBuffer([]byte{}),
		)

		resp := getResponse(t, server, req)

		assertBadRequestResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", invalidCustomerID),
		})
	})

	t.Run("test edit customer details but not found", func(t *testing.T) {
		invalidCustomerID := "5936ca64-3c2c-4ada-89f9-27fece73a0a8"
		server.repository = &StubCustomerRepository{
			customers: []database.Customer{
				{ID: "7dd4ace2-d792-4532-bda2-c986a9a04363", FirstName: "Jane", LastName: "Doe", TelephoneNumber: "123567848"},
			},
		}
		body, err := json.Marshal(
			map[string]string{
				"first_name":       "John",
				"last_name":        "Doe",
				"telephone_number": "123456891",
			},
		)
		assert.NoError(t, err)
		req := createRequest(
			t,
			http.MethodPut,
			fmt.Sprintf("/api/customers/%s", invalidCustomerID),
			bytes.NewBuffer(body),
		)

		resp := getResponse(t, server, req)

		assertNotFoundResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf("customer with given id '%s' does not exists", invalidCustomerID),
		})
	})

	t.Run("test  edit customer details invalid content-type header", func(t *testing.T) {
		invalidContentType := "text/html"
		server.repository = &StubCustomerRepository{
			customers: []database.Customer{
				{ID: "8a5cae65-222c-4164-a08b-9983af7e366c", FirstName: "Bob", LastName: "Toe", TelephoneNumber: "367654567"},
			},
		}
		req := createRequest(t, http.MethodPut, "/api/customers/8a5cae65-222c-4164-a08b-9983af7e366c", nil)
		req.Header.Set("Content-Type", invalidContentType)

		resp := getResponse(t, server, req)

		assertBadRequestResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf(
				"invalid content-type header specified: '%s', allowed: 'application/json'",
				invalidContentType,
			),
		})
	})

}
