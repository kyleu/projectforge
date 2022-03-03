package websocket

import (
	"encoding/json"

	"github.com/fevo-tech/nuevo/app/util"
	"github.com/google/uuid"
)

// Common message struct for passing a service, command and parameter.
type Message struct {
	From    *uuid.UUID      `json:"from,omitempty"`
	Channel string          `json:"channel,omitempty"`
	Cmd     string          `json:"cmd"`
	Param   json.RawMessage `json:"param"`
}

// Constructor.
func NewMessage(userID *uuid.UUID, ch string, cmd string, param interface{}) *Message {
	return &Message{From: userID, Channel: ch, Cmd: cmd, Param: util.ToJSONBytes(param, true)}
}

// Returns a string in "cmd" format, ignoring the param.
func (m *Message) String() string {
	return m.Channel + ":" + m.Cmd
}

// Message for updates of a user's online status.
type OnlineUpdate struct {
	UserID    uuid.UUID `json:"userID"`
	Connected bool      `json:"connected"`
}
