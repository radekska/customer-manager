package server

import (
	"bytes"
	"customer-manager/database"
	"customer-manager/repositories"
	"encoding/json"
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
	err, customerToDelete := s.GetByID(customerID)
	if err != nil {
		return err
	}
	for index, customer := range s.customers {
		if customer.ID == customerToDelete.ID {
			s.customers = append(s.customers[:index], s.customers[index+1:]...)
			return nil
		}
	}
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
	return &repositories.CustomerNotFoundError{CustomerID: customerID}, nil
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

func makeRequest(t *testing.T, method string, path string, body io.Reader) *http.Request {
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
		req := makeRequest(t, http.MethodGet, "/api/customers", nil)

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
		req := makeRequest(t, http.MethodGet, "/api/customers", nil)

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assert.ElementsMatch(t, expectedCustomers, decodeCustomers(t, resp.Body))
	})

	t.Run("test create new customer", func(t *testing.T) {
		server.repository = &StubCustomerRepository{}
		body, _ := json.Marshal(&customer)
		req := makeRequest(t, http.MethodPost, "/api/customers", bytes.NewBuffer(body))

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		_, currentCustomers := server.repository.GetAll()
		assert.ElementsMatch(t, []database.Customer{customer}, currentCustomers)
	})

	t.Run("test create new customer invalid payload", func(t *testing.T) {
		server.repository = &StubCustomerRepository{}
		body, _ := json.Marshal(map[string]string{"invalid": "invalid"}) // TODO test unique tel. number
		req := makeRequest(t, http.MethodPost, "/api/customers", bytes.NewBuffer(body))

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
		req := makeRequest(t, http.MethodPost, "/api/customers", nil)
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
		req := makeRequest(t, http.MethodGet, "/api/customers/8a5cae65-222c-4164-a08b-9983af7e366c", nil)

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
		req := makeRequest(t, http.MethodGet, fmt.Sprintf("/api/customers/%s", invalidCustomerID), nil)

		resp := getResponse(t, server, req)

		assertNotFoundResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf("customer with given id '%s' does not exists", invalidCustomerID),
		})
	})

	t.Run("test get customer by invalid id", func(t *testing.T) {
		server.repository = &StubCustomerRepository{}
		invalidCustomerID := "im-not-uuid"
		req := makeRequest(t, http.MethodGet, fmt.Sprintf("/api/customers/%s", invalidCustomerID), nil)

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
		req := makeRequest(
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
		req := makeRequest(
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
		req := makeRequest(
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
		req := makeRequest(t, http.MethodPut, "/api/customers/8a5cae65-222c-4164-a08b-9983af7e366c", nil)
		req.Header.Set("Content-Type", invalidContentType)

		resp := getResponse(t, server, req)

		assertBadRequestResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf(
				"invalid content-type header specified: '%s', allowed: 'application/json'",
				invalidContentType,
			),
		})
	})

	t.Run("test delete customer and it's relations", func(t *testing.T) {
		customerOneID := "8a5cae65-222c-4164-a08b-9983af7e366c"
		customerTwo := database.Customer{
			ID:              "aba3a29f-f8c4-4173-acc0-b01a2f18c1bb",
			FirstName:       "John",
			LastName:        "Brown",
			TelephoneNumber: "12345",
		}
		server.repository = &StubCustomerRepository{
			customers: []database.Customer{
				{
					ID:              customerOneID,
					FirstName:       "Bob",
					LastName:        "Toe",
					TelephoneNumber: "367654567",
					Purchases: []database.Purchase{
						{
							ID:         "d11aeae2-d18a-4b6f-8ed5-2223c015adfd",
							FrameModel: "Solano",
							CustomerID: customerOneID,
						},
					},
					Repairs: []database.Repair{
						{
							ID:          "b483c02c-d4d0-4da9-8601-50a72c1eac14",
							Description: "Repair 1",
							Cost:        123.34,
							CustomerID:  customerOneID,
						},
					},
				},
				customerTwo,
			},
		}
		req := makeRequest(t, http.MethodDelete, fmt.Sprintf("/api/customers/%s", customerOneID), bytes.NewBuffer([]byte{}))

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		err, customers := server.repository.GetAll()
		assert.NoError(t, err)
		assert.Equal(t, []database.Customer{customerTwo}, customers)
		assert.Equal(t, "", resp.Header.Get("Content-Type"))
	})

	t.Run("test delete customer but not found", func(t *testing.T) {
		invalidID := "37567fea-71ab-4677-9b19-708370034a66"
		server.repository = &StubCustomerRepository{}
		req := makeRequest(
			t,
			http.MethodDelete,
			fmt.Sprintf("/api/customers/%s", invalidID),
			bytes.NewBuffer([]byte{}),
		)

		resp := getResponse(t, server, req)

		assertNotFoundResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf("customer with given id '%s' does not exists", invalidID),
		})
	})

	t.Run("test delete customer but invalid id", func(t *testing.T) {
		invalidID := "invalid-id"
		server.repository = &StubCustomerRepository{}
		req := makeRequest(
			t,
			http.MethodDelete,
			fmt.Sprintf("/api/customers/%s", invalidID),
			bytes.NewBuffer([]byte{}),
		)

		resp := getResponse(t, server, req)

		assertBadRequestResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", invalidID),
		})
	})

}
