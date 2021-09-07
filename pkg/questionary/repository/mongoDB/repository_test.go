package mongoDB_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	gomock "github.com/golang/mock/gomock"
	"github.com/ismaeljpv/qa-api/pkg/questionary/domain"
	"github.com/ismaeljpv/qa-api/pkg/questionary/repository/mocks"
	"github.com/stretchr/testify/assert"
)

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
		{value: "1", expected: "No Question Found"},
		{value: "2", expected: "No Question Found"},
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
		{
			value:    domain.Answer{ID: "wq1233", Answer: "This is second anwser", UserID: "1", QuestionID: "1"},
			expected: "wq1233",
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
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	mockRepo.EXPECT().FindByID(ctx, "1").Return(domain.QuestionInfo{Question: domain.Question{ID: "1"}}, nil).Times(1)
	mockRepo.EXPECT().FindByID(ctx, "2").Return(domain.QuestionInfo{Question: domain.Question{ID: "2"}}, nil).Times(1)

	for _, data := range findByIDDataSuccess {
		id := fmt.Sprintf("%v", data.value)
		question, err := mockRepo.FindByID(ctx, id)
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, data.expected, question.Question.ID)
	}
}

func TestFindByID_NotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	mockRepo.EXPECT().FindByID(ctx, "1").Return(domain.QuestionInfo{}, errors.New("No Question Found")).Times(1)
	mockRepo.EXPECT().FindByID(ctx, "2").Return(domain.QuestionInfo{}, errors.New("No Question Found")).Times(1)

	for _, data := range findByIDDataNotFound {
		id := fmt.Sprintf("%v", data.value)
		_, err := mockRepo.FindByID(ctx, id)
		if err == nil {
			t.Errorf("Error = [%v] expected", data.expected)
		}
		assert.Equal(t, data.expected, err.Error())
	}
}

func TestFindByUser_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	mockRepo.EXPECT().FindByUser(ctx, "1").Return([]domain.QuestionInfo{{}, {}}, nil).Times(1)
	mockRepo.EXPECT().FindByUser(ctx, "2").Return([]domain.QuestionInfo{{}}, nil).Times(1)

	for _, data := range findByUserDataSuccess {
		userID := fmt.Sprintf("%v", data.value)
		questions, err := mockRepo.FindByUser(ctx, userID)
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, data.expected, len(questions))
	}
}

func TestCreateQuestion_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(domain.Question{ID: "1212"}, nil).Times(1)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(domain.Question{ID: "12124"}, nil).Times(1)

	for _, data := range createQuestionDataSuccess {
		createdQuestion, err := mockRepo.Create(ctx, data.value.(domain.Question))
		if err != nil {
			t.Error(err.Error())
		}
		assert.Equal(t, data.expected, createdQuestion.ID)
	}
}

func TestCreateQuestion_IDConflict(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(domain.Question{}, errors.New("Question Already Exists")).Times(1)
	mockRepo.EXPECT().Create(ctx, gomock.Any()).Return(domain.Question{}, errors.New("Question Already Exists")).Times(1)

	for _, data := range createQuestionDataIDConflict {
		_, err := mockRepo.Create(ctx, data.value.(domain.Question))
		if err == nil {
			t.Errorf("Error = [%v] expected", data.expected)
		}
		assert.Equal(t, data.expected, err.Error())
	}
}

func TestAddAnswer_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	mockRepo.EXPECT().AddAnswer(ctx, gomock.Any()).Return(domain.QuestionInfo{Answer: domain.Answer{ID: "wq1212"}}, nil).Times(1)
	mockRepo.EXPECT().AddAnswer(ctx, gomock.Any()).Return(domain.QuestionInfo{Answer: domain.Answer{ID: "wq1233"}}, nil).Times(1)

	for _, data := range addAnswerDataSuccess {
		createdAnswer, err := mockRepo.AddAnswer(ctx, data.value.(domain.Answer))
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, data.expected, createdAnswer.Answer.ID)
	}
}

func TestAddAnswer_QuestionAnsweredConflict(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	mockRepo.EXPECT().AddAnswer(ctx, gomock.Any()).Return(domain.QuestionInfo{}, errors.New("The question already has an answer!")).Times(1)
	mockRepo.EXPECT().AddAnswer(ctx, gomock.Any()).Return(domain.QuestionInfo{}, errors.New("The question already has an answer!")).Times(1)

	for _, data := range addAnswerDataQuestionAnsweredConflict {
		_, err := mockRepo.AddAnswer(ctx, data.value.(domain.Answer))
		if err == nil {
			t.Errorf("Error = [%v] expected", data.expected)
		}
		assert.Equal(t, data.expected, err.Error())
	}
}

func TestUpdateQuestionInfo_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)

	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(domain.QuestionInfo{
		Question: domain.Question{
			ID:        "1",
			Statement: "Do You Think That GOPHERS Rocks?"},
		Answer: domain.Answer{
			ID:         "1",
			Answer:     "Answered!",
			QuestionID: "1",
		},
	}, nil).Times(1)

	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(domain.QuestionInfo{
		Question: domain.Question{
			ID:        "3",
			Statement: "What is a chanel in GOLANG?"},
		Answer: domain.Answer{
			ID:         "2",
			Answer:     "Answered!",
			QuestionID: "3",
		},
	}, nil).Times(1)

	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(domain.QuestionInfo{
		Question: domain.Question{
			ID:        "1",
			Statement: "Do You Think That GOPHERS Rocks?"},
		Answer: domain.Answer{
			ID:         "1",
			Answer:     "Answered over and over!",
			QuestionID: "1",
		},
	}, nil).Times(1)

	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(domain.QuestionInfo{
		Question: domain.Question{
			ID:        "3",
			Statement: "What is a chanel in GOLANG?"},
		Answer: domain.Answer{
			ID:         "2",
			Answer:     "Answered Again!",
			QuestionID: "3",
		},
	}, nil).Times(1)

	t.Run("TestUpdateQuestionStatement", func(t *testing.T) {
		for _, data := range updateQuestionStatementSuccess {
			updatedInfo, err := mockRepo.Update(ctx, data.value.(domain.QuestionInfo))
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, data.expected, updatedInfo.Question.Statement)
		}
	})

	t.Run("TestUpdateQuestionAnswer", func(t *testing.T) {
		for _, data := range updateQuestionAnswerSuccess {
			updatedInfo, err := mockRepo.Update(ctx, data.value.(domain.QuestionInfo))
			if err != nil {
				t.Error(err)
			}
			assert.Equal(t, data.expected, updatedInfo.Answer.Answer)
		}
	})
}

func TestUpdateQuestionInfo_NotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	mockRepo.EXPECT().Update(ctx, gomock.Any()).Return(domain.QuestionInfo{}, errors.New("No Question Found To Update")).Times(2)

	for _, data := range updateQuestionInfoNotFound {
		_, err := mockRepo.Update(ctx, data.value.(domain.QuestionInfo))
		if err == nil {
			t.Errorf("Error = [%v] expected", data.expected)
		}
		assert.Equal(t, data.expected, err.Error())
	}
}

func TestDeleteQuestion_Success(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	mockRepo.EXPECT().Delete(ctx, "3").Return("Question Deleted Successfully!", nil).Times(1)
	mockRepo.EXPECT().Delete(ctx, "2").Return("Question Deleted Successfully!", nil).Times(1)

	for _, data := range deleteQuestionSuccess {
		id := fmt.Sprintf("%v", data.value)
		msg, err := mockRepo.Delete(ctx, id)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, data.expected, msg)
	}
}

func TestDeleteQuestion_NotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	mockRepo := mocks.NewMockRepository(mockCtrl)
	mockRepo.EXPECT().Delete(ctx, "333").Return("", errors.New("No Question Found")).Times(1)
	mockRepo.EXPECT().Delete(ctx, "222").Return("", errors.New("No Question Found")).Times(1)

	for _, data := range deleteQuestionNotFound {
		id := fmt.Sprintf("%v", data.value)
		_, err := mockRepo.Delete(ctx, id)
		if err == nil {
			t.Errorf("Error = [%v] expected", data.expected)
		}
		assert.Equal(t, data.expected, err.Error())
	}
}
