package web

// FieldError is used to indicate an error with a specific request field.
type FieldError struct {
	Field string `json:"field"`
	Error string `json:"error"`
}

// ErrorResponse how we respond to clients when something goes wrong.
type ErrorResponse struct {
	Error  string       `json:"error"`
	Fields []FieldError `json:"fields,omitempty"`
}

// Error is used to add web information to a request error.
type Error struct {
	Err    error
	Status int
	Fields []FieldError
}

func (e *Error) Error() string {
	return e.Err.Error()
}

// NewRequestError is used when a know error condition is encountered.
func NewRequestError(err error, status int) error {
	return &Error{Err: err, Status: status}
}
