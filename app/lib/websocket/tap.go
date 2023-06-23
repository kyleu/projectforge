// Content managed by Project Forge, see [projectforge.md] for details.
package websocket

import (
	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"

	"projectforge.dev/projectforge/app/util"
)

func (s *Service) RegisterTap(rc *fasthttp.RequestCtx, logger util.Logger) (uuid.UUID, error) {
	id := util.UUID()
	err := upgrader.Upgrade(rc, func(conn *websocket.Conn) {
		s.tapsMu.Lock()
		defer s.tapsMu.Unlock()
		s.taps[id] = conn
		onMessage := func(m *Message) error {
			logger.Errorf("message received from tap socket")
			return nil
		}
		onDisconnect := func() error {
			s.RemoveTap(id, logger)
			return nil
		}
		err := ReadSocketLoop(id, conn, onMessage, onDisconnect, logger)
		if err != nil {
			logger.Errorf("error registering tap socket [%s]: %s", id.String(), err.Error())
		}
	})
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
