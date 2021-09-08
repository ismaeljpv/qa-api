package http

import (
	"context"
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	transport "github.com/ismaeljpv/qa-api/pkg/questionary/transport/http"
)

//This is the HTTP Server that will handle all avaliable operations of the API
//The endpoints configuration is used to instantiate all routes.
func NewHTTPServer(ctx context.Context, endpoints transport.Endpoints) http.Handler {

	router := mux.NewRouter()
	router.Use(commonMiddleware)
	serverOpts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(transport.ErrorHandler),
	}

	router.Methods("GET").Path("/question").Handler(httptransport.NewServer(
		endpoints.FindAllQuestions,
		transport.DecodeRequest,
		transport.EncodeResponse,
		serverOpts...,
	))

	router.Methods("GET").Path("/question/{id}").Handler(httptransport.NewServer(
		endpoints.FindQuestionById,
		transport.DecodeIDParamRequest,
		transport.EncodeResponse,
		serverOpts...,
	))

	router.Methods("GET").Path("/question/user/{userId}").Handler(httptransport.NewServer(
		endpoints.FindQuestionsByUser,
		transport.DecodeFindQuestionByUserRequest,
		transport.EncodeResponse,
		serverOpts...,
	))

	router.Methods("POST").Path("/question").Handler(httptransport.NewServer(
		endpoints.CreateQuestion,
		transport.DecodeCreateQuestionRequest,
		transport.EncodeResponse,
		serverOpts...,
	))

	router.Methods("POST").Path("/question/answer").Handler(httptransport.NewServer(
		endpoints.AddAnswer,
		transport.DecodeAddAnswerRequest,
		transport.EncodeResponse,
		serverOpts...,
	))

	router.Methods("PUT").Path("/question/{id}").Handler(httptransport.NewServer(
		endpoints.UpdateQuestion,
		transport.DecodeUpdateQuestionRequest,
		transport.EncodeResponse,
		serverOpts...,
	))

	router.Methods("DELETE").Path("/question/{id}").Handler(httptransport.NewServer(
		endpoints.DeleteQuestion,
		transport.DecodeIDParamRequest,
		transport.EncodeResponse,
		serverOpts...,
	))

	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}
