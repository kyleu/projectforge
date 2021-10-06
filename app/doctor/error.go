package doctor

import (
	"fmt"
)

type Error struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Args    []interface{} `json:"args,omitempty"`
}

func NewError(code string, message string, args ...interface{}) *Error {
	return &Error{Code: code, Message: message, Args: args}
}

func (e *Error) String() string {
	return fmt.Sprintf("%s: %s", e.Code, fmt.Sprintf(e.Message, e.Args...))
}

type Errors []*Error
