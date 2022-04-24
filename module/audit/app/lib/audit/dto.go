package audit

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"{{{ .Package }}}/app/util"
)

var (
	table         = "audit"
	tableQuoted   = fmt.Sprintf("%q", table)
	columns       = []string{"id", "app", "act", "client", "server", "user", "metadata", "message", "started", "completed"}
	columnsQuoted = util.StringArrayQuoted(columns)
	columnsString = strings.Join(columnsQuoted, ", ")
	defaultWC     = "\"id\" = $1"
)

type dto struct {
	ID        uuid.UUID       `db:"id"`
	App       string          `db:"app"`
	Act       string          `db:"act"`
	Client    string          `db:"client"`
	Server    string          `db:"server"`
	User      string          `db:"user"`
	Metadata  json.RawMessage `db:"metadata"`
	Message   string          `db:"message"`
	Started   time.Time       `db:"started"`
	Completed time.Time       `db:"completed"`
}

func (d *dto) ToAudit() *Audit {
	if d == nil {
		return nil
	}
	metadataArg := util.ValueMap{}
	_ = util.FromJSON(d.Metadata, &metadataArg)
	return &Audit{
		ID: d.ID, App: d.App, Act: d.Act, Client: d.Client, Server: d.Server, User: d.User,
		Metadata: metadataArg, Message: d.Message, Started: d.Started, Completed: d.Completed,
	}
}

type dtos []*dto

func (x dtos) ToAudits() Audits {
	ret := make(Audits, 0, len(x))
	for _, d := range x {
		ret = append(ret, d.ToAudit())
	}
	return ret
}
