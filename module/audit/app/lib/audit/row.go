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

type row struct {
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

func (r *row) ToAudit() *Audit {
	if r == nil {
		return nil
	}
	metadataArg := util.ValueMap{}
	_ = util.FromJSON(r.Metadata, &metadataArg)
	return &Audit{
		ID: r.ID, App: r.App, Act: r.Act, Client: r.Client, Server: r.Server, User: r.User,
		Metadata: metadataArg, Message: r.Message, Started: r.Started, Completed: r.Completed,
	}
}

type rows []*row

func (x rows) ToAudits() Audits {
	ret := make(Audits, 0, len(x))
	for _, r := range x {
		ret = append(ret, r.ToAudit())
	}
	return ret
}
