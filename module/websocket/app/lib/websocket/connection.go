package websocket

import (
	"fmt"
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/user"{{{ if .HasModule "user" }}}
	dbuser "{{{ .Package }}}/app/user"{{{ end }}}
	"{{{ .Package }}}/app/util"
)

// Represents a user's WebSocket session.
type Connection struct {
	ID       uuid.UUID     `json:"id"`{{{ if .HasModule "user" }}}
	User     *dbuser.User  `json:"user,omitempty"`{{{ end }}}
	Profile  *user.Profile `json:"profile,omitempty"`{{{ if .HasModule "oauth" }}}
	Accounts user.Accounts `json:"accounts,omitempty"`{{{ end }}}
	Svc      string        `json:"svc,omitempty"`
	ModelID  *uuid.UUID    `json:"modelID,omitempty"`
	Channels []string      `json:"channels,omitempty"`
	socket   *websocket.Conn
	mu       sync.Mutex
}

// Creates a new Connection.
func NewConnection(svc string{{{ if .HasModule "user" }}}, user *dbuser.User{{{ end }}}, profile *user.Profile{{{ if .HasModule "oauth" }}}, accounts user.Accounts{{{ end }}}, socket *websocket.Conn) *Connection {
	return &Connection{ID: util.UUID(){{{ if .HasModule "user" }}}, User: user{{{ end }}}, Profile: profile{{{ if .HasModule "oauth" }}}, Accounts: accounts{{{ end }}}, Svc: svc, socket: socket}
}

// Transforms this Connection to a serializable Status object.
func (c *Connection) ToStatus() *Status {
	if c.Channels == nil {
		return &Status{ID: c.ID, Username: c.Profile.Name, Channels: nil}
	}
	return &Status{ID: c.ID, Username: c.Profile.Name, Channels: c.Channels}
}{{{ if .HasModule "user" }}}

func (c *Connection) Username() string {
	if c.User != nil {
		return c.User.Name
	}
	return c.Profile.Name
}{{{ else }}}

func (c *Connection) Username() string {
	return c.Profile.Name
}{{{ end }}}

// Writes bytes to this Connection, you should probably use a helper method.
func (c *Connection) Write(b []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.socket.WriteMessage(websocket.TextMessage, b)
	if err != nil {
		return errors.Wrap(err, "unable to write to websocket")
	}
	return nil
}

// Reads bytes from this Connection, you should probably use a helper method.
func (c *Connection) Read() ([]byte, error) {
	_, message, err := c.socket.ReadMessage()
	return message, errors.Wrap(err, "unable to write to websocket")
}

// Closes the backing socket.
func (c *Connection) Close() error {
	return c.socket.Close()
}

func (c *Connection) String() string {
	return fmt.Sprintf("[%s][%s::%s][%s]", c.ID, c.Svc, c.ModelID, c.Profile.String())
}
