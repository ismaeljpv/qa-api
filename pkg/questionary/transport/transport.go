package transport

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
)

func RequestError(e string) error {
	return errors.New(e)
}

type (
	GenericRequest struct{}

	FindQuestionByIDRequest struct {
		ID string `json:"id"`
	}
	FindQuestionByIDResponse struct {
		QuestionInfo domain.QuestionInfo `json:"questionInfo"`
	}

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

	DeleteQuestionRequest struct {
		ID string `json:"ID"`
	}

	DeleteQuestionResponse struct {
		Message string `json:"message"`
	}
)

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

func DecodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req GenericRequest
	return req, nil
}

func DecodeFindQuestionByIDRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, ok := mux.Vars(r)["id"]

	if !ok {
		return nil, RequestError("Please pass the ID of the question")
	}
	return FindQuestionByIDRequest{ID: id}, nil
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

	return req, nil
}

func DecodeDeleteQuestionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, ok := mux.Vars(r)["id"]

	if !ok {
		return nil, RequestError("Please pass the ID of the question")
	}
	return DeleteQuestionRequest{ID: id}, nil
}
