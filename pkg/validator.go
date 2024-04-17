package pkg

import (
	"strings"

	"github.com/go-playground/validator/v10"
)

type (
	ErrorResponse struct {
		FailedField string
		Tag         string
		Value       interface{}
	}

	XValidator struct {
		validator *validator.Validate
	}
)

var validate = validator.New(validator.WithRequiredStructEnabled())

var Valtor = &XValidator{
	validator: validate,
}

func (v XValidator) Validate(data interface{}) error {

	if err := validate.Struct(data); err != nil {
		return err
	}

	return nil
}

func (v XValidator) TransErrorToResponse(err *error) []ErrorResponse {
	validationErrors := []ErrorResponse{}

	errs := (*err).(validator.ValidationErrors)

	for _, err := range errs {
		var elem ErrorResponse

		elem.FailedField = strings.ToLower(err.Field()) // Export struct field name
		elem.Tag = err.Tag()                            // Export struct tag
		elem.Value = err.Value()                        // Export field value

		validationErrors = append(validationErrors, elem)
	}

	return validationErrors
}
