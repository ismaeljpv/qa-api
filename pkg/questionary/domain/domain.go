package domain

type Question struct {
	ID        string `json:"id"`
	Statement string `json:"statement"`
	UserID    string `json:"userId"`
	CreatedOn int64  `json:"createdOn"`
}

type Answer struct {
	ID         string `json:"id,omitempty"`
	Answer     string `json:"anwser"`
	QuestionID string `json:"questionId"`
	UserID     string `json:"userId"`
	CreatedOn  int64  `json:"createdOn"`
}

type QuestionInfo struct {
	Question Question `json:"questionInfo"`
	Answer   Answer   `json:"answer,omitempty"`
}
