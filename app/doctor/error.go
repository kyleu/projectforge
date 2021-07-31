package doctor

type Error struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Args    []string `json:"args,omitempty"`
}

func NewError(code string, message string, args ...string) *Error {
	return &Error{Code: code, Message: message, Args: args}
}

func (e *Error) Solution() *Solution {
	return FindSolution(e.Code, e.Args...)
}

type Errors []*Error
