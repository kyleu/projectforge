// Content managed by Project Forge, see [projectforge.md] for details.
package websocket

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/exp/slices"

	"projectforge.dev/projectforge/app/util"
)

type Channel struct {
	Key          string      `json:"key"`
	ConnIDs      []uuid.UUID `json:"connIDs"`
	LastUpdate   time.Time   `json:"lastUpdate"`
	MessageCount int         `json:"messageCount"`
}

func newChannel(key string) *Channel {
	return &Channel{Key: key, LastUpdate: time.Now()}
}

func (s *Service) Join(connID uuid.UUID, ch string, logger util.Logger) (bool, error) {
	conn, ok := s.connections[connID]
	if !ok {
		return false, invalidConnection(connID)
	}
	if !slices.Contains(conn.Channels, ch) {
		conn.Channels = append(conn.Channels, ch)
	}

	s.channelsMu.Lock()
	defer s.channelsMu.Unlock()

	created := false
	curr, ok := s.channels[ch]
	if !ok {
		curr = newChannel(ch)
		s.channels[ch] = curr
		created = true
	}
	if !slices.Contains(curr.ConnIDs, connID) {
		curr.ConnIDs = append(curr.ConnIDs, connID)
	}
	return created, s.sendOnlineUpdate(ch, conn.ID, conn.ID, true, logger)
}

func (s *Service) Leave(connID uuid.UUID, ch string, logger util.Logger) (bool, error) {
	conn, ok := s.connections[connID]
	if !ok {
		return false, invalidConnection(connID)
	}
	conn.Channels = chanWithout(conn.Channels, ch)

	s.channelsMu.Lock()
	defer s.channelsMu.Unlock()

	curr, ok := s.channels[ch]
	if !ok {
		curr = newChannel(ch)
	}

	filteredConns := make([]uuid.UUID, 0, len(curr.ConnIDs))
	for _, i := range curr.ConnIDs {
		if i != connID {
			filteredConns = append(filteredConns, i)
		}
	}
	if len(filteredConns) == 0 {
		delete(s.channels, ch)
		return true, nil
	}

	s.channels[ch].ConnIDs = filteredConns
	return false, s.sendOnlineUpdate(ch, conn.ID, conn.ID, false, logger)
}

func (s *Service) GetChannel(ch string) *Channel {
	ret := s.channels[ch]
	return ret
}

func (s *Service) GetConnection(id uuid.UUID) *Connection {
	ret := s.connections[id]
	return ret
}

func chanWithout(c []string, ch string) []string {
	ret := make([]string, 0, len(c))
	for _, x := range c {
		if x != ch {
			ret = append(ret, x)
		}
	}
	return ret
}
