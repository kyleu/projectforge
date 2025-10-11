package websocket

import (
	"fmt"
	"time"

	"{{{ .Package }}}/app/util"
)

type Stats struct {
	MessagesSent     int       `json:"messagesSent"`
	MessagesReceived int       `json:"messagesReceived"`
	BytesSent        int       `json:"bytesSent"`
	BytesReceived    int       `json:"bytesReceived"`
	Started          time.Time `json:"started,omitempty"`
}

func NewStats() *Stats {
	return &Stats{Started: time.Now()}
}

func (s *Stats) Sent(b int) {
	s.MessagesSent++
	s.BytesSent += b
}

func (s *Stats) Received(b int) {
	s.MessagesReceived++
	s.BytesReceived += b
}

func (s *Stats) String() string {
	bs, br := util.ByteSizeSI(int64(s.BytesSent)), util.ByteSizeSI(int64(s.BytesReceived))
	return fmt.Sprintf("[%d msgs/%s] sent and [%d msgs/%s] received", s.MessagesSent, bs, s.MessagesReceived, br)
}

func (s *Stats) Clone() *Stats {
	return &Stats{MessagesSent: s.MessagesSent, MessagesReceived: s.MessagesReceived, BytesSent: s.BytesSent, BytesReceived: s.BytesReceived, Started: s.Started}
}
