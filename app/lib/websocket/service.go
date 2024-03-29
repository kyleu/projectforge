// Package websocket - Content managed by Project Forge, see [projectforge.md] for details.
package websocket

import (
	"cmp"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"sync"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/lib/user"
	"projectforge.dev/projectforge/app/util"
)

type Handler func(ctx context.Context, s *Service, conn *Connection, svc string, cmd string, param json.RawMessage, logger util.Logger) error

type ConnectEvent func(s *Service, conn *Connection, logger util.Logger) error

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

func (s *Service) ChannelList(params *filter.Params) []string {
	params = filter.ParamsWithDefaultOrdering("channel", params)
	ret := &util.StringSlice{}
	idx := 0
	lo.ForEach(lo.Keys(s.channels), func(conn string, _ int) {
		if idx >= params.Offset && (params.Limit == 0 || idx < params.Limit) {
			ret.Push(conn)
		}
		idx++
	})
	return util.ArraySorted(ret.Slice)
}

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

func (s *Service) Count() int {
	return len(s.connections)
}

func (s *Service) Status() ([]string, []*Connection, []uuid.UUID) {
	s.connectionsMu.Lock()
	defer s.connectionsMu.Unlock()
	conns := lo.Values(s.connections)
	slices.SortFunc(conns, func(l *Connection, r *Connection) int {
		return cmp.Compare(l.ID.String(), r.ID.String())
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

var upgrader = websocket.Upgrader{EnableCompression: true}

func (s *Service) Upgrade(
	ctx context.Context, w http.ResponseWriter, r *http.Request, channel string, profile *user.Profile, logger util.Logger,
) error {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return err
	}
	cx, err := s.Register(profile, conn, logger)
	if err != nil {
		logger.Warnf("unable to register websocket connection: %+v", err)
		return nil
	}
	joined, err := s.Join(cx.ID, channel, logger)
	if err != nil {
		logger.Error(fmt.Sprintf("error processing socket join (%v): %+v", joined, err))
		return nil
	}
	err = s.ReadLoop(ctx, cx.ID, logger)
	if err != nil {
		if !strings.Contains(err.Error(), "1001") {
			logger.Error(fmt.Sprintf("error processing socket read loop for connection [%s]: %+v", cx.ID.String(), err))
		}
		return nil
	}
	return nil
}
