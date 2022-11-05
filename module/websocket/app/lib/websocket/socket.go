package websocket

import (
	"fmt"
	"sync/atomic"

	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"golang.org/x/exp/slices"

	"{{{ .Package }}}/app/util"
)

var count int64

func (s *Service) Write(connID uuid.UUID, message string) error {
	if connID == systemID {
		s.Logger.Warn("--- admin message sent ---")
		s.Logger.Warn(message)
		return nil
	}

	s.connectionsMu.Lock()
	c, ok := s.connections[connID]
	s.connectionsMu.Unlock()

	atomic.AddInt64(&count, 1)

	if !ok {
		return errors.New("cannot load connection [" + connID.String() + "]")
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.socket.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return errors.Wrap(err, "unable to write to websocket")
	}
	return nil
}

func (s *Service) WriteMessage(connID uuid.UUID, message *Message) error {
	return s.Write(connID, util.ToJSON(message))
}

func (s *Service) WriteChannel(message *Message, except ...uuid.UUID) error {
	if message == nil {
		return errors.New("no message provided")
	}
	if message.Channel == "" {
		return errors.New("no channel provided")
	}
	conns, ok := s.channels[message.Channel]
	if !ok {
		return nil
	}

	s.Logger.Debug(fmt.Sprintf("sending message [%v::%v] to [%v] connections", message.Channel, message.Cmd, len(conns.MemberIDs)))
	for _, conn := range conns.MemberIDs {
		if !slices.Contains(except, conn) {
			connID := conn

			go func() {
				defer func() { _ = recover() }()
				_ = s.Write(connID, util.ToJSON(message))
			}()
		}
	}
	return nil
}

func (s *Service) ReadLoop(connID uuid.UUID, onDisconnect func(conn *Connection) error) error {
	c, ok := s.connections[connID]
	if !ok {
		return errors.New("cannot load connection [" + connID.String() + "]")
	}

	defer func() {
		_ = c.socket.Close()
		if onDisconnect != nil {
			err := onDisconnect(c)
			if err != nil {
				s.Logger.Warn(fmt.Sprintf("error running onDisconnect for [%v]: %+v", connID, err))
			}
		}
		_, err := s.Disconnect(connID)
		if err != nil {
			s.Logger.Warn(fmt.Sprintf("error running disconnect for [%v]: %+v", connID, err))
		}
		s.Logger.Debug(fmt.Sprintf("closed websocket [%v]", connID.String()))
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		if err != nil {
			return err
		}

		var m Message
		err = util.FromJSON(message, &m)
		if err != nil {
			return errors.Wrap(err, "error decoding websocket message")
		}

		err = OnMessage(s, connID, &m)
		if err != nil {
			_ = s.WriteMessage(c.ID, NewMessage(nil, "system", "error", err.Error()))
			return errors.Wrap(err, "error handling websocket message")
		}
	}
}
