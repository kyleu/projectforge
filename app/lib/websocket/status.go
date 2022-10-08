// Content managed by Project Forge, see [projectforge.md] for details.
package websocket

import (
	"github.com/google/uuid"
)

// Serializable representation of a Connection.
type Status struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Channels []string  `json:"channels"`
}

type Statuses = []*Status
