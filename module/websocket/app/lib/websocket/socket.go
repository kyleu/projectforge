package websocket

import (
	"context"
	"slices"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func (s *Service) write(connID uuid.UUID, message string, logger util.Logger) error {
	if connID == systemID {
		logger.Warnf("admin message sent: %s", message)
		return nil
	}

	s.connectionsMu.Lock()
	c, ok := s.connections[connID]
	if ok {
		c.Stats.Sent(len(message))
	}
	s.stats.Sent(len(message))
	s.connectionsMu.Unlock()

	if !ok {
		return errors.Errorf("cannot load connection [%s] for writing", connID.String())
	}
	c.mu.Lock()
	defer c.mu.Unlock()

	err := c.socket.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		if err.Error() == "websocket: close sent" {
			return nil
		}
		return errors.Wrap(err, "unable to write to websocket")
	}
	return nil
}

func (s *Service) WriteMessage(connID uuid.UUID, message *Message, logger util.Logger) error {
	s.WriteTap(message, logger)
	return s.write(connID, util.ToJSON(message), logger)
}

func (s *Service) WriteCloseRequest(connID uuid.UUID, logger util.Logger) error {
	message := &Message{Cmd: "close-connection"}
	return s.WriteMessage(connID, message, logger)
}

func (s *Service) WriteChannel(message *Message, logger util.Logger, except ...uuid.UUID) error {
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
	json := util.ToJSON(message)
	lo.ForEach(conns.ConnIDs, func(conn uuid.UUID, _ int) {
		if !lo.Contains(except, conn) {
			connID := conn

			go func() {
				defer func() { _ = recover() }()
				_ = s.write(connID, json, logger)
			}()
		}
	})
	s.WriteTap(message, logger)
	return nil
}

func (s *Service) Broadcast(message *Message, logger util.Logger, except ...uuid.UUID) error {
	logger.Debugf("broadcasting message [%v::%v] to [%v] connections", message.Channel, message.Cmd, len(s.connections))
	for id := range s.connections {
		if !slices.Contains(except, id) {
			closureID := id
			go func() {
				_ = s.WriteMessage(closureID, message, logger)
			}()
		}
	}
	return nil
}

func (s *Service) ReadLoop(ctx context.Context, connID uuid.UUID, logger util.Logger) error {
	c, ok := s.connections[connID]
	if !ok {
		return errors.Errorf("cannot load connection [%s] for reading", connID.String())
	}
	d := func() error {
		_, err := s.Disconnect(connID, logger)
		if err != nil {
			return err
		}
		if s.onClose != nil {
			return s.onClose(s, c, logger)
		}
		return nil
	}
	m := func(m *Message) error {
		return OnMessage(ctx, s, connID, m, logger)
	}
	return ReadSocketLoop(connID, c.socket, m, d, logger, s.stats, c.Stats)
}

func ReadSocketLoop(
	connID uuid.UUID, sock *websocket.Conn, onMessage func(m *Message) error, onDisconnect func() error, logger util.Logger, stats ...*Stats,
) error {
	defer func() {
		_ = sock.Close()
		if onDisconnect != nil {
			err := onDisconnect()
			if err != nil {
				logger.Warnf("error running onDisconnect for [%v]: %+v", connID, err)
			}
		}
	}()

	for {
		_, message, err := sock.ReadMessage()
		if err != nil {
			errMsg := err.Error()
			if strings.Contains(errMsg, "1001") || strings.Contains(errMsg, "close 1005") {
				return nil
			}
			return errors.Wrapf(err, "error processing socket read loop for connection [%s]", connID.String())
		}
		for _, x := range stats {
			x.Received(len(message))
		}
		m, err := util.FromJSONObj[*Message](message)
		if err != nil {
			return errors.Wrap(err, "error decoding websocket message")
		}
		if onMessage == nil {
			logger.Warnf("received [%s] message with command [%s] from connection [%s], which has no handler", m.Channel, m.Cmd, m.From)
		} else {
			err = onMessage(m)
			if err != nil {
				return errors.Wrap(err, "error handling websocket message")
			}
		}
	}
}
