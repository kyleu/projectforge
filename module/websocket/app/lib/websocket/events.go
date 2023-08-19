package websocket

import (
	"context"

	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
	"github.com/pkg/errors"

	"{{{ .Package }}}/app/lib/telemetry"
	"{{{ .Package }}}/app/lib/user"{{{ if .HasUser }}}
	dbuser "{{{ .Package }}}/app/user"{{{ end }}}
	"{{{ .Package }}}/app/util"
)

// Registers a new Connection for this Service using the provided user.Profile and websocket.Conn.
func (s *Service) Register({{{ if .HasUser }}}u *dbuser.User, {{{ end }}}profile *user.Profile{{{ if .HasModule "oauth" }}}, accts user.Accounts{{{ end }}}, c *websocket.Conn, logger util.Logger) (*Connection, error) {
	conn := NewConnection("system", {{{ if .HasUser }}}u, {{{ end }}}profile{{{ if .HasModule "oauth" }}}, accts{{{ end }}}, c)
	s.connectionsMu.Lock()
	defer s.connectionsMu.Unlock()
	s.connections[conn.ID] = conn
	if s.onOpen != nil {
		err := s.onOpen(s, conn, logger)
		if err != nil {
			return nil, err
		}
	}
	s.WriteTap(NewMessage(nil, "connection", "open", conn), logger)
	return conn, nil
}

// Sends a message to a provided Connection ID.
func OnMessage(ctx context.Context, s *Service, connID uuid.UUID, message *Message, logger util.Logger) error {
	ctx, span, logger := telemetry.StartSpan(ctx, "message::"+message.Cmd, logger)
	defer span.Complete()

	if connID == systemID {
		logger.Warnf("admin message received: %s", util.ToJSON(message))
		return nil
	}
	s.connectionsMu.Lock()
	c, ok := s.connections[connID]
	s.connectionsMu.Unlock()
	if !ok {
		return invalidConnection(connID)
	}
	s.WriteTap(message, logger)
	return s.handler(ctx, s, c, message.Channel, message.Cmd, message.Param, logger)
}

// Removes a Connection from this Service.
func (s *Service) Disconnect(connID uuid.UUID, logger util.Logger) (bool, error) {
	conn, ok := s.connections[connID]
	if !ok {
		return false, errors.Errorf("no connection found with id [%s]", connID.String())
	}
	left := false
	if conn.Channels != nil {
		left = true
		for _, x := range conn.Channels {
			_, err := s.Leave(connID, x, logger)
			if err != nil {
				return left, errors.Wrapf(err, "error leaving channel [%s]", x)
			}
		}
	}
	s.WriteTap(NewMessage(nil, "connection", "close", conn), logger)
	s.connectionsMu.Lock()
	defer s.connectionsMu.Unlock()
	delete(s.connections, connID)
	return left, nil
}

func invalidConnection(id uuid.UUID) error {
	return errors.Errorf("no connection found with id [%s]", id.String())
}
