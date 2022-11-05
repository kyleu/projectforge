// Content managed by Project Forge, see [projectforge.md] for details.
package websocket

import (
	"github.com/google/uuid"
	"golang.org/x/exp/slices"
)

type UpdateMemberParams struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
}

func (s *Service) GetAllMembers(key string) []*Connection {
	ch, ok := s.channels[key]
	if !ok {
		return nil
	}
	ret := make([]*Connection, 0, len(ch.MemberIDs))
	for _, cID := range ch.MemberIDs {
		c, ok := s.connections[cID]
		if ok && c != nil {
			ret = append(ret, c)
		}
	}
	return ret
}

func (s *Service) GetOnline(key string) []uuid.UUID {
	ch, ok := s.channels[key]
	if !ok {
		return nil
	}
	online := make([]uuid.UUID, 0)
	for _, cID := range ch.MemberIDs {
		c, ok := s.connections[cID]
		if ok && c != nil && (!slices.Contains(online, c.ID)) {
			online = append(online, c.ID)
		}
	}
	return online
}

func (s *Service) sendOnlineUpdate(ch string, connID uuid.UUID, userID uuid.UUID, connected bool) error {
	p := OnlineUpdate{UserID: userID, Connected: connected}
	onlineMsg := NewMessage(&userID, ch, "online-update", p)
	return s.WriteChannel(onlineMsg, connID)
}
