package grpc

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
	"github.com/ismaeljpv/qa-api/pkg/questionary/transport/grpc/protobuff"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

//This is the decode/encode handlers that will decode the request and encode the response returned by the API in gRPC protocol

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
)

func validateStruct(v *validator.Validate, s interface{}) error {
	errs := validate.Struct(s)
	if errs != nil {
		for _, err := range errs.(validator.ValidationErrors) {
			return errors.New(fmt.Sprintf("Error validating body, error => %v", err.Error()))
		}
	}
	return nil
}

func DecodeIDParamRequest(ctx context.Context, request interface{}) (interface{}, error) {
	id, ok := request.(*wrapperspb.StringValue)
	if !ok || id == nil {
		return nil, errors.New("Question ID is required")
	}
	return IDParamRequest{ID: id.GetValue()}, nil
}

func DecodeRequest(ctx context.Context, request interface{}) (interface{}, error) {
	var req GenericRequest
	return req, nil
}

func DecodeFindQuestionByUserRequest(ctx context.Context, request interface{}) (interface{}, error) {
	userId, ok := request.(*wrapperspb.StringValue)
	if !ok || userId == nil {
		return nil, errors.New("UserID ID is required")
	}
	return FindQuestionsByUserRequest{UserID: userId.GetValue()}, nil
}

func DecodeCreateQuestionRequest(ctx context.Context, request interface{}) (interface{}, error) {
	var newQuestion domain.Question
	body, ok := request.(*protobuff.Question)
	if !ok || body == nil {
		return nil, errors.New("No body found in the request")
	}

	newQuestion.Statement = body.GetStatement()
	newQuestion.UserID = body.GetUserID()

	valErr := validateStruct(validate, &newQuestion)
	if valErr != nil {
		return nil, valErr
	}

	return newQuestion, nil
}

func DecodeAddAnswerRequest(ctx context.Context, request interface{}) (interface{}, error) {
	var newAnswer domain.Answer
	body, ok := request.(*protobuff.Answer)
	if !ok || body == nil {
		return nil, errors.New("No body found in the request")
	}

	newAnswer.Answer = body.GetAnswer()
	newAnswer.UserID = body.GetUserID()
	newAnswer.QuestionID = body.GetQuestionID()

	valErr := validateStruct(validate, &newAnswer)
	if valErr != nil {
		return nil, valErr
	}

	return newAnswer, nil
}

func DecodeUpdateQuestionRequest(ctx context.Context, request interface{}) (interface{}, error) {
	var req UpdateQuestionRequest
	var info domain.QuestionInfo
	questionUpdate, ok := request.(*protobuff.QuestionUpdate)
	if !ok || questionUpdate == nil {
		return nil, errors.New("No body found in the request")
	}

	if questionUpdate.QuestionID == "" {
		return nil, errors.New("No Question ID passed")
	}

	info.Question.ID = questionUpdate.GetQuestionInfo().GetQuestion().GetID()
	info.Question.Statement = questionUpdate.GetQuestionInfo().GetQuestion().GetStatement()
	info.Question.CreatedOn = questionUpdate.GetQuestionInfo().GetQuestion().GetCreatedOn()
	info.Question.UserID = questionUpdate.GetQuestionInfo().GetQuestion().GetUserID()

	info.Answer.ID = questionUpdate.GetQuestionInfo().GetAnswer().GetID()
	info.Answer.Answer = questionUpdate.GetQuestionInfo().GetAnswer().GetAnswer()
	info.Answer.QuestionID = questionUpdate.GetQuestionInfo().GetAnswer().GetQuestionID()
	info.Answer.UserID = questionUpdate.GetQuestionInfo().GetAnswer().GetUserID()
	info.Answer.CreatedOn = questionUpdate.GetQuestionInfo().GetAnswer().GetCreatedOn()

	req.ID = questionUpdate.QuestionID
	req.QuestionInfo = info
	valErr := validateStruct(validate, &req.QuestionInfo)
	if valErr != nil {
		return nil, valErr
	}

	return req, nil
}

func EncodeGetQuestionsResponse(_ context.Context, response interface{}) (interface{}, error) {
	var result protobuff.Questions
	var resQuestions = make([]*protobuff.QuestionInfo, 0)
	questions, ok := response.([]domain.QuestionInfo)
	if !ok {
		return []*protobuff.QuestionInfo{}, errors.New("Error parsing the response for gRPC Questions message")
	}

	for _, question := range questions {
		var info protobuff.QuestionInfo
		info.Question = &protobuff.Question{}
		info.Answer = &protobuff.Answer{}

		info.Question.ID = question.Question.ID
		info.Question.Statement = question.Question.Statement
		info.Question.UserID = question.Question.UserID
		info.Question.CreatedOn = question.Question.CreatedOn

		info.Answer.ID = question.Answer.ID
		info.Answer.Answer = question.Answer.Answer
		info.Answer.QuestionID = question.Answer.QuestionID
		info.Answer.UserID = question.Answer.UserID

		resQuestions = append(resQuestions, &info)
	}

	result.Questions = resQuestions
	return &result, nil
}

func EncodeQuestionInfoResponse(_ context.Context, response interface{}) (interface{}, error) {
	var info protobuff.QuestionInfo
	question, ok := response.(domain.QuestionInfo)
	if !ok {
		return &protobuff.QuestionInfo{}, errors.New("Error parsing the response for gRPC QuestionInfo message")
	}

	info.Question = &protobuff.Question{}
	info.Answer = &protobuff.Answer{}

	info.Question.ID = question.Question.ID
	info.Question.Statement = question.Question.Statement
	info.Question.UserID = question.Question.UserID
	info.Question.CreatedOn = question.Question.CreatedOn

	info.Answer.ID = question.Answer.ID
	info.Answer.Answer = question.Answer.Answer
	info.Answer.QuestionID = question.Answer.QuestionID
	info.Answer.UserID = question.Answer.UserID

	return &info, nil
}

func EncodeQuestionResponse(_ context.Context, response interface{}) (interface{}, error) {
	var info *protobuff.Question
	question, ok := response.(domain.Question)
	if !ok {
		return &protobuff.Question{}, errors.New("Error parsing the response for gRPC Question message")
	}
	info = &protobuff.Question{}
	info.ID = question.ID
	info.Statement = question.Statement
	info.UserID = question.UserID
	info.CreatedOn = question.CreatedOn
	return info, nil
}

func EncodeGenericMessageResponse(_ context.Context, response interface{}) (interface{}, error) {
	var resp *protobuff.GenericMessage
	message, ok := response.(string)
	if !ok {
		return &protobuff.GenericMessage{}, errors.New("Error parsing the response for gRPC GenericMessage message")
	}
	resp = &protobuff.GenericMessage{}
	resp.Message = message
	return resp, nil
}
