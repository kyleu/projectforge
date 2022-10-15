package websocket

import (
	"time"

	"github.com/google/uuid"
)

type Channel struct {
	Key          string      `json:"key"`
	MemberIDs    []uuid.UUID `json:"memberIDs"`
	LastUpdate   time.Time   `json:"lastUpdate"`
	MessageCount int         `json:"messageCount"`
}

func newChannel(key string) *Channel {
	return &Channel{Key: key, MemberIDs: []uuid.UUID{}, LastUpdate: time.Now()}
}

// Adds a Connection to this Channel.
func (s *Service) Join(connID uuid.UUID, ch string) (bool, error) {
	conn, ok := s.connections[connID]
	if !ok {
		return false, invalidConnection(connID)
	}
	if !chanContains(conn.Channels, ch) {
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
	if !containsUUID(curr.MemberIDs, connID) {
		curr.MemberIDs = append(curr.MemberIDs, connID)
	}
	return created, nil
}

// Removes a Connection from this Channel.
func (s *Service) Leave(connID uuid.UUID, ch string) (bool, error) {
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
	filtered := make([]uuid.UUID, 0)
	for _, i := range curr.MemberIDs {
		if i != connID {
			filtered = append(filtered, i)
		}
	}

	if len(filtered) == 0 {
		delete(s.channels, ch)
		return true, nil
	}

	s.channels[ch].MemberIDs = filtered
	return false, s.sendOnlineUpdate(ch, conn.ID, conn.ID, false)
}

func containsUUID(s []uuid.UUID, e uuid.UUID) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func chanContains(c []string, id string) bool {
	for _, x := range c {
		if x == id {
			return true
		}
	}
	return false
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
