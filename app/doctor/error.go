package doctor

import (
	"fmt"

	"github.com/samber/lo"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Args    []any  `json:"args,omitempty"`
}

func NewError(code string, message string, args ...any) *Error {
	return &Error{Code: code, Message: message, Args: args}
}

func (e *Error) String() string {
	return fmt.Sprintf(e.Message, e.Args...)
}

type Errors []*Error

func (e Errors) Find(code string) *Error {
	return lo.FindOrElse(e, nil, func(x *Error) bool {
		return x.Code == code
	})
}
