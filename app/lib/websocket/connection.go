// Content managed by Project Forge, see [projectforge.md] for details.
package websocket

import (
	"sync"

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
	Accounts user.Accounts `json:"accounts,omitempty"`
	Svc      string        `json:"svc,omitempty"`
	ModelID  *uuid.UUID    `json:"modelID,omitempty"`
	Channels []string      `json:"channels,omitempty"`
	socket   *websocket.Conn
	mu       sync.Mutex
}

// Creates a new Connection.
func NewConnection(svc string, profile *user.Profile, accounts user.Accounts, socket *websocket.Conn) *Connection {
	return &Connection{ID: util.UUID(), Profile: profile, Accounts: accounts, Svc: svc, socket: socket}
}

// Transforms this Connection to a serializable Status object.
func (c *Connection) ToStatus() *Status {
	if c.Channels == nil {
		return &Status{ID: c.ID, Username: c.Profile.Name, Channels: nil}
	}
	return &Status{ID: c.ID, Username: c.Profile.Name, Channels: c.Channels}
}

// Writes bytes to this Connection, you should probably use a helper method.
func (c *Connection) Write(b []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.socket.WriteMessage(websocket.TextMessage, b)
	return errors.Wrap(err, "unable to write to websocket")
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

// Registers a new Connection for this Service using the provided fevouser.Profile and websocket.Conn.
func (s *Service) Register(profile *user.Profile, c *websocket.Conn) (*Connection, error) {
	conn := &Connection{ID: util.UUID(), Profile: profile, Svc: "system", socket: c}

	s.connectionsMu.Lock()
	defer s.connectionsMu.Unlock()

	s.connections[conn.ID] = conn
	if s.onOpen != nil {
		err := s.onOpen(s, conn)
		if err != nil {
			return nil, err
		}
	}
	return conn, nil
}

// Removes a Connection from this Service.
func (s *Service) Disconnect(connID uuid.UUID) (bool, error) {
	conn, ok := s.connections[connID]
	if !ok {
		return false, errors.Errorf("no connection found with id [%s]", connID.String())
	}
	left := false

	if conn.Channels != nil {
		left = true
		for _, x := range conn.Channels {
			_, err := s.Leave(connID, x)
			if err != nil {
				return left, errors.Wrap(err, "error leaving channel ["+x+"]")
			}
		}
	}

	s.connectionsMu.Lock()
	defer s.connectionsMu.Unlock()

	delete(s.connections, connID)
	return left, nil
}

func invalidConnection(id uuid.UUID) error {
	return errors.Errorf("no connection found with id [%s]", id.String())
}
