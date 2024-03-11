package queue

import (
	"github.com/google/uuid"

	"{{{ .Package }}}/app/util"
)

type Message struct {
	ID      uuid.UUID `json:"id"`
	Topic   string    `json:"topic"`
	Param   any       `json:"param"`
	Retries int       `json:"retries,omitempty"`
}

func NewMessage(topic string, param any) *Message {
	return &Message{ID: util.UUID(), Topic: topic, Param: param}
}
