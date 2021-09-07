package domain

//Domains that represents the structures of the database schemas

type Question struct {
	ID        string `json:"id,omitempty"`
	Statement string `json:"statement" validate:"required"`
	UserID    string `json:"userId" validate:"required"`
	CreatedOn int64  `json:"createdOn,omitempty"`
}

type Answer struct {
	ID         string `json:"id,omitempty"`
	Answer     string `json:"anwser,omitempty" validate:"required"`
	QuestionID string `json:"questionId,omitempty" validate:"required"`
	UserID     string `json:"userId,omitempty" validate:"required"`
	CreatedOn  int64  `json:"createdOn,omitempty"`
}

type QuestionInfo struct {
	Question Question `json:"question" validate:"required"`
	Answer   Answer   `json:"answer" validate:"required"`
}
