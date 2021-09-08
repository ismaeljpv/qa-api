package mockDB

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
	repo "github.com/ismaeljpv/qa-api/pkg/questionary/repository"
	httpError "github.com/ismaeljpv/qa-api/pkg/questionary/transport/http/error"
)

//This is a Mock implementation of a real database using in memory data
//The Repository interface is implemented
var questionData = []domain.QuestionInfo{
	{
		Question: domain.Question{
			ID:        "1",
			Statement: "Do You Think That GO Rocks?",
			UserID:    "1",
			CreatedOn: time.Now().Unix(),
		},
		Answer: domain.Answer{
			ID:         "1",
			Answer:     "Yes!",
			QuestionID: "1",
			UserID:     "2",
			CreatedOn:  time.Now().Unix(),
		},
	},
	{
		Question: domain.Question{
			ID:        "2",
			Statement: "Where are all the gophers?",
			UserID:    "1",
			CreatedOn: time.Now().Unix(),
		},
	},
	{
		Question: domain.Question{
			ID:        "3",
			Statement: "What is a chanel in GO?",
			UserID:    "2",
			CreatedOn: time.Now().Unix(),
		},
		Answer: domain.Answer{
			ID:         "2",
			Answer:     "Is the conduct that let goroutines to comunicate to each other",
			QuestionID: "3",
			UserID:     "1",
			CreatedOn:  time.Now().Unix(),
		},
	},
}

type repository struct {
	db     []domain.QuestionInfo
	logger log.Logger
}

func NewRepository(logger log.Logger) repo.Repository {
	return &repository{
		db:     questionData,
		logger: logger,
	}
}

func (r *repository) FindAll(ctx context.Context) ([]domain.QuestionInfo, error) {
	DBQuestions := r.db
	return DBQuestions, nil
}

func (r *repository) FindByID(ctx context.Context, id string) (domain.QuestionInfo, error) {
	DBQuestions := r.db
	for _, questionInfo := range DBQuestions {
		if questionInfo.Question.ID == id {
			return questionInfo, nil
		}
	}
	level.Warn(r.logger).Log("msg", fmt.Sprintf("No Question Found by ID %v, method FindByID", id))
	return domain.QuestionInfo{}, httpError.NewClientError(errors.New(fmt.Sprintf("No question found by ID %v", id)),
		http.StatusNotFound,
		"No Question Found")
}

func (r *repository) FindByUser(ctx context.Context, userId string) ([]domain.QuestionInfo, error) {
	userQuestions := []domain.QuestionInfo{}
	DBQuestions := r.db
	for _, questionInfo := range DBQuestions {
		if questionInfo.Question.UserID == userId {
			userQuestions = append(userQuestions, questionInfo)
		}
	}
	return userQuestions, nil
}

func (r *repository) Create(ctx context.Context, question domain.Question) (domain.Question, error) {
	for _, questionInfo := range r.db {
		if questionInfo.Question.ID == question.ID {
			return domain.Question{}, httpError.NewClientError(errors.New("Conflict - Question already exists"),
				http.StatusConflict,
				"Question Already Exists")
		}
	}
	r.db = append(r.db, domain.QuestionInfo{Question: question})
	return question, nil
}

func (r *repository) Update(ctx context.Context, questionInfo domain.QuestionInfo) (domain.QuestionInfo, error) {
	var updated bool
	for i, questionData := range r.db {
		if questionData.Question.ID == questionInfo.Question.ID {

			if strings.Compare(questionData.Question.Statement, questionInfo.Question.Statement) != 0 {
				r.db[i].Question.Statement = questionInfo.Question.Statement
				updated = true
			}

			if r.db[i].Answer.ID != "" && r.db[i].Answer.ID == questionInfo.Answer.ID {
				if strings.Compare(questionData.Answer.Answer, questionInfo.Answer.Answer) != 0 {
					r.db[i].Answer.Answer = questionInfo.Answer.Answer
					updated = true
				}
			}

			if updated {
				return r.db[i], nil
			}
		}
	}

	level.Warn(r.logger).Log("msg", fmt.Sprintf("No Question Found by ID %v, method Update", questionInfo.Question.ID))
	return domain.QuestionInfo{}, httpError.NewClientError(errors.New(fmt.Sprintf("No question found by ID %v", questionInfo.Question.ID)),
		http.StatusNotFound,
		"No Question Found To Update")
}

func (r *repository) Delete(ctx context.Context, id string) (string, error) {
	var deleted bool
	for i, questionInfo := range r.db {
		if questionInfo.Question.ID == id {
			r.db = append(r.db[:i], r.db[i+1:]...)
			deleted = true
		}
	}

	if deleted {
		return "Question Deleted Successfully!", nil
	} else {
		level.Warn(r.logger).Log("msg", fmt.Sprintf("No Question Found by ID %v, method Delete", id))
		return "", httpError.NewClientError(errors.New(fmt.Sprintf("No question found by ID %v", id)),
			http.StatusNotFound,
			"No Question Found")
	}
}

func (r repository) AddAnswer(ctx context.Context, answer domain.Answer) (domain.QuestionInfo, error) {
	for i, questionInfo := range r.db {
		if questionInfo.Question.ID == answer.QuestionID {
			if r.db[i].Answer.ID == "" {
				r.db[i].Answer = answer
				return r.db[i], nil
			} else {
				return domain.QuestionInfo{}, httpError.NewClientError(errors.New("Question is already answered"),
					http.StatusConflict,
					"The question already has an answer!")
			}
		}
	}
	level.Warn(r.logger).Log("msg", fmt.Sprintf("No Question Found by ID %v, method AddAnswer", answer.QuestionID))
	return domain.QuestionInfo{}, httpError.NewClientError(errors.New(fmt.Sprintf("No question found by ID %v", answer.QuestionID)),
		http.StatusNotFound,
		"No Question Found")
}
