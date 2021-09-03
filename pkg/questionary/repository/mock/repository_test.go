package mock_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
	"github.com/ismaeljpv/qa-api/pkg/questionary/repository"
	"github.com/ismaeljpv/qa-api/pkg/questionary/repository/mock"
	"github.com/stretchr/testify/assert"
)

var logger log.Logger

func NewMockRepository(logger log.Logger) repository.Repository {
	logger = log.NewLogfmtLogger(os.Stderr)
	logger = log.NewSyncLogger(logger)
	logger = log.With(logger,
		"service", "repository_test",
		"time:", log.DefaultTimestampUTC,
		"caller", log.DefaultCaller,
	)
	return mock.NewRepository(logger)
}

type testBody struct {
	value    interface{}
	expected interface{}
}

var ctx = context.Background()

var (
	findByIDDataSuccess = []testBody{
		{value: "1", expected: "1"},
		{value: "2", expected: "2"},
	}

	findByIDDataNotFound = []testBody{
		{value: "asas", expected: "No Question Found"},
		{value: "smas1", expected: "No Question Found"},
	}

	findByUserDataSuccess = []testBody{
		{value: "1", expected: 2},
		{value: "2", expected: 1},
	}

	createQuestionDataSuccess = []testBody{
		{
			value:    domain.Question{ID: "1212", Statement: "Is This A Test?", UserID: "9", CreatedOn: time.Now().Unix()},
			expected: "1212",
		},
		{
			value:    domain.Question{ID: "12124", Statement: "Second Test?", UserID: "10", CreatedOn: time.Now().Unix()},
			expected: "12124",
		},
	}

	createQuestionDataIDConflict = []testBody{
		{
			value:    domain.Question{ID: "1"},
			expected: "Question Already Exists",
		},
		{
			value:    domain.Question{ID: "2"},
			expected: "Question Already Exists",
		},
	}

	addAnswerDataSuccess = []testBody{
		{
			value:    domain.Answer{ID: "wq1212", Answer: "This is an anwser", UserID: "2", QuestionID: "2"},
			expected: "wq1212",
		},
	}

	addAnswerDataQuestionAnsweredConflict = []testBody{
		{
			value:    domain.Answer{ID: "w0212", QuestionID: "1"},
			expected: "The question already has an answer!",
		},
		{
			value:    domain.Answer{ID: "e312", QuestionID: "3"},
			expected: "The question already has an answer!",
		},
	}

	updateQuestionStatementSuccess = []testBody{
		{
			value: domain.QuestionInfo{
				Question: domain.Question{
					ID:        "1",
					Statement: "Do You Think That GOPHERS Rocks?"},
				Answer: domain.Answer{
					ID:         "1",
					Answer:     "Answered!",
					QuestionID: "1",
				},
			},
			expected: "Do You Think That GOPHERS Rocks?",
		},
		{
			value: domain.QuestionInfo{
				Question: domain.Question{
					ID:        "3",
					Statement: "What is a chanel in GOLANG?"},
				Answer: domain.Answer{
					ID:         "2",
					Answer:     "Answered!",
					QuestionID: "3",
				},
			},
			expected: "What is a chanel in GOLANG?",
		},
	}

	updateQuestionAnswerSuccess = []testBody{
		{
			value: domain.QuestionInfo{
				Question: domain.Question{
					ID:        "1",
					Statement: "Do You Think That GOPHERS Rocks?"},
				Answer: domain.Answer{
					ID:         "1",
					Answer:     "Answered over and over!",
					QuestionID: "1",
				},
			},
			expected: "Answered over and over!",
		},
		{
			value: domain.QuestionInfo{
				Question: domain.Question{
					ID:        "3",
					Statement: "What is a chanel in GOLANG?"},
				Answer: domain.Answer{
					ID:         "2",
					Answer:     "Answered Again!",
					QuestionID: "3",
				},
			},
			expected: "Answered Again!",
		},
	}

	updateQuestionInfoNotFound = []testBody{
		{
			value: domain.QuestionInfo{
				Question: domain.Question{
					ID:        "1212",
					Statement: "Do You Think That GOPHERS Rocks?"},
				Answer: domain.Answer{
					ID:         "1",
					Answer:     "Answered!",
					QuestionID: "1",
				},
			},
			expected: "No Question Found To Update",
		},
		{
			value: domain.QuestionInfo{
				Question: domain.Question{
					ID:        "3212",
					Statement: "What is a chanel in GO?"},
				Answer: domain.Answer{
					ID:         "2",
					Answer:     "Answered!",
					QuestionID: "3",
				},
			},
			expected: "No Question Found To Update",
		},
	}

	deleteQuestionSuccess = []testBody{
		{value: "3", expected: "Question Deleted Successfully!"},
		{value: "2", expected: "Question Deleted Successfully!"},
	}

	deleteQuestionNotFound = []testBody{
		{value: "333", expected: "No Question Found"},
		{value: "222", expected: "No Question Found"},
	}
)

func TestFindByID_Success(t *testing.T) {
	repo := NewMockRepository(logger)
	for _, data := range findByIDDataSuccess {
		id := fmt.Sprintf("%v", data.value)
		question, err := repo.FindByID(ctx, id)
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, data.expected, question.Question.ID)
	}
}

func TestFindByID_NotFound(t *testing.T) {
	repo := NewMockRepository(logger)
	for _, data := range findByIDDataNotFound {
		id := fmt.Sprintf("%v", data.value)
		_, err := repo.FindByID(ctx, id)
		if err == nil {
			t.Errorf("Error = [%v] expected", data.expected)
		}
		assert.Equal(t, data.expected, err.Error())
	}
}

func TestFindByUser_Success(t *testing.T) {
	repo := NewMockRepository(logger)
	for _, data := range findByUserDataSuccess {
		userID := fmt.Sprintf("%v", data.value)
		questions, err := repo.FindByUser(ctx, userID)
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, data.expected, len(questions))
	}
}

func TestCreateQuestion_Success(t *testing.T) {
	repo := NewMockRepository(logger)
	for _, data := range createQuestionDataSuccess {
		createdQuestion, err := repo.Create(ctx, data.value.(domain.Question))
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, data.expected, createdQuestion.ID)
	}
}

func TestCreateQuestion_IDConflict(t *testing.T) {
	repo := NewMockRepository(logger)
	for _, data := range createQuestionDataIDConflict {
		_, err := repo.Create(ctx, data.value.(domain.Question))
		if err == nil {
			t.Errorf("Error = [%v] expected", data.expected)
		}
		assert.Equal(t, data.expected, err.Error())
	}
}

func TestAddAnswer_Success(t *testing.T) {
	repo := NewMockRepository(logger)
	for _, data := range addAnswerDataSuccess {
		createdAnswer, err := repo.AddAnswer(ctx, data.value.(domain.Answer))
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, data.expected, createdAnswer.Answer.ID)
	}
}

func TestAddAnswer_QuestionAnsweredConflict(t *testing.T) {
	repo := NewMockRepository(logger)
	for _, data := range addAnswerDataQuestionAnsweredConflict {
		_, err := repo.AddAnswer(ctx, data.value.(domain.Answer))
		if err == nil {
			t.Errorf("Error = [%v] expected", data.expected)
		}
		assert.Equal(t, data.expected, err.Error())
	}
}

func TestUpdateQuestionInfo_Success(t *testing.T) {
	repo := NewMockRepository(logger)

	t.Run("TestUpdateQuestionStatement", func(t *testing.T) {
		for _, data := range updateQuestionStatementSuccess {
			updatedInfo, err := repo.Update(ctx, data.value.(domain.QuestionInfo))
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, data.expected, updatedInfo.Question.Statement)
		}
	})

	t.Run("TestUpdateQuestionAnswer", func(t *testing.T) {
		for _, data := range updateQuestionAnswerSuccess {
			updatedInfo, err := repo.Update(ctx, data.value.(domain.QuestionInfo))
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, data.expected, updatedInfo.Answer.Answer)
		}
	})
}

func TestUpdateQuestionInfo_NotFound(t *testing.T) {
	repo := NewMockRepository(logger)
	for _, data := range updateQuestionInfoNotFound {
		_, err := repo.Update(ctx, data.value.(domain.QuestionInfo))
		if err == nil {
			t.Errorf("Error = [%v] expected", data.expected)
		}
		assert.Equal(t, data.expected, err.Error())
	}
}

func TestDeleteQuestion_Success(t *testing.T) {
	repo := NewMockRepository(logger)
	for _, data := range deleteQuestionSuccess {
		id := fmt.Sprintf("%v", data.value)
		msg, err := repo.Delete(ctx, id)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, data.expected, msg)
	}
}

func TestDeleteQuestion_NotFound(t *testing.T) {
	repo := NewMockRepository(logger)
	for _, data := range deleteQuestionNotFound {
		id := fmt.Sprintf("%v", data.value)
		_, err := repo.Delete(ctx, id)
		if err == nil {
			t.Errorf("Error = [%v] expected", data.expected)
		}
		assert.Equal(t, data.expected, err.Error())
	}
}
