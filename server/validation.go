package server

import (
	"strings"
	"time"

	"github.com/gookit/validate"
)

type EditCustomerDetailsRequest struct {
	FirstName       string `json:"first_name"       validate:"required"`
	LastName        string `json:"last_name"        validate:"required"`
	TelephoneNumber string `json:"telephone_number" validate:"required"`
}

type CreateCustomerRequest = EditCustomerDetailsRequest

type Date time.Time

func (d *Date) UnmarshalJSON(b []byte) error {
	t, err := time.Parse("2006-01-02", strings.ReplaceAll(string(b), "\"", ""))
	if err != nil {
		return err
	}
	*d = Date(t)
	return nil
}

func (d *Date) MarshalJSON() ([]byte, error) {
	return []byte(time.Time(*d).Format("2006-01-02")), nil
}

type CreatePurchaseRequest struct {
	FrameModel   string `json:"frame_model"   validate:"required"`
	LensType     string `json:"lens_type"     validate:"required"`
	LensPower    string `json:"lens_power"    validate:"required"`
	PD           string `json:"pd"            validate:"required"`
	PurchaseType string `json:"purchase_type" validate:"required"`
	// TODO - when invalid date specified it returns field is required
	PurchasedAt Date `json:"purchased_at"  validate:"required"`
}

type EditPurchaseRequest = CreatePurchaseRequest

type CreateRepairRequest struct {
	Description string  `json:"description" validate:"required"`
	Cost        float64 `json:"cost" validate:"required"`
	// TODO - when invalid date specified it returns field is required
	ReportedAt Date `json:"reported_at"  validate:"required"`
}

func getValidator(s interface{}) *validate.Validation {
	validate.Config(func(opt *validate.GlobalOption) {
		opt.StopOnError = false
	})
	v := validate.New(s)
	v.AddMessages(map[string]string{
		"required": "The '{field}' is required",
	})
	return v
}
