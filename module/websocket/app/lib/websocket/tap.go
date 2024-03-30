package websocket

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

func (s *Service) RegisterTap(w http.ResponseWriter, r *http.Request, logger util.Logger) (uuid.UUID, error) {
	id := util.UUID()
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return id, err
	}

	func() {
		s.tapsMu.Lock()
		defer s.tapsMu.Unlock()
		s.taps[id] = conn
	}()
	onMessage := func(m *Message) error {
		logger.Errorf("message [%s:%s] received from tap socket", m.Channel, m.Cmd)
		return nil
	}
	onDisconnect := func() error {
		s.RemoveTap(id, logger)
		return nil
	}
	err = ReadSocketLoop(id, conn, onMessage, onDisconnect, logger)
	if err != nil {
		logger.Errorf("error registering tap socket [%s]: %s", id.String(), err.Error())
	}
	if err != nil {
		return id, errors.Wrapf(err, "error registering tap [%s]: %s", id.String(), err.Error())
	}
	return id, nil
}

func (s *Service) RemoveTap(connID uuid.UUID, logger util.Logger) {
	_, ok := s.taps[connID]
	if !ok {
		logger.Warnf("unable to find tap [%s]", connID.String())
		return
	}
	s.tapsMu.Lock()
	defer s.tapsMu.Unlock()
	delete(s.taps, connID)
}

func (s *Service) WriteTap(message *Message, logger util.Logger) {
	if len(s.taps) == 0 {
		return
	}
	b := util.ToJSONBytes(message, true)
	lo.ForEach(lo.Values(s.taps), func(tap *websocket.Conn, _ int) {
		err := tap.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			logger.Warnf("error writing tap message: %s", err.Error())
		}
	})
}
