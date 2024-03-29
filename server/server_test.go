package server

import (
	"bytes"
	"customer-manager/database"
	"customer-manager/repositories"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

type StubCustomerRepository struct {
	customerIDToCreate string
	customers          []database.Customer
}

func (s *StubCustomerRepository) Create(customer *database.Customer) (error, *database.Customer) {
	for _, c := range s.customers {
		if c.TelephoneNumber == customer.TelephoneNumber {
			return &repositories.DuplicatedTelephoneNumberError{Customer: customer}, nil
		}
	}
	customer.ID = s.customerIDToCreate
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

func (s *StubCustomerRepository) ListBy(
	firstName string,
	lastName string,
	limit int,
	offset int,
) (error, []database.Customer, int) {
	total := len(s.customers)
	if limit > len(s.customers) {
		return nil, s.customers, total
	}
	customers := s.customers[offset : offset+limit]
	return nil, customers, total
}

func (s *StubCustomerRepository) GetByID(customerID string) (error, *database.Customer) {
	for _, customer := range s.customers {
		if customer.ID == customerID {
			return nil, &customer
		}
	}
	return &repositories.CustomerNotFoundError{CustomerID: customerID}, nil
}

func (s *StubCustomerRepository) Update(
	customerDetails *database.Customer,
) (error, *database.Customer) {
	err, customer := s.GetByID(customerDetails.ID)
	if err != nil {
		return err, nil
	}
	customer.FirstName = customerDetails.FirstName
	customer.LastName = customerDetails.LastName
	customer.TelephoneNumber = customerDetails.TelephoneNumber
	return nil, customer
}

type StubPurchaseRepository struct {
	purchaseIDToCreate string
	purchases          []database.Purchase
}

func (s *StubPurchaseRepository) Create(
	customer *database.Customer,
	purchase *database.Purchase,
) (error, *database.Purchase) {
	purchase.CustomerID = customer.ID
	if purchase.ID == "" {
		purchase.ID = s.purchaseIDToCreate
	}
	customer.Purchases = append(customer.Purchases, *purchase)
	s.purchases = append(s.purchases, *purchase)
	return nil, purchase
}

func (s *StubPurchaseRepository) GetAll(customerID string) (error, []database.Purchase) {
	var customerPurchases []database.Purchase
	for _, purchase := range s.purchases {
		if purchase.CustomerID == customerID {
			customerPurchases = append(customerPurchases, purchase)
		}
	}
	return nil, customerPurchases
}

func (s *StubPurchaseRepository) DeleteByID(purchaseID string) error {
	for idx, purchase := range s.purchases {
		if purchase.ID == purchaseID {
			s.purchases = append(s.purchases[:idx], s.purchases[idx+1:]...)
			return nil
		}
	}
	return &repositories.PurchaseNotFoundError{PurchaseID: purchaseID}
}

func (s *StubPurchaseRepository) Update(purchase *database.Purchase) (error, *database.Purchase) {
	err, purchases := s.GetAll(purchase.CustomerID)
	if err != nil {
		return err, nil
	}
	for idx, currentPurchase := range purchases {
		if currentPurchase.ID == purchase.ID {
			purchases[idx] = *purchase
			return nil, purchase
		}
	}
	return &repositories.PurchaseNotFoundError{PurchaseID: purchase.ID}, nil
}

type StubRepairRepository struct {
	repairIDToCreate string
	repairs          []database.Repair
}

func (s *StubRepairRepository) Create(
	customer *database.Customer,
	repair *database.Repair,
) (error, *database.Repair) {
	repair.CustomerID = customer.ID
	if repair.ID == "" {
		repair.ID = s.repairIDToCreate
	}
	customer.Repairs = append(customer.Repairs, *repair)
	s.repairs = append(s.repairs, *repair)
	return nil, repair
}

func (s *StubRepairRepository) GetAll(customerID string) (error, []database.Repair) {
	var customerRepairs []database.Repair
	for _, repair := range s.repairs {
		if repair.CustomerID == customerID {
			customerRepairs = append(customerRepairs, repair)
		}
	}
	return nil, customerRepairs
}

func (s *StubRepairRepository) DeleteByID(repairID string) error {
	for idx, repair := range s.repairs {
		if repair.ID == repairID {
			s.repairs = append(s.repairs[:idx], s.repairs[idx+1:]...)
			return nil
		}
	}
	return &repositories.RepairNotFoundError{RepairID: repairID}
}

func getCustomer() database.Customer {
	return database.Customer{
		ID:              "ec8f6cb1-61f6-4dfc-b970-9dd81ff2547f",
		FirstName:       "John",
		LastName:        "Doe",
		TelephoneNumber: "123-456-789",
	}
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

func getResponse(
	t *testing.T,
	server *CustomerManagerServer,
	request *http.Request,
) *http.Response {
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

func assertCustomerDetailsResponse(
	t *testing.T,
	resp *http.Response,
	expectedCustomerDetails map[string]string,
) {
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

func assertBadRequestResponse(
	t *testing.T,
	resp *http.Response,
	expectedDetails map[string]string,
) {
	t.Helper()
	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	assertResponse(t, resp, expectedDetails)
}

func assertNotFoundResponse(t *testing.T, resp *http.Response, expectedDetails map[string]string) {
	t.Helper()
	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)
	assertResponse(t, resp, expectedDetails)
}

func TestCustomerHandlers(t *testing.T) {
	customer := database.Customer{
		FirstName:       "John",
		LastName:        "Doe",
		TelephoneNumber: "123-456-789",
	}
	server := NewCustomerManagerServer(
		fiber.New(),
		&StubCustomerRepository{},
		&StubPurchaseRepository{},
		&StubRepairRepository{},
	)

	t.Run("test get all customers", func(t *testing.T) {
		server.customerRepository = &StubCustomerRepository{customers: []database.Customer{
			{
				ID:              "7dd4ace2-d792-4532-bda2-c986a9a04363",
				FirstName:       "Jane",
				LastName:        "Doe",
				TelephoneNumber: "123567848",
			},
			{
				ID:              "8a5cae65-222c-4164-a08b-9983af7e366c",
				FirstName:       "Bob",
				LastName:        "Toe",
				TelephoneNumber: "367654567",
			},
		}}
		req := makeRequest(t, http.MethodGet, "/api/customers", nil)

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var actualCustomers fiber.Map
		err := json.NewDecoder(resp.Body).Decode(&actualCustomers)
		assert.NoError(t, err)
		assert.Equal(
			t,
			fiber.Map{
				"data": []interface{}{
					map[string]interface{}{
						"created_at":       "0001-01-01T00:00:00Z",
						"first_name":       "Jane",
						"id":               "7dd4ace2-d792-4532-bda2-c986a9a04363",
						"last_name":        "Doe",
						"telephone_number": "123567848",
						"updated_at":       "0001-01-01T00:00:00Z",
					},
					map[string]interface{}{
						"created_at":       "0001-01-01T00:00:00Z",
						"first_name":       "Bob",
						"id":               "8a5cae65-222c-4164-a08b-9983af7e366c",
						"last_name":        "Toe",
						"telephone_number": "367654567",
						"updated_at":       "0001-01-01T00:00:00Z",
					},
				},
				"total": 2.0,
			},
			actualCustomers,
		)
	})

	t.Run("test get paginated customers", func(t *testing.T) {
		server.customerRepository = &StubCustomerRepository{customers: []database.Customer{
			{
				ID:              "7dd4ace2-d792-4532-bda2-c986a9a04363",
				FirstName:       "Jane",
				LastName:        "Doe",
				TelephoneNumber: "123567848",
			},
			{
				ID:              "8a5cae65-222c-4164-a08b-9983af7e366c",
				FirstName:       "Bob",
				LastName:        "Toe",
				TelephoneNumber: "367654567",
			},
			{
				ID:              "325cae65-222c-4164-a08b-9983af7e366c",
				FirstName:       "Joe",
				LastName:        "Doe",
				TelephoneNumber: "567231123",
			},
		}}
		req := makeRequest(t, http.MethodGet, "/api/customers?limit=2&offset=1", nil)

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var actualCustomers fiber.Map
		err := json.NewDecoder(resp.Body).Decode(&actualCustomers)
		assert.NoError(t, err)
		assert.Equal(t, fiber.Map{
			"data": []interface{}{
				map[string]interface{}{
					"created_at":       "0001-01-01T00:00:00Z",
					"first_name":       "Bob",
					"id":               "8a5cae65-222c-4164-a08b-9983af7e366c",
					"last_name":        "Toe",
					"telephone_number": "367654567",
					"updated_at":       "0001-01-01T00:00:00Z",
				},
				map[string]interface{}{
					"created_at":       "0001-01-01T00:00:00Z",
					"first_name":       "Joe",
					"id":               "325cae65-222c-4164-a08b-9983af7e366c",
					"last_name":        "Doe",
					"telephone_number": "567231123",
					"updated_at":       "0001-01-01T00:00:00Z",
				},
			},
			"total": 3.0,
		}, actualCustomers)
	})

	t.Run("test get no customers as empty repository", func(t *testing.T) {
		var expectedCustomers []database.Customer
		server.customerRepository = &StubCustomerRepository{customers: expectedCustomers}
		req := makeRequest(t, http.MethodGet, "/api/customers", nil)

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		var actualCustomers fiber.Map
		err := json.NewDecoder(resp.Body).Decode(&actualCustomers)
		assert.NoError(t, err)
		assert.Equal(t, fiber.Map{
			"data":  nil,
			"total": 0.0,
		}, actualCustomers)
	})

	t.Run("test create new customer", func(t *testing.T) {
		server.customerRepository = &StubCustomerRepository{
			customerIDToCreate: "67a85348-2afe-4677-99ce-ed7cdc17e525",
		}
		body, _ := json.Marshal(&customer)
		req := makeRequest(t, http.MethodPost, "/api/customers", bytes.NewBuffer(body))

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		assertCustomerDetailsResponse(t, resp, map[string]string{
			"id":               "67a85348-2afe-4677-99ce-ed7cdc17e525",
			"first_name":       "John",
			"last_name":        "Doe",
			"telephone_number": "123-456-789",
			"created_at":       "0001-01-01T00:00:00Z",
			"updated_at":       "0001-01-01T00:00:00Z",
		})
		_, currentCustomers, total := server.customerRepository.ListBy("", "", 10, 0)
		customer.ID = "67a85348-2afe-4677-99ce-ed7cdc17e525"
		assert.ElementsMatch(t, []database.Customer{customer}, currentCustomers)
		assert.Equal(t, 1, total)
	})

	t.Run("test create new customer duplicated telephone number", func(t *testing.T) {
		server.customerRepository = &StubCustomerRepository{
			customers: []database.Customer{{FirstName: "Bob", LastName: "Smith", TelephoneNumber: "123-456-789"}},
		}
		body, _ := json.Marshal(&customer)
		req := makeRequest(t, http.MethodPost, "/api/customers", bytes.NewBuffer(body))

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
		assertCustomerDetailsResponse(t, resp, map[string]string{
			"detail": "customer 'John Doe' cannot have telephone number '123-456-789' as already taken.",
		})
		_, _, total := server.customerRepository.ListBy("", "", 10, 0)
		assert.Equal(t, 1, total)
	})

	t.Run("test create new customer invalid payload", func(t *testing.T) {
		server.customerRepository = &StubCustomerRepository{}
		body, _ := json.Marshal(
			map[string]string{"invalid": "invalid"},
		) // TODO test unique tel. number
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
		_, currentCustomers, total := server.customerRepository.ListBy("", "", 10, 0)
		assert.ElementsMatch(t, []database.Customer{}, currentCustomers)
		assert.Equal(t, 0, total)
	})

	t.Run("test create new customer invalid content-type header", func(t *testing.T) {
		invalidContentType := "text/html"
		server.customerRepository = &StubCustomerRepository{}
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
		server.customerRepository = &StubCustomerRepository{
			customers: []database.Customer{
				{
					ID:              "7dd4ace2-d792-4532-bda2-c986a9a04363",
					FirstName:       "Jane",
					LastName:        "Doe",
					TelephoneNumber: "123567848",
				},
				{
					ID:              "8a5cae65-222c-4164-a08b-9983af7e366c",
					FirstName:       "Bob",
					LastName:        "Toe",
					TelephoneNumber: "367654567",
				},
			},
		}
		req := makeRequest(
			t,
			http.MethodGet,
			"/api/customers/8a5cae65-222c-4164-a08b-9983af7e366c",
			nil,
		)

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
		server.customerRepository = &StubCustomerRepository{
			customers: []database.Customer{
				{
					ID:              "7dd4ace2-d792-4532-bda2-c986a9a04363",
					FirstName:       "Jane",
					LastName:        "Doe",
					TelephoneNumber: "123567848",
				},
			},
		}
		invalidCustomerID := "8a5cae65-222c-4164-a08b-9983af7e366c"
		req := makeRequest(
			t,
			http.MethodGet,
			fmt.Sprintf("/api/customers/%s", invalidCustomerID),
			nil,
		)

		resp := getResponse(t, server, req)

		assertNotFoundResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf("customer with given id '%s' does not exists", invalidCustomerID),
		})
	})

	t.Run("test get customer by invalid id", func(t *testing.T) {
		server.customerRepository = &StubCustomerRepository{}
		invalidCustomerID := "im-not-uuid"
		req := makeRequest(
			t,
			http.MethodGet,
			fmt.Sprintf("/api/customers/%s", invalidCustomerID),
			nil,
		)

		resp := getResponse(t, server, req)

		assertBadRequestResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf("given customer id '%s' is not a valid UUID", invalidCustomerID),
		})
	})

	t.Run("test edit customer details", func(t *testing.T) {
		server.customerRepository = &StubCustomerRepository{customers: []database.Customer{
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
		server.customerRepository = &StubCustomerRepository{}
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
		server.customerRepository = &StubCustomerRepository{
			customers: []database.Customer{
				{
					ID:              "7dd4ace2-d792-4532-bda2-c986a9a04363",
					FirstName:       "Jane",
					LastName:        "Doe",
					TelephoneNumber: "123567848",
				},
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
		server.customerRepository = &StubCustomerRepository{
			customers: []database.Customer{
				{
					ID:              "8a5cae65-222c-4164-a08b-9983af7e366c",
					FirstName:       "Bob",
					LastName:        "Toe",
					TelephoneNumber: "367654567",
				},
			},
		}
		req := makeRequest(
			t,
			http.MethodPut,
			"/api/customers/8a5cae65-222c-4164-a08b-9983af7e366c",
			nil,
		)
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
		server.customerRepository = &StubCustomerRepository{
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
		req := makeRequest(
			t,
			http.MethodDelete,
			fmt.Sprintf("/api/customers/%s", customerOneID),
			bytes.NewBuffer([]byte{}),
		)

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		err, customers, total := server.customerRepository.ListBy("", "", 10, 0)
		assert.NoError(t, err)
		assert.Equal(t, []database.Customer{customerTwo}, customers)
		assert.Equal(t, "", resp.Header.Get("Content-Type"))
		assert.Equal(t, 1, total)
	})

	t.Run("test delete customer but not found", func(t *testing.T) {
		invalidID := "37567fea-71ab-4677-9b19-708370034a66"
		server.customerRepository = &StubCustomerRepository{}
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
		server.customerRepository = &StubCustomerRepository{}
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

func TestPurchaseHandlers(t *testing.T) {
	customer := getCustomer()
	t.Run("test get all purchases", func(t *testing.T) {
		server := NewCustomerManagerServer(
			fiber.New(),
			&StubCustomerRepository{},
			&StubPurchaseRepository{},
			&StubRepairRepository{},
		)
		err, _ := server.purchasesRepository.Create(
			&customer,
			&database.Purchase{
				ID:           "ca1224cb-c993-4d45-8053-73c56aaf2c77",
				FrameModel:   "Model1",
				LensType:     "Lens1",
				LensPower:    "Power1",
				PD:           "PD1",
				PurchaseType: "PurchaseType1",
				PurchasedAt:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		)
		assert.NoError(t, err)
		err, _ = server.purchasesRepository.Create(
			&customer,
			&database.Purchase{
				ID:           "5b521e40-e0f1-47fd-a832-fe6ea3fba22c",
				FrameModel:   "Model2",
				LensType:     "Lens2",
				LensPower:    "Power2",
				PD:           "PD2",
				PurchaseType: "PurchaseType2",
				PurchasedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		)
		assert.NoError(t, err)
		req := makeRequest(
			t,
			http.MethodGet,
			"/api/customers/ec8f6cb1-61f6-4dfc-b970-9dd81ff2547f/purchases",
			nil,
		)

		resp := getResponse(t, server, req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var actualPurchases []map[string]string
		err = json.NewDecoder(resp.Body).Decode(&actualPurchases)
		assert.NoError(t, err)
		assert.Equal(
			t,
			[]map[string]string{
				{
					"created_at":    "0001-01-01T00:00:00Z",
					"customer_id":   "ec8f6cb1-61f6-4dfc-b970-9dd81ff2547f",
					"frame_model":   "Model1",
					"id":            "ca1224cb-c993-4d45-8053-73c56aaf2c77",
					"lens_power":    "Power1",
					"lens_type":     "Lens1",
					"pd":            "PD1",
					"purchase_type": "PurchaseType1",
					"purchased_at":  "2022-01-01T00:00:00Z",
					"updated_at":    "0001-01-01T00:00:00Z",
				},
				{
					"created_at":    "0001-01-01T00:00:00Z",
					"customer_id":   "ec8f6cb1-61f6-4dfc-b970-9dd81ff2547f",
					"frame_model":   "Model2",
					"id":            "5b521e40-e0f1-47fd-a832-fe6ea3fba22c",
					"lens_power":    "Power2",
					"lens_type":     "Lens2",
					"pd":            "PD2",
					"purchase_type": "PurchaseType2",
					"purchased_at":  "2021-01-01T00:00:00Z",
					"updated_at":    "0001-01-01T00:00:00Z",
				},
			},
			actualPurchases,
		)
	})

	t.Run("test create purchase for a customer", func(t *testing.T) {
		server := NewCustomerManagerServer(fiber.New(), &StubCustomerRepository{
			customers: []database.Customer{customer},
		}, &StubPurchaseRepository{purchaseIDToCreate: "80dfb090-deea-4672-873d-a9cf8d4103e0"}, &StubRepairRepository{})
		req := makeRequest(
			t,
			http.MethodPost,
			fmt.Sprintf("/api/customers/%s/purchases", customer.ID),
			bytes.NewBuffer([]byte(`{
				"frame_model": "Model1",
				"lens_type": "Lens1",
				"lens_power": "Power1",
				"pd": "PD1",
				"purchase_type": "PurchaseType1",
				"purchased_at": "2021-01-01"
			}`)),
		)

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusCreated, resp.StatusCode)
		var actualPurchase map[string]string
		err := json.NewDecoder(resp.Body).Decode(&actualPurchase)
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]string{
				"created_at":    "0001-01-01T00:00:00Z",
				"customer_id":   "ec8f6cb1-61f6-4dfc-b970-9dd81ff2547f",
				"frame_model":   "Model1",
				"id":            "80dfb090-deea-4672-873d-a9cf8d4103e0",
				"lens_power":    "Power1",
				"lens_type":     "Lens1",
				"pd":            "PD1",
				"purchase_type": "PurchaseType1",
				"purchased_at":  "2021-01-01T00:00:00Z",
				"updated_at":    "0001-01-01T00:00:00Z",
			},
			actualPurchase,
		)
	})

	t.Run("test update purchase for a customer", func(t *testing.T) {
		purchases := []database.Purchase{{
			ID:           "ca1224cb-c993-4d45-8053-73c56aaf2c77",
			FrameModel:   "Model1",
			LensType:     "Lens1",
			LensPower:    "Power1",
			PD:           "PD1",
			PurchaseType: "PurchaseType1",
			PurchasedAt:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			CustomerID:   customer.ID,
		}, {
			ID:           "5b521e40-e0f1-47fd-a832-fe6ea3fba22c",
			FrameModel:   "Model2",
			LensType:     "Lens2",
			LensPower:    "Power2",
			PD:           "PD2",
			PurchaseType: "PurchaseType2",
			PurchasedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			CustomerID:   customer.ID,
		}}
		server := NewCustomerManagerServer(fiber.New(), &StubCustomerRepository{
			customers: []database.Customer{customer},
		}, &StubPurchaseRepository{purchases: purchases}, &StubRepairRepository{})

		body, _ := json.Marshal(
			map[string]string{
				"frame_model":   "UpdatedModel1",
				"lens_type":     "UpdatedLens1",
				"lens_power":    "UpdatedPower1",
				"pd":            "UpdatedPD1",
				"purchase_type": "UpdatedPurchaseType1",
				"purchased_at":  "2025-01-01",
			},
		)
		req := makeRequest(
			t,
			http.MethodPut,
			fmt.Sprintf(
				"/api/customers/%s/purchases/%s",
				customer.ID,
				"ca1224cb-c993-4d45-8053-73c56aaf2c77",
			),
			bytes.NewBuffer(body),
		)

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		assertResponse(t, resp, map[string]string{
			"created_at":    "0001-01-01T00:00:00Z",
			"customer_id":   customer.ID,
			"frame_model":   "UpdatedModel1",
			"id":            "ca1224cb-c993-4d45-8053-73c56aaf2c77",
			"lens_power":    "UpdatedPower1",
			"lens_type":     "UpdatedLens1",
			"pd":            "UpdatedPD1",
			"purchase_type": "UpdatedPurchaseType1",
			"purchased_at":  "2025-01-01T00:00:00Z",
			"updated_at":    "0001-01-01T00:00:00Z",
		})
	})

	t.Run("test delete purchase for a customer", func(t *testing.T) {
		customerID := "ec8f6cb1-61f6-4dfc-b970-9dd81ff2547f"
		purchases := []database.Purchase{{
			ID:           "ca1224cb-c993-4d45-8053-73c56aaf2c77",
			FrameModel:   "Model1",
			LensType:     "Lens1",
			LensPower:    "Power1",
			PD:           "PD1",
			PurchaseType: "PurchaseType1",
			PurchasedAt:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			CustomerID:   customerID,
		}, {
			ID:           "5b521e40-e0f1-47fd-a832-fe6ea3fba22c",
			FrameModel:   "Model2",
			LensType:     "Lens2",
			LensPower:    "Power2",
			PD:           "PD2",
			PurchaseType: "PurchaseType2",
			PurchasedAt:  time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC),
			CustomerID:   customerID,
		}}
		server := NewCustomerManagerServer(fiber.New(), &StubCustomerRepository{
			customers: []database.Customer{customer},
		}, &StubPurchaseRepository{purchases: purchases}, &StubRepairRepository{})
		req := makeRequest(
			t,
			http.MethodDelete,
			fmt.Sprintf(
				"/api/customers/%s/purchases/%s",
				customerID,
				"ca1224cb-c993-4d45-8053-73c56aaf2c77",
			),
			nil,
		)

		resp := getResponse(t, server, req)

		assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)
		err, currentPurchases := server.purchasesRepository.GetAll(customerID)
		assert.NoError(t, err)
		assert.Equal(t, []database.Purchase{purchases[1]}, currentPurchases)
	})

	t.Run("test delete purchase but not found", func(t *testing.T) {
		invalidID := "37567fea-71ab-4677-9b19-708370034a66"
		server := NewCustomerManagerServer(
			fiber.New(),
			&StubCustomerRepository{},
			&StubPurchaseRepository{},
			&StubRepairRepository{},
		)
		req := makeRequest(
			t,
			http.MethodDelete,
			fmt.Sprintf(
				"/api/customers/33c2cb49-6156-4efe-b282-b0ba553d883f/purchases/%s",
				invalidID,
			),
			bytes.NewBuffer([]byte{}),
		)

		resp := getResponse(t, server, req)

		assertNotFoundResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf("purchase with given id '%s' does not exists", invalidID),
		})
	})

	t.Run("test delete purchase for a customer with invalid purchase id", func(t *testing.T) {
		invalidID := "invalid-id"
		server := NewCustomerManagerServer(
			fiber.New(),
			&StubCustomerRepository{},
			&StubPurchaseRepository{},
			&StubRepairRepository{},
		)
		req := makeRequest(
			t,
			http.MethodDelete,
			fmt.Sprintf(
				"/api/customers/33c2cb49-6156-4efe-b282-b0ba553d883f/purchases/%s",
				invalidID,
			),
			bytes.NewBuffer([]byte{}),
		)

		resp := getResponse(t, server, req)

		assertBadRequestResponse(t, resp, map[string]string{
			"detail": fmt.Sprintf("given purchase id '%s' is not a valid UUID", invalidID),
		})
	})
}

func TestRepairHandler(t *testing.T) {
	customer := getCustomer()

	t.Run("test create repair for custoemr", func(t *testing.T) {
		server := NewCustomerManagerServer(
			fiber.New(),
			&StubCustomerRepository{customers: []database.Customer{customer}},
			&StubPurchaseRepository{},
			&StubRepairRepository{},
		)
		req := makeRequest(t,
			http.MethodPost,
			fmt.Sprintf("/api/customers/%s/repairs", customer.ID),
			bytes.NewBuffer([]byte(`{
				"description": "repair I",
        "cost": "1.50",
        "reported_at": "2021-01-01"
			}`)),
		)

		resp := getResponse(t, server, req)

		assert.Equal(t, http.StatusCreated, resp.StatusCode)
		var createdRepair map[string]any
		err := json.NewDecoder(resp.Body).Decode(&createdRepair)
		assert.NoError(t, err)
		assert.Equal(
			t,
			map[string]any{
				"cost":        1.5,
				"created_at":  "0001-01-01T00:00:00Z",
				"customer_id": "ec8f6cb1-61f6-4dfc-b970-9dd81ff2547f",
				"description": "repair I",
				"id":          "",
				"reported_at": "2021-01-01T00:00:00Z",
			},
			createdRepair,
		)
	})

	t.Run("test list repairs for a customer", func(t *testing.T) {
		server := NewCustomerManagerServer(
			fiber.New(),
			&StubCustomerRepository{},
			&StubPurchaseRepository{},
			&StubRepairRepository{},
		)
		err, _ := server.repairsRepository.Create(
			&customer,
			&database.Repair{
				ID:          "ca1224cb-c993-4d45-8053-73c56aaf2c77",
				Description: "To be repaired",
				Cost:        12.65,
				CreatedAt:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
				ReportedAt:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		)
		assert.NoError(t, err)
		err, _ = server.repairsRepository.Create(
			&customer,
			&database.Repair{
				ID:          "5b521e40-e0f1-47fd-a832-fe6ea3fba22c",
				Description: "To be repaired II",
				Cost:        2.65,
				CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
				ReportedAt:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		)

		assert.NoError(t, err)
		req := makeRequest(
			t,
			http.MethodGet,
			fmt.Sprintf("/api/customers/%s/repairs", customer.ID),
			nil,
		)

		resp := getResponse(t, server, req)

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var actualRepairs []map[string]any
		err = json.NewDecoder(resp.Body).Decode(&actualRepairs)
		assert.NoError(t, err)
		assert.Equal(
			t,
			[]map[string]any{
				{
					"cost":        12.65,
					"created_at":  "2022-01-01T00:00:00Z",
					"reported_at": "2022-01-01T00:00:00Z",
					"customer_id": "ec8f6cb1-61f6-4dfc-b970-9dd81ff2547f",
					"description": "To be repaired",
					"id":          "ca1224cb-c993-4d45-8053-73c56aaf2c77",
				},
				{
					"cost":        2.65,
					"created_at":  "2020-01-01T00:00:00Z",
					"reported_at": "2022-01-01T00:00:00Z",
					"customer_id": "ec8f6cb1-61f6-4dfc-b970-9dd81ff2547f",
					"description": "To be repaired II",
					"id":          "5b521e40-e0f1-47fd-a832-fe6ea3fba22c",
				},
			},
			actualRepairs,
		)
	})

	t.Run("test delete repair for customer", func(t *testing.T) {
		server := NewCustomerManagerServer(
			fiber.New(),
			&StubCustomerRepository{customers: []database.Customer{customer}},
			&StubPurchaseRepository{},
			&StubRepairRepository{},
		)
		repairOne := &database.Repair{
			ID:          "ca1224cb-c993-4d45-8053-73c56aaf2c77",
			Description: "To be repaired",
			Cost:        12.65,
			CreatedAt:   time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
			ReportedAt:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		}
		repairTwo := &database.Repair{
			ID:          "5b521e40-e0f1-47fd-a832-fe6ea3fba22c",
			Description: "To be repaired II",
			Cost:        2.65,
			CreatedAt:   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
			ReportedAt:  time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC),
		}

		err, _ := server.repairsRepository.Create(&customer, repairOne)
		assert.NoError(t, err)
		err, _ = server.repairsRepository.Create(&customer, repairTwo)
		assert.NoError(t, err)

		req := makeRequest(
			t,
			http.MethodDelete,
			fmt.Sprintf("/api/customers/%s/repairs/%s", customer.ID, repairOne.ID),
			nil,
		)

		resp := getResponse(t, server, req)

		assert.Equal(t, http.StatusNoContent, resp.StatusCode)
		err, currentRepairs := server.repairsRepository.GetAll(customer.ID)
		assert.NoError(t, err)
		assert.Len(t, currentRepairs, 1)
		assert.Equal(t, []database.Repair{*repairTwo}, currentRepairs)
	})
}
