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
	httpError "github.com/ismaeljpv/qa-api/pkg/questionary/transport/error"
)

//This is the decode/encode handlers that will decode the request and encode the response returned by the API

type errorHandler func(http.ResponseWriter, *http.Request) error

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

func validateStruct(v *validator.Validate, s interface{}) error {
	errors := validate.Struct(s)
	if errors != nil {
		for _, err := range errors.(validator.ValidationErrors) {
			return httpError.NewClientError(err,
				http.StatusBadRequest,
				fmt.Sprintf("Error on field %v, data is %v", err.Field(), err.ActualTag()),
			)
		}
	}
	return nil
}

func DecodeIDParamRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		return nil, httpError.NewClientError(errors.New("Question ID is required"),
			http.StatusBadRequest,
			"Question ID is required")
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
		return nil, httpError.NewClientError(errors.New("User ID is required"),
			http.StatusBadRequest,
			"User ID is required")
	}
	return FindQuestionsByUserRequest{UserID: userId}, nil
}

func DecodeCreateQuestionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body domain.Question
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	valErr := validateStruct(validate, body)
	if valErr != nil {
		return nil, valErr
	}

	return body, nil
}

func DecodeAddAnswerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body domain.Answer
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	valErr := validateStruct(validate, &body)
	if valErr != nil {
		return nil, valErr
	}

	return body, nil
}

func DecodeUpdateQuestionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req UpdateQuestionRequest
	var info domain.QuestionInfo

	quetionId, ok := mux.Vars(r)["id"]
	if !ok {
		return nil, httpError.NewClientError(errors.New("Question ID is required"),
			http.StatusBadRequest,
			"Question ID is required")
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

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func ErrorHandler(ctx context.Context, err error, w http.ResponseWriter) {

	switch er := err.(type) {
	case httpError.ClientError:
		body, err := er.ResponseBody()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("There was an error procesing your request"))
			return
		}
		status, headers := er.ResponseHeaders()
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(status)
		w.Write(body)
	case httpError.InternalServerError:
		body, err := er.ResponseBody()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("There was an error procesing your request"))
			return
		}
		_, headers := er.ResponseHeaders()
		for k, v := range headers {
			w.Header().Set(k, v)
		}
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(body)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("There was an error procesing your request"))
	}
}
