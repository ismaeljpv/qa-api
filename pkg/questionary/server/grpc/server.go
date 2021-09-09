package grpc

import (
	"context"
	"errors"

	"github.com/go-kit/kit/log"

	"github.com/go-kit/kit/transport/grpc"
	transport "github.com/ismaeljpv/qa-api/pkg/questionary/transport/grpc"
	pb "github.com/ismaeljpv/qa-api/pkg/questionary/transport/grpc/protobuff"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

// This is the gRPC server configuration and initialization layer

type gRPCServer struct {
	findAll    grpc.Handler
	findByID   grpc.Handler
	findByUser grpc.Handler
	create     grpc.Handler
	addAnswer  grpc.Handler
	update     grpc.Handler
	delete     grpc.Handler
	pb.UnimplementedQuestionaryServiceServer
}

func NewGRPCServer(endpoints transport.Endpoints, logger log.Logger) pb.QuestionaryServiceServer {

	return &gRPCServer{
		findAll: grpc.NewServer(
			endpoints.FindAllQuestions,
			transport.DecodeRequest,
			transport.EncodeGetQuestionsResponse,
		),
		findByID: grpc.NewServer(
			endpoints.FindQuestionById,
			transport.DecodeIDParamRequest,
			transport.EncodeQuestionInfoResponse,
		),
		findByUser: grpc.NewServer(
			endpoints.FindQuestionsByUser,
			transport.DecodeFindQuestionByUserRequest,
			transport.EncodeGetQuestionsResponse,
		),
		create: grpc.NewServer(
			endpoints.CreateQuestion,
			transport.DecodeCreateQuestionRequest,
			transport.EncodeQuestionResponse,
		),
		addAnswer: grpc.NewServer(
			endpoints.AddAnswer,
			transport.DecodeAddAnswerRequest,
			transport.EncodeQuestionInfoResponse,
		),
		update: grpc.NewServer(
			endpoints.UpdateQuestion,
			transport.DecodeUpdateQuestionRequest,
			transport.EncodeQuestionInfoResponse,
		),
		delete: grpc.NewServer(
			endpoints.DeleteQuestion,
			transport.DecodeIDParamRequest,
			transport.EncodeGenericMessageResponse,
		),
	}
}

func (server *gRPCServer) FindAll(ctx context.Context, msg *pb.EmptyMessage) (*pb.Questions, error) {
	_, resp, err := server.findAll.ServeGRPC(ctx, msg)
	if err != nil {
		return &pb.Questions{}, err
	}

	questions, ok := resp.(*pb.Questions)
	if !ok {
		return &pb.Questions{}, errors.New("Error parsing the response for FindAll() method")
	}
	return questions, nil
}

func (server *gRPCServer) FindByID(ctx context.Context, id *wrapperspb.StringValue) (*pb.QuestionInfo, error) {
	_, resp, err := server.findByID.ServeGRPC(ctx, id)
	if err != nil {
		return &pb.QuestionInfo{}, err
	}

	questionInfo, ok := resp.(*pb.QuestionInfo)
	if !ok {
		return &pb.QuestionInfo{}, errors.New("Error parsing the response for FindByID() method")
	}
	return questionInfo, nil
}

func (server *gRPCServer) FindByUser(ctx context.Context, id *wrapperspb.StringValue) (*pb.Questions, error) {
	_, resp, err := server.findByUser.ServeGRPC(ctx, id)
	if err != nil {
		return &pb.Questions{}, err
	}

	questions, ok := resp.(*pb.Questions)
	if !ok {
		return &pb.Questions{}, errors.New("Error parsing the response for FindByUser() method")
	}
	return questions, nil
}

func (server *gRPCServer) Create(ctx context.Context, question *pb.Question) (*pb.Question, error) {
	_, resp, err := server.create.ServeGRPC(ctx, question)
	if err != nil {
		return &pb.Question{}, err
	}

	newQuestion, ok := resp.(*pb.Question)
	if !ok {
		return &pb.Question{}, errors.New("Error parsing the response for Create() method")
	}
	return newQuestion, nil
}

func (server *gRPCServer) AddAnswer(ctx context.Context, answer *pb.Answer) (*pb.QuestionInfo, error) {
	_, resp, err := server.addAnswer.ServeGRPC(ctx, answer)
	if err != nil {
		return &pb.QuestionInfo{}, err
	}

	newAnswer, ok := resp.(*pb.QuestionInfo)
	if !ok {
		return &pb.QuestionInfo{}, errors.New("Error parsing the response for AddAnswer() method")
	}
	return newAnswer, nil
}

func (server *gRPCServer) Update(ctx context.Context, questionInfo *pb.QuestionUpdate) (*pb.QuestionInfo, error) {
	_, resp, err := server.update.ServeGRPC(ctx, questionInfo)
	if err != nil {
		return &pb.QuestionInfo{}, err
	}

	updatedInfo, ok := resp.(*pb.QuestionInfo)
	if !ok {
		return &pb.QuestionInfo{}, errors.New("Error parsing the response for Update() method")
	}
	return updatedInfo, nil
}

func (server *gRPCServer) Delete(ctx context.Context, id *wrapperspb.StringValue) (*pb.GenericMessage, error) {
	_, resp, err := server.delete.ServeGRPC(ctx, id)
	if err != nil {
		return &pb.GenericMessage{}, err
	}

	message, ok := resp.(*pb.GenericMessage)
	if !ok {
		return &pb.GenericMessage{}, errors.New("Error parsing the response for Delete() method")
	}
	return message, nil
}
