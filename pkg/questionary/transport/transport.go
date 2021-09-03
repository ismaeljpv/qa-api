package transport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
)

var validate *validator.Validate = validator.New()

func RequestError(e string) error {
	return errors.New(e)
}

func validateStruct(v *validator.Validate, s interface{}) error {
	errors := validate.Struct(s)
	if errors != nil {
		for _, err := range errors.(validator.ValidationErrors) {
			return RequestError(fmt.Sprintf("Error on field %v, data is %v", err.Field(), err.ActualTag()))
		}
	}
	return nil
}

type (
	GenericRequest struct{}

	FindQuestionsByUserRequest struct {
		UserID string `json:"userId"`
	}

	CreateQuestionRequest struct {
		Question domain.Question `json:"question"`
	}

	AddAnswerRequest struct {
		Answer domain.Answer `json:"answer"`
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

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func DecodeIDParamRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		return nil, RequestError("Please pass the ID of the question")
	}
	return IDParamRequest{ID: id}, nil
}

func DecodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GenericRequest
	return req, nil
}

func DecodeFindQuestionByUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	userId, ok := mux.Vars(r)["userId"]

	if !ok {
		return nil, RequestError("Please pass the User ID of the question")
	}
	return FindQuestionsByUserRequest{UserID: userId}, nil
}

func DecodeCreateQuestionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req CreateQuestionRequest
	var body domain.Question
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	req.Question = body
	valErr := validateStruct(validate, &req.Question)
	if valErr != nil {
		return nil, valErr
	}

	return req, nil
}

func DecodeAddAnswerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req AddAnswerRequest
	var ans domain.Answer

	err := json.NewDecoder(r.Body).Decode(&ans)
	if err != nil {
		return nil, err
	}

	req.Answer = ans
	valErr := validateStruct(validate, &req.Answer)
	if valErr != nil {
		return nil, valErr
	}

	return req, nil
}

func DecodeUpdateQuestionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req UpdateQuestionRequest
	var info domain.QuestionInfo

	quetionId, ok := mux.Vars(r)["id"]
	if !ok {
		return nil, RequestError("Please pass the ID of the question")
	}

	err := json.NewDecoder(r.Body).Decode(&info)
	if err != nil {
		return nil, err
	}

	req.ID = quetionId
	req.QuestionInfo = info
	valErr := validateStruct(validate, &req.QuestionInfo)
	if valErr != nil {
		return nil, valErr
	}

	return req, nil
}
