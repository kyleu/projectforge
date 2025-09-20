package websocket

import (
	"encoding/json/jsontext"

	"github.com/google/uuid"

	"{{{ .Package }}}/app/util"
)

type Message struct {
	From    *uuid.UUID     `json:"from,omitempty"`
	Channel string         `json:"channel,omitempty"`
	Cmd     string         `json:"cmd"`
	Param   jsontext.Value `json:"param"`
}

func NewMessage(userID *uuid.UUID, ch string, cmd string, param any) *Message {
	return &Message{From: userID, Channel: ch, Cmd: cmd, Param: util.ToJSONBytes(param, true)}
}

func (m *Message) String() string {
	return m.Channel + ":" + m.Cmd
}

type OnlineUpdate struct {
	UserID    uuid.UUID `json:"userID"`
	Connected bool      `json:"connected"`
}
