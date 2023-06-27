package server

import (
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
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

var valid *validator.Validate = validator.New()

type CreateRepairRequest struct {
	Description string `json:"description" validate:"required"`
	Cost        string `json:"cost" validate:"required"`
	ReportedAt  string `json:"reported_at"  validate:"required,datetime=2006-01-02"`
}

func registerValidators() {
	valid.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
	valid.RegisterStructValidation(func(sl validator.StructLevel) {}, CreateRepairRequest{})
}

func validateRequest(r interface{}) map[string]string {
	err := valid.Struct(r)
	if err == nil {
		return nil
	}

	validationErrors := make(map[string]string)

	for _, err := range err.(validator.ValidationErrors) {
		validationErrors[err.Field()] = err.Error()
	}
	return validationErrors
}
