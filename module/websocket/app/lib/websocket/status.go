package websocket

import (
	"cmp"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/lib/filter"
	"{{{ .Package }}}/app/util"
)

type Status struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Channels []string  `json:"channels"`
}

type Statuses = []*Status

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
