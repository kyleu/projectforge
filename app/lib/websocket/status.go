package websocket

import (
	"cmp"
	"slices"

	"github.com/google/uuid"
	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/filter"
	"projectforge.dev/projectforge/app/util"
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
	taps := util.ArrayCopy(util.MapKeys(s.taps))
	return s.ChannelList(nil), conns, taps
}

func (s *Service) UserList(params *filter.Params) Statuses {
	params = filter.ParamsWithDefaultOrdering("connection", params)
	ret := make(Statuses, 0)
	ret = append(ret, systemStatus)
	var idx int
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
	var idx int
	lo.ForEach(util.MapKeysSorted(s.channels), func(conn string, _ int) {
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
		logger.Errorf("error getting connection by id [%v]", id)
		return nil
	}
	return conn.ToStatus()
}

func (s *Service) Count() int {
	return len(s.connections)
}
