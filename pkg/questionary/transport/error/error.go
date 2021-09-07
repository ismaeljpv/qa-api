package error

//This interface represents the errors produced by the data recieved in the request od the client
//The methods of the interface represents a standar HTTP error
type ClientError interface {
	Error() string
	// ResponseBody returns response body.
	ResponseBody() ([]byte, error)
	// ResponseHeaders returns http status code and headers.
	ResponseHeaders() (int, map[string]string)
}

//This interface represents the errors produced by the internally when the server is unable to handle the return a proper response
//The methods of the interface represents a standar HTTP error
type InternalServerError interface {
	Error() string
	// ResponseBody returns response body.
	ResponseBody() ([]byte, error)
	// ResponseHeaders returns http status code and headers.
	ResponseHeaders() (int, map[string]string)
}
