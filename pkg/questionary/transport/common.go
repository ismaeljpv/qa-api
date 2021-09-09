package transport

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
)

var validate *validator.Validate = validator.New()

type (
	GenericRequest struct{}

	FindQuestionsByUserRequest struct {
		UserID string `json:"userId"`
	}

	UpdateQuestionRequest struct {
		ID           string              `json:"ID"`
		QuestionInfo domain.QuestionInfo `json:"questionInfo"`
	}

	IDParamRequest struct {
		ID string `json:"ID"`
	}

	GenericMessageResponse struct {
		Message string `json:"message"`
		Status  string `json:"status"`
		Code    int64  `json:"code"`
	}
)

func ValidateStruct(s interface{}) error {
	errs := validate.Struct(s)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			return errors.New(fmt.Sprintf("Error on field %v, data is %v", err.Field(), err.ActualTag()))
		}
	}
	return nil
}
