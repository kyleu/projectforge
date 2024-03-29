package schedule

import (
	"time"

	"github.com/google/uuid"
)

type Result struct {
	Job           uuid.UUID `json:"job"`
	Returned      any       `json:"returned,omitempty"`
	Error         string    `json:"error,omitempty"`
	Occurred      time.Time `json:"occurred"`
	DurationMicro int       `json:"durationMicro"`
}
