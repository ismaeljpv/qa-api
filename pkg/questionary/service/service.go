package service

import (
	"context"

	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
)

type Service interface {
	FindAll(ctx context.Context) ([]domain.QuestionInfo, error)
	FindByID(ctx context.Context, id string) (domain.QuestionInfo, error)
	FindByUser(ctx context.Context, userId string) ([]domain.QuestionInfo, error)
	Create(ctx context.Context, question domain.Question) (domain.Question, error)
	Update(ctx context.Context, questionInfo domain.QuestionInfo, id string) (domain.QuestionInfo, error)
	Delete(ctx context.Context, id string) (string, error)
	AddAnswer(ctx context.Context, answer domain.Answer) (domain.QuestionInfo, error)
}
