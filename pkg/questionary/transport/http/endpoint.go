package http

import (
	"context"
	"net/http"

	"github.com/go-kit/kit/endpoint"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
	"github.com/ismaeljpv/qa-api/pkg/questionary/service"
	"github.com/ismaeljpv/qa-api/pkg/questionary/transport"
)

//This is the enpoint configuration to instantiate all avaliable routes for the HTTP server.
//Its injects the service to handle the data recieved in the request and process the response.
type Endpoints struct {
	FindAllQuestions    endpoint.Endpoint
	FindQuestionById    endpoint.Endpoint
	FindQuestionsByUser endpoint.Endpoint
	CreateQuestion      endpoint.Endpoint
	AddAnswer           endpoint.Endpoint
	UpdateQuestion      endpoint.Endpoint
	DeleteQuestion      endpoint.Endpoint
}

func MakeEndpoints(s service.Service) Endpoints {
	return Endpoints{
		FindAllQuestions:    makeFindAllQuestionsEndpoint(s),
		FindQuestionById:    makeFindQuestionByIDEndpoint(s),
		FindQuestionsByUser: makeFindQuestiosnByUserEndpoint(s),
		CreateQuestion:      makeCreateQuestionEndpoint(s),
		AddAnswer:           makeAddAnswerEndpoint(s),
		UpdateQuestion:      makeUpdateQuestionEndPoint(s),
		DeleteQuestion:      makeDeleteQuestionEndpoint(s),
	}
}

func makeFindAllQuestionsEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		questions, err := s.FindAll(ctx)
		return questions, err
	}
}

func makeFindQuestionByIDEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.IDParamRequest)
		question, err := s.FindByID(ctx, req.ID)
		return question, err
	}
}

func makeFindQuestiosnByUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.FindQuestionsByUserRequest)
		questions, err := s.FindByUser(ctx, req.UserID)
		return questions, err
	}
}

func makeCreateQuestionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		question := request.(domain.Question)
		question, err := s.Create(ctx, question)
		return question, err
	}
}

func makeAddAnswerEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		answer := request.(domain.Answer)
		questionInfo, err := s.AddAnswer(ctx, answer)
		return questionInfo, err
	}
}

func makeUpdateQuestionEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.UpdateQuestionRequest)
		questionInfo, err := s.Update(ctx, req.QuestionInfo, req.ID)
		return questionInfo, err
	}
}

func makeDeleteQuestionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(transport.IDParamRequest)
		msg, err := s.Delete(ctx, req.ID)

		return transport.GenericMessageResponse{
			Message: msg,
			Status:  http.StatusText(http.StatusOK),
			Code:    http.StatusOK,
		}, err
	}
}
