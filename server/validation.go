package server

import "time"

type EditCustomerDetailsRequest struct {
	FirstName       string `json:"first_name"       validate:"required"`
	LastName        string `json:"last_name"        validate:"required"`
	TelephoneNumber string `json:"telephone_number" validate:"required"`
}

type CreateCustomerRequest = EditCustomerDetailsRequest

type CreatePurchaseRequest struct {
	FrameModel   string    `json:"frame_model"   validate:"required"`
	LensType     string    `json:"lens_type"     validate:"required"`
	LensPower    string    `json:"lens_power"    validate:"required"`
	PD           string    `json:"pd"            validate:"required"`
	PurchaseType string    `json:"purchase_type" validate:"required"`
	PurchasedAt  time.Time `json:"purchased_at"  validate:"required"`
}
