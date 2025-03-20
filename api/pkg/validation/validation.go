package validation

import (
	"github.com/go-playground/validator/v10"

	"errors"
	"fmt"
	"regexp"
	"strings"
)

const tagTrimmed = "trimmed"

type (
	Validate          = validator.Validate
	ErrorHTTPResponse struct {
		Message string            `json:"message"`
		Fields  map[string]string `json:"fields"`
	}
)

var (
	NewValidator = validator.New
)

// ValidateTrimmedWhitespaces checks if string has white space or end of line char symbols at start and end
func ValidateTrimmedWhitespaces(fl validator.FieldLevel) bool {
	s := fl.Field().String()
	return s == strings.TrimSpace(s)
}

func RegisterValidationTagTrimmed(v *validator.Validate) error {
	err := v.RegisterValidation(tagTrimmed, ValidateTrimmedWhitespaces)
	return err
}

func FormatValidationErr(err error) ErrorHTTPResponse {
	response := ErrorHTTPResponse{
		Message: "validation failed",
		Fields:  make(map[string]string),
	}
	// todo: get clear understanding what here happens with validationErrors errorsInsert
	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, fieldError := range validationErrors {
			regx := regexp.MustCompile("([a-z])([A-Z])")
			fieldName := strings.ToLower(regx.ReplaceAllString(fieldError.Field(), "${1}_${2}"))
			switch fieldError.Tag() {
			case "required":
				response.Fields[fieldName] = fmt.Sprintf(
					"field '%s' is required",
					fieldName,
				)
			case "email":
				response.Fields[fieldName] = "invalid email format"
			case "max":
				response.Fields[fieldName] = fmt.Sprintf(
					"max length is %s",
					fieldError.Param(),
				)
			case "min":
				response.Fields[fieldName] = fmt.Sprintf(
					"minimum length is %s",
					fieldError.Param(),
				)
			case "oneof":
				response.Fields[fieldName] = "must be one of values"
			case tagTrimmed:
				response.Fields[fieldName] = "must not contain leading and trailing whitespace"
			default:
				response.Fields[fieldError.Field()] = fmt.Sprintf(
					"validation failed for field '%s'",
					fieldName,
				)
			}
		}
	}

	return response
}
