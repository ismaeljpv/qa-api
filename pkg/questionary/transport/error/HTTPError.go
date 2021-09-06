package error

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// HTTPError implements ClientError interface.
type HTTPError struct {
	Cause  error  `json:"-"`
	Detail string `json:"detail"`
	Status int    `json:"status"`
}

func NewClientError(err error, status int, detail string) ClientError {
	return &HTTPError{
		Cause:  err,
		Detail: detail,
		Status: status,
	}
}

func NewServerError(err error, detail string) InternalServerError {
	return &HTTPError{
		Cause:  err,
		Detail: detail,
		Status: http.StatusInternalServerError,
	}
}

func (e *HTTPError) Error() string {
	return e.Detail
}

// ResponseBody returns JSON response body.
func (e *HTTPError) ResponseBody() ([]byte, error) {
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("Error while parsing response body: %v", err)
	}
	return body, nil
}

// ResponseHeaders returns http status code and headers.
func (e *HTTPError) ResponseHeaders() (int, map[string]string) {
	return e.Status, map[string]string{
		"Content-Type": "application/json; charset=utf-8",
	}
}