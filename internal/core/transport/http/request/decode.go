package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/M1sterZag/Dont_Play_Separately/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func init() {
	if err := requestValidator.RegisterValidation("password", passwordValidator); err != nil {
		panic(err)
	}
}

type validatable interface {
	Validate() error
}

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v, %w", err, core_errors.ErrInvalidArgument)
	}

	var err error
	v, ok := dest.(validatable)
	if ok {
		err = v.Validate()
	} else {
		err = requestValidator.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf("validate request: %v, %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}

func passwordValidator(fl validator.FieldLevel) bool {
	pw := fl.Field().String()
	if len(pw) < 8 || len(pw) > 72 {
		return false
	}

	var hasUpper, hasLower, hasDigit, hasSpecial bool
	for _, ch := range pw {
		switch {
		case 'A' <= ch && ch <= 'Z':
			hasUpper = true
		case 'a' <= ch && ch <= 'z':
			hasLower = true
		case '0' <= ch && ch <= '9':
			hasDigit = true
		default:
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}
