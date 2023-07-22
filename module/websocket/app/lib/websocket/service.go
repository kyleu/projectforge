package websocket

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	"github.com/fasthttp/websocket"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/valyala/fasthttp"
	"golang.org/x/exp/slices"

	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/lib/user"{{{ if .HasModule "user" }}}
	dbuser "{{{ .Package }}}/app/user"{{{ end }}}
	"{{{ .Package }}}/app/util"
)

// Function used to handle incoming messages.
type Handler func(ctx context.Context, s *Service, conn *Connection, svc string, cmd string, param json.RawMessage, logger util.Logger) error

// Function used to handle incoming connections.
type ConnectEvent func(s *Service, conn *Connection, logger util.Logger) error

// Manages all Connection objects.
type Service struct {
	connections   map[uuid.UUID]*Connection
	connectionsMu sync.Mutex
	channels      map[string]*Channel
	channelsMu    sync.Mutex
	taps          map[uuid.UUID]*websocket.Conn
	tapsMu        sync.Mutex
	onOpen        ConnectEvent
	handler       Handler
	onClose       ConnectEvent
}

// Creates a new service with the provided handler functions.
func NewService(onOpen ConnectEvent, handler Handler, onClose ConnectEvent) *Service {
	return &Service{
		connections: make(map[uuid.UUID]*Connection),
		channels:    make(map[string]*Channel),
		taps:        make(map[uuid.UUID]*websocket.Conn),
		onOpen:      onOpen,
		handler:     handler,
		onClose:     onClose,
	}
}

var (
	systemID     = *util.UUIDFromString("FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF")
	systemStatus = &Status{
		ID: systemID, Username: util.GetEnv("system_username", "System Broadcast"), Channels: []string{systemID.String()},
	}
)

// Returns an array of Connection statuses based on the parameters.
func (s *Service) UserList(params *filter.Params) Statuses {
	params = filter.ParamsWithDefaultOrdering("connection", params)
	ret := make(Statuses, 0)
	ret = append(ret, systemStatus)
	idx := 0
	lo.ForEach(lo.Values(s.connections), func(conn *Connection, _ int) {
		if idx >= params.Offset && (params.Limit == 0 || idx < params.Limit) {
			ret = append(ret, conn.ToStatus())
		}
		idx++
	})
	return ret
}

// Returns an array of channels based on the parameters.
func (s *Service) ChannelList(params *filter.Params) []string {
	params = filter.ParamsWithDefaultOrdering("channel", params)
	ret := make([]string, 0)
	idx := 0
	lo.ForEach(lo.Keys(s.channels), func(conn string, _ int) {
		if idx >= params.Offset && (params.Limit == 0 || idx < params.Limit) {
			ret = append(ret, conn)
		}
		idx++
	})
	return util.ArraySorted(ret)
}

// Returns a Status representing the Connection with the provided ID.
func (s *Service) GetByID(id uuid.UUID, logger util.Logger) *Status {
	if id == systemID {
		return systemStatus
	}
	conn, ok := s.connections[id]
	if !ok {
		logger.Error(fmt.Sprintf("error getting connection by id [%v]", id))
		return nil
	}
	return conn.ToStatus()
}

// Total number of all connections managed by this service.
func (s *Service) Count() int {
	return len(s.connections)
}

func (s *Service) Status() ([]string, []*Connection, []uuid.UUID) {
	s.connectionsMu.Lock()
	defer s.connectionsMu.Unlock()
	conns := lo.Values(s.connections)
	slices.SortFunc(conns, func(l *Connection, r *Connection) bool {
		return l.ID.String() < r.ID.String()
	})
	taps := slices.Clone(lo.Keys(s.taps))
	return s.ChannelList(nil), conns, taps
}

func (s *Service) Close() {
	s.connectionsMu.Lock()
	defer s.connectionsMu.Unlock()
	lo.ForEach(lo.Values(s.connections), func(v *Connection, _ int) {
		_ = v.Close()
	})
}

var upgrader = websocket.FastHTTPUpgrader{EnableCompression: true}

func (s *Service) Upgrade(
	ctx context.Context, rc *fasthttp.RequestCtx, channel string{{{ if .HasModule "user" }}}, u *dbuser.User{{{ end }}}, profile *user.Profile{{{ if .HasModule "oauth" }}}, accts user.Accounts{{{ end }}}, logger util.Logger,
) error {
	return upgrader.Upgrade(rc, func(conn *websocket.Conn) {
		cx, err := s.Register({{{ if .HasModule "user" }}}u, {{{ end }}}profile{{{ if .HasModule "oauth" }}}, accts{{{ end }}}, conn, logger)
		if err != nil {
			logger.Warn("unable to register websocket connection")
			return
		}
		joined, err := s.Join(cx.ID, channel, logger)
		if err != nil {
			logger.Error(fmt.Sprintf("error processing socket join (%v): %+v", joined, err))
			return
		}
		err = s.ReadLoop(ctx, cx.ID, logger)
		if err != nil {
			if !strings.Contains(err.Error(), "1001") {
				logger.Error(fmt.Sprintf("error processing socket read loop for connection [%s]: %+v", cx.ID.String(), err))
			}
			return
		}
	})
}
