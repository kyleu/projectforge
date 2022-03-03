package websocket

import (
	"github.com/google/uuid"
)

type UpdateMemberParams struct {
	ID   uuid.UUID `json:"id"`
	Role string    `json:"role"`
}

// Returns membership details for the provided Channel.
func (s *Service) GetOnline(ch string) []uuid.UUID {
	connections, ok := s.channels[ch]
	if !ok {
		connections = newChannel(ch)
	}
	online := make([]uuid.UUID, 0)
	for _, cID := range connections.MemberIDs {
		c, ok := s.connections[cID]
		if ok && c != nil && (!containsUUID(online, c.ID)) {
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
