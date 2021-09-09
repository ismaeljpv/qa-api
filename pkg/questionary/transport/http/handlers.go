package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
	"github.com/ismaeljpv/qa-api/pkg/questionary/transport"
	httpError "github.com/ismaeljpv/qa-api/pkg/questionary/transport/http/error"
)

//
//This is the decode/encode handlers that will decode the request and encode the response returned by the API in HTTP protocol
//

func DecodeIDParamRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	id, ok := mux.Vars(r)["id"]
	if !ok {
		return nil, httpError.NewClientError(errors.New("Question ID is required"),
			http.StatusBadRequest,
			"Question ID is required")
	}
	return transport.IDParamRequest{ID: id}, nil
}

func DecodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req transport.GenericRequest
	return req, nil
}

func DecodeFindQuestionByUserRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	userId, ok := mux.Vars(r)["userId"]

	if !ok {
		return nil, httpError.NewClientError(errors.New("User ID is required"),
			http.StatusBadRequest,
			"User ID is required")
	}
	return transport.FindQuestionsByUserRequest{UserID: userId}, nil
}

func DecodeCreateQuestionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body domain.Question
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	valErr := transport.ValidateStruct(&body)
	if valErr != nil {
		return nil, httpError.NewClientError(valErr,
			http.StatusBadRequest,
			valErr.Error(),
		)
	}

	return body, nil
}

func DecodeAddAnswerRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var body domain.Answer
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		return nil, err
	}

	valErr := transport.ValidateStruct(&body)
	if valErr != nil {
		return nil, httpError.NewClientError(valErr,
			http.StatusBadRequest,
			valErr.Error(),
		)
	}

	return body, nil
}

func DecodeUpdateQuestionRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req transport.UpdateQuestionRequest
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
	valErr := transport.ValidateStruct(&req.QuestionInfo)
	if valErr != nil {
		return nil, httpError.NewClientError(valErr,
			http.StatusBadRequest,
			valErr.Error(),
		)
	}

	return req, nil
}

func EncodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.WriteHeader(http.StatusOK)
	return json.NewEncoder(w).Encode(response)
}

func HTTPErrorHandler(ctx context.Context, err error, w http.ResponseWriter) {

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
