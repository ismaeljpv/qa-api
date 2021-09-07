package service

import (
	"context"

	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
)

//This is the service implementation interface
//In it you will find all the methods required to implement the business logic of the Questionary Microservice
//Each methods has its own validations.
type Service interface {

	//Method that returns all question avaliables in the database.
	FindAll(ctx context.Context) ([]domain.QuestionInfo, error)

	//Method that find and return a question with its anwers by its unique ID
	FindByID(ctx context.Context, id string) (domain.QuestionInfo, error)

	//Method that finds all questions asked by a user related by the User ID
	FindByUser(ctx context.Context, userId string) ([]domain.QuestionInfo, error)

	//Method that create a new Question
	Create(ctx context.Context, question domain.Question) (domain.Question, error)

	//Method that Update a Question and/or its anwser
	Update(ctx context.Context, questionInfo domain.QuestionInfo, id string) (domain.QuestionInfo, error)

	//Method that delete a Question by its unique ID
	Delete(ctx context.Context, id string) (string, error)

	//Method that adds an anwer to a existing Question
	AddAnswer(ctx context.Context, answer domain.Answer) (domain.QuestionInfo, error)
}
