package doctor

import (
	"fmt"
)

type Error struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Args    []interface{} `json:"args,omitempty"`
}

func NewError(code string, message string, args ...interface{}) *Error {
	return &Error{Code: code, Message: message, Args: args}
}

func (e *Error) String() string {
	msg := fmt.Sprintf(e.Message, e.Args...)
	return e.Code + ": " + msg
}

type Errors []*Error

func (e Errors) Find(code string) *Error {
	for _, x := range e {
		if x.Code == code {
			return x
		}
	}
	return nil
}
