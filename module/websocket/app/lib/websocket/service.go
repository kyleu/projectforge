package websocket

import (
	"encoding/json"
	"fmt"
	"sort"
	"sync"

	"github.com/fevo-tech/nuevo/app/lib/filter"
	"github.com/fevo-tech/nuevo/app/util"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Function used to handle incoming messages.
type Handler func(s *Service, conn *Connection, svc string, cmd string, param json.RawMessage) error

// Function used to handle incoming connections.
type ConnectEvent func(s *Service, conn *Connection) error

// Manages all Connection objects.
type Service struct {
	connections   map[uuid.UUID]*Connection
	connectionsMu sync.Mutex
	channels      map[string]*Channel
	channelsMu    sync.Mutex
	Logger        *zap.SugaredLogger
	onOpen        ConnectEvent
	handler       Handler
	onClose       ConnectEvent
	Context       any
}

// Creates a new service with the provided handler functions.
func NewService(logger *zap.SugaredLogger, onOpen ConnectEvent, handler Handler, onClose ConnectEvent, ctx any) *Service {
	return &Service{
		connections: make(map[uuid.UUID]*Connection),
		channels:    make(map[string]*Channel),
		Logger:      logger,
		handler:     handler,
		onOpen:      onOpen,
		Context:     ctx,
	}
}

var (
	systemID     = *util.UUIDFromString("FFFFFFFF-FFFF-FFFF-FFFF-FFFFFFFFFFFF")
	systemStatus = &Status{ID: systemID, Username: "System Broadcast", Channels: []string{systemID.String()}}
)

// Returns an array of Connection statuses based on the parameters.
func (s *Service) UserList(params *filter.Params) Statuses {
	params = filter.ParamsWithDefaultOrdering("connection", params)
	ret := make(Statuses, 0)
	ret = append(ret, systemStatus)
	idx := 0
	for _, conn := range s.connections {
		if idx >= params.Offset && (params.Limit == 0 || idx < params.Limit) {
			ret = append(ret, conn.ToStatus())
		}
		idx++
	}
	return ret
}

// Returns an array of channels based on the parameters.
func (s *Service) ChannelList(params *filter.Params) []string {
	params = filter.ParamsWithDefaultOrdering("channel", params)
	ret := make([]string, 0)
	idx := 0
	for conn := range s.channels {
		if idx >= params.Offset && (params.Limit == 0 || idx < params.Limit) {
			ret = append(ret, conn)
		}
		idx++
	}
	sort.Strings(ret)
	return ret
}

// Returns a Status representing the Connection with the provided ID.
func (s *Service) GetByID(id uuid.UUID) *Status {
	if id == systemID {
		return systemStatus
	}
	conn, ok := s.connections[id]
	if !ok {
		s.Logger.Error(fmt.Sprintf("error getting connection by id [%v]", id))
		return nil
	}
	return conn.ToStatus()
}

// Total number of all connections managed by this service.
func (s *Service) Count() int {
	return len(s.connections)
}

// Callback for when the backing connection is re-established.
func (s *Service) OnOpen(connID uuid.UUID) error {
	c, ok := s.connections[connID]
	if !ok {
		return invalidConnection(connID)
	}
	return s.onOpen(s, c)
}

// Sends a message to a provided Connection ID.
func OnMessage(s *Service, connID uuid.UUID, message *Message) error {
	if connID == systemID {
		s.Logger.Warn("--- admin message received ---")
		s.Logger.Warn(message.String())
		return nil
	}
	s.connectionsMu.Lock()
	c, ok := s.connections[connID]
	s.connectionsMu.Unlock()
	if !ok {
		return invalidConnection(connID)
	}

	return s.handler(s, c, message.Channel, message.Cmd, message.Param)
}

// Callback for when the backing connection is closed.
func (s *Service) OnClose(connID uuid.UUID) error {
	c, ok := s.connections[connID]
	if !ok {
		return invalidConnection(connID)
	}
	if s.onClose != nil {
		return s.onClose(s, c)
	}
	return nil
}
