//go:generate mockgen -source=repository.go -destination=mocks/repository_mock.go -package mocks
package repository

import (
	"context"

	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
)

//This is the repository interface
//In it you will find all the methods required to interact with the database an do the CRUD operations.
//Each methods has its own validations and error handling.
type Repository interface {

	//Method that search all Questions in the database and returns the result.
	FindAll(ctx context.Context) ([]domain.QuestionInfo, error)

	//Method that search a Question in the database filter by its Unique ID
	FindByID(ctx context.Context, id string) (domain.QuestionInfo, error)

	//Method that find all Questions in the database filter by the User ID
	FindByUser(ctx context.Context, userId string) ([]domain.QuestionInfo, error)

	//Method that Save a new Question in the database
	Create(ctx context.Context, question domain.Question) (domain.Question, error)

	//Method that update the statement and/or answer of an existing Question in the database
	Update(ctx context.Context, questionInfo domain.QuestionInfo) (domain.QuestionInfo, error)

	//Method that delete a Question filter by its unique ID
	Delete(ctx context.Context, id string) (string, error)

	//Method that add a new answer to an existing Question
	AddAnswer(ctx context.Context, answer domain.Answer) (domain.QuestionInfo, error)
}
