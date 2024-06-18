package websocket

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/robert-nix/ansihtml"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
)

type ConnectEvent func(s *Service, conn *Connection, logger util.Logger) error

type Service struct {
	connections   map[uuid.UUID]*Connection
	connectionsMu sync.Mutex
	channels      map[string]*Channel
	channelsMu    sync.Mutex
	taps          map[uuid.UUID]*websocket.Conn
	tapsMu        sync.Mutex
	onOpen        ConnectEvent
	onClose       ConnectEvent
}

func NewService(onOpen ConnectEvent, onClose ConnectEvent) *Service {
	return &Service{
		connections: make(map[uuid.UUID]*Connection),
		channels:    make(map[string]*Channel),
		taps:        make(map[uuid.UUID]*websocket.Conn),
		onOpen:      onOpen,
		onClose:     onClose,
	}
}

var (
	systemID     = *util.UUIDFromString("FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF")
	systemStatus = &Status{
		ID: systemID, Username: util.GetEnv("system_username", "System Broadcast"), Channels: []string{systemID.String()},
	}
)

func (s *Service) ReplaceHandlers(onOpen ConnectEvent, onClose ConnectEvent) {
	s.onOpen = onOpen
	s.onClose = onClose
}

func (s *Service) Close() {
	s.connectionsMu.Lock()
	defer s.connectionsMu.Unlock()
	lo.ForEach(lo.Values(s.connections), func(v *Connection, _ int) {
		_ = v.Close()
	})
}

var upgrader = websocket.Upgrader{EnableCompression: true}

func (s *Service) Upgrade(
	ctx context.Context, w http.ResponseWriter, r *http.Request, channel string,
	profile *user.Profile, handler Handler, logger util.Logger,
) (uuid.UUID, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return uuid.Nil, err
	}
	cx, err := s.Register(profile, conn, handler, logger)
	if err != nil {
		logger.Warnf("unable to register websocket connection: %+v", err)
		return uuid.Nil, nil
	}
	joined, err := s.Join(cx.ID, channel, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("error processing socket join (%v): %+v", joined, err))
		return uuid.Nil, nil
	}
	return cx.ID, nil
}

func (s *Service) Terminal(ch string, logger util.Logger) func(_ string, b []byte) error {
	return func(_ string, b []byte) error {
		html := string(ansihtml.ConvertToHTML(b))
		m := util.ValueMap{"msg": string(b), "html": strings.TrimSpace(html)}
		msg := NewMessage(nil, ch, "output", m)
		return s.WriteChannel(msg, logger)
	}
}
