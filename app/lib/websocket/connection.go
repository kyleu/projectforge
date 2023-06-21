// Content managed by Project Forge, see [projectforge.md] for details.
package websocket

import (
	"fmt"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
)

// Represents a user's WebSocket session.
type Connection struct {
	ID       uuid.UUID     `json:"id"`
	Profile  *user.Profile `json:"profile,omitempty"`
	Svc      string        `json:"svc,omitempty"`
	ModelID  *uuid.UUID    `json:"modelID,omitempty"`
	Channels []string      `json:"channels,omitempty"`
	Started  time.Time     `json:"started,omitempty"`
	socket   *websocket.Conn
	mu       sync.Mutex
}

// Creates a new Connection.
func NewConnection(svc string, profile *user.Profile, socket *websocket.Conn) *Connection {
	return &Connection{ID: util.UUID(), Profile: profile, Svc: svc, Started: time.Now(), socket: socket}
}

// Transforms this Connection to a serializable Status object.
func (c *Connection) ToStatus() *Status {
	if c.Channels == nil {
		return &Status{ID: c.ID, Username: c.Profile.Name, Channels: nil}
	}
	return &Status{ID: c.ID, Username: c.Profile.Name, Channels: c.Channels}
}

func (c *Connection) Username() string {
	return c.Profile.Name
}

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
