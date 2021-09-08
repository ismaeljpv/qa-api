package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ismaeljpv/qa-api/pkg/questionary/transport/grpc/protobuff"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func main() {
	var conn *grpc.ClientConn
	conn, err := grpc.Dial(":50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Connection error %s", err)
	}
	defer conn.Close()

	var ctx context.Context
	ctx = context.Background()
	client := protobuff.NewQuestionaryServiceClient(conn)

	questions, err := client.FindAll(ctx, &protobuff.EmptyMessage{})
	if err != nil {
		fmt.Println(err.Error())
	}

	questionsByUser, err := client.FindByUser(ctx, wrapperspb.String("1"))
	if err != nil {
		fmt.Println(err.Error())
	}

	question, err := client.FindByID(ctx, wrapperspb.String("de8b8892-d64e-4651-9416-b0b7472b4737"))
	if err != nil {
		fmt.Println(err.Error())
	}

	newQuestion := &protobuff.Question{
		Statement: "is gRPC great?",
		UserID:    "3",
	}
	createdQuestion, err := client.Create(ctx, newQuestion)
	if err != nil {
		fmt.Println(err.Error())
	}

	newAnswer := &protobuff.Answer{
		Answer:     "gRPC is awesome!",
		UserID:     "33",
		QuestionID: createdQuestion.GetID(),
	}
	updatedInfo, err := client.AddAnswer(ctx, newAnswer)
	if err != nil {
		fmt.Println(err.Error())
	}

	updatedInfo.Question.Statement = "do you think gRPC great?"
	update := &protobuff.QuestionUpdate{
		QuestionInfo: updatedInfo,
		QuestionID:   updatedInfo.GetQuestion().GetID(),
	}

	updatedQuestion, err := client.Update(ctx, update)
	if err != nil {
		fmt.Println(err.Error())
	}

	deleted, err := client.Delete(ctx, wrapperspb.String(createdQuestion.GetID()))
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Printf("Questionary - method findAll(), result len = %v \n\n", len(questions.Questions))
	fmt.Printf("Questionary - method findByUser(), result len = %v \n\n", len(questionsByUser.GetQuestions()))
	fmt.Printf("Questionary - method findByID(), statement = %v\n\n", question.GetQuestion().GetStatement())
	fmt.Printf("Questionary - method Create(), NEW question ID = %v, statement = %v\n\n", createdQuestion.GetID(), createdQuestion.GetStatement())
	fmt.Printf("Questionary - method AddAnswer(), NEW answer ID = %v, answer = %v\n\n", updatedInfo.GetAnswer().GetID(), updatedInfo.GetAnswer().GetAnswer())
	fmt.Printf("Questionary - method Update(), updated ID = %v, info = %v\n\n", updatedQuestion.GetQuestion().GetID(), updatedInfo.String())
	fmt.Printf("Questionary - method Delete(), message = %v\n\n", deleted.GetMessage())
}
