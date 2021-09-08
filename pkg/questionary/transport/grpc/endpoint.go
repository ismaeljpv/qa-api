package grpc

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
	"github.com/ismaeljpv/qa-api/pkg/questionary/service"
)

//This is the enpoint configuration to instantiate all avaliable routes for the gRPC server.
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
		if err != nil {
			return []domain.QuestionInfo{}, errors.New(err.Error())
		}
		return questions, nil
	}
}

func makeFindQuestionByIDEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(IDParamRequest)
		question, err := s.FindByID(ctx, req.ID)
		if err != nil {
			return domain.QuestionInfo{}, errors.New(err.Error())
		}
		return question, nil
	}
}

func makeFindQuestiosnByUserEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(FindQuestionsByUserRequest)
		questions, err := s.FindByUser(ctx, req.UserID)
		if err != nil {
			return []domain.QuestionInfo{}, errors.New(err.Error())
		}
		return questions, nil
	}
}

func makeCreateQuestionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		question := request.(domain.Question)
		question, err := s.Create(ctx, question)
		if err != nil {
			return domain.Question{}, errors.New(err.Error())
		}
		return question, nil
	}
}

func makeAddAnswerEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		answer := request.(domain.Answer)
		questionInfo, err := s.AddAnswer(ctx, answer)
		if err != nil {
			return domain.QuestionInfo{}, errors.New(err.Error())
		}
		return questionInfo, nil
	}
}

func makeUpdateQuestionEndPoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateQuestionRequest)
		questionInfo, err := s.Update(ctx, req.QuestionInfo, req.ID)
		if err != nil {
			return domain.QuestionInfo{}, errors.New(err.Error())
		}
		return questionInfo, nil
	}
}

func makeDeleteQuestionEndpoint(s service.Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(IDParamRequest)
		msg, err := s.Delete(ctx, req.ID)
		if err != nil {
			return "", errors.New(err.Error())
		}
		return msg, nil
	}
}
