package http

type Query struct {
	Address string
}

type Error struct {
	error
	Code    int
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func NewError(code int, message string) *Error {

	e := Error{
		Code:    code,
		Message: message,
	}

	return &e
}
