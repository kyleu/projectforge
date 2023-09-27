package websocket

import (
	"github.com/google/uuid"
)

type Status struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Channels []string  `json:"channels"`
}

type Statuses = []*Status
