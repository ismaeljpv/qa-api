package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gofrs/uuid"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
	"github.com/ismaeljpv/qa-api/pkg/questionary/repository"
)

func ImplementationError(e string) error {
	return errors.New(e)
}

type service struct {
	repository repository.Repository
	logger     log.Logger
}

func NewService(rep repository.Repository, logger log.Logger) Service {
	return &service{
		repository: rep,
		logger:     logger,
	}
}

func (s *service) FindAll(ctx context.Context) ([]domain.QuestionInfo, error) {
	questions, err := s.repository.FindAll(ctx)
	if err != nil {
		return []domain.QuestionInfo{}, err
	}
	return questions, nil
}

func (s *service) FindByID(ctx context.Context, id string) (domain.QuestionInfo, error) {
	questionInfo, err := s.repository.FindByID(ctx, id)
	if err != nil {
		return domain.QuestionInfo{}, err
	}
	return questionInfo, nil
}

func (s *service) FindByUser(ctx context.Context, userId string) ([]domain.QuestionInfo, error) {
	userQuestions, err := s.repository.FindByUser(ctx, userId)
	if err != nil {
		return []domain.QuestionInfo{}, err
	}
	return userQuestions, nil
}

func (s *service) Create(ctx context.Context, question domain.Question) (domain.Question, error) {
	uuid, idErr := uuid.NewV4()
	if idErr != nil {
		level.Warn(s.logger).Log("msg", "Error creating uuid for Question, method Create")
		return domain.Question{}, idErr
	}

	question.ID = uuid.String()
	question.CreatedOn = time.Now().Unix()
	createdQuestion, err := s.repository.Create(ctx, question)
	if err != nil {
		return createdQuestion, err
	}
	return createdQuestion, nil
}

func (s *service) Update(ctx context.Context, questionInfo domain.QuestionInfo, id string) (domain.QuestionInfo, error) {
	if questionInfo.Question.ID != id {
		level.Warn(s.logger).Log("msg", fmt.Sprintf("The Path Param ID doesnt match with the body ID [%v!=%v], method update", questionInfo.Question.ID, id))
		return domain.QuestionInfo{}, ImplementationError("There is a inconsistency with the information of the request")
	}

	if questionInfo.Answer.ID == "" {
		level.Warn(s.logger).Log("msg", "The answer provided in the request doesnt have an ID, method update")
		return domain.QuestionInfo{}, ImplementationError("The answer passed to update is not valid")
	}

	updatedQuestion, err := s.repository.Update(ctx, questionInfo)
	if err != nil {
		return updatedQuestion, err
	}

	return updatedQuestion, nil
}

func (s *service) Delete(ctx context.Context, id string) (string, error) {
	msg, err := s.repository.Delete(ctx, id)
	if err != nil {
		return "", err
	}
	return msg, nil
}

func (s *service) AddAnswer(ctx context.Context, answer domain.Answer) (domain.QuestionInfo, error) {
	uuid, idErr := uuid.NewV4()
	if idErr != nil {
		level.Warn(s.logger).Log("msg", "Error creating uuid for Answer, method AddAnswer")
		return domain.QuestionInfo{}, idErr
	}

	answer.ID = uuid.String()
	answer.CreatedOn = time.Now().Unix()
	questionInfo, err := s.repository.AddAnswer(ctx, answer)
	if err != nil {
		return questionInfo, err
	}
	return questionInfo, nil
}
