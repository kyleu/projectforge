package audit

import (
	{{{ if .SQLServer }}}{{{ else }}}"encoding/json"
	{{{ end }}}"fmt"
	"strings"
	"time"

	{{{ if .SQLServer }}}mssql "github.com/denisenkom/go-mssqldb"{{{ else }}}"github.com/google/uuid"{{{ end }}}
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

var (
	table         = "audit"
	tableQuoted   = fmt.Sprintf("%q", table)
	columns       = []string{"id", "app", "act", "client", "server", "user", "metadata", "message", "started", "completed"}
	columnsQuoted = util.StringArrayQuoted(columns)
	columnsString = strings.Join(columnsQuoted, ", ")
	defaultWC     = "\"id\" = {{{ .Placeholder 1 }}}"
)

{{{ if .SQLServer }}}type row struct {
	ID        mssql.UniqueIdentifier `db:"id"`
	App       string                 `db:"app"`
	Act       string                 `db:"act"`
	Client    string                 `db:"client"`
	Server    string                 `db:"server"`
	User      string                 `db:"user"`
	Metadata  string                 `db:"metadata"`
	Message   string                 `db:"message"`
	Started   time.Time              `db:"started"`
	Completed time.Time              `db:"completed"`
}{{{ else }}}type row struct {
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
}{{{ end }}}

func (r *row) ToAudit() *Audit {
	if r == nil {
		return nil
	}
	metadataArg := util.ValueMap{}
	_ = util.FromJSON({{{ if .SQLServer }}}[]byte(r.Metadata){{{ else }}}r.Metadata{{{ end }}}, &metadataArg)
	return &Audit{
		ID: {{{ if .SQLServer }}}util.UUIDFromStringOK(r.ID.String()){{{ else }}}r.ID{{{ end }}}, App: r.App, Act: r.Act, Client: r.Client, Server: r.Server, User: r.User,
		Metadata: metadataArg, Message: r.Message, Started: r.Started, Completed: r.Completed,
	}
}

type rows []*row

func (x rows) ToAudits() Audits {
	return lo.Map(x, func(r *row, _ int) *Audit {
		return r.ToAudit()
	})
}
