package server

type editCustomerDetailsRequest struct {
	FirstName       string `json:"first_name"       validate:"required"`
	LastName        string `json:"last_name"        validate:"required"`
	TelephoneNumber string `json:"telephone_number" validate:"required"`
}

type createCustomerRequest = editCustomerDetailsRequest
