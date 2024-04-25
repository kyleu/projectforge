package audit

import (
	{{{ if .SQLServerOnly }}}{{{ else }}}{{{ if .PostgreSQL }}}"encoding/json"
	{{{ end }}}{{{ end }}}"fmt"
	"strings"
	"time"

	{{{ if .SQLServerOnly }}}mssql "github.com/denisenkom/go-mssqldb"{{{ else }}}"github.com/google/uuid"{{{ end }}}
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

{{{ if .PostgreSQL }}}type row struct {
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
}{{{ else }}}{{{ if .SQLite }}}type row struct {
	ID        uuid.UUID `db:"id"`
	App       string    `db:"app"`
	Act       string    `db:"act"`
	Client    string    `db:"client"`
	Server    string    `db:"server"`
	User      string    `db:"user"`
	Metadata  string    `db:"metadata"`
	Message   string    `db:"message"`
	Started   time.Time `db:"started"`
	Completed time.Time `db:"completed"`
}{{{ else }}}{{{ if .SQLServer }}}type row struct {
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
}{{{ else }}}type row struct {}{{{ end }}}{{{ end }}}{{{ end }}}

func (r *row) ToAudit() *Audit {
	if r == nil {
		return nil
	}{{{ if .SQLServerOnly }}}
	metadataArg, _ := util.FromJSONMap([]byte(r.Metadata))
	return &Audit{
		ID: util.UUIDFromStringOK(r.ID.String()), App: r.App, Act: r.Act, Client: r.Client, Server: r.Server, User: r.User,
		Metadata: metadataArg, Message: r.Message, Started: r.Started, Completed: r.Completed,
	}{{{ else }}}
	metadataArg, _ := util.FromJSONMap({{{ if .PostgreSQL }}}r.Metadata{{{ else }}}[]byte(r.Metadata){{{ end }}})
	return &Audit{
		ID: r.ID, App: r.App, Act: r.Act, Client: r.Client, Server: r.Server, User: r.User,
		Metadata: metadataArg, Message: r.Message, Started: r.Started, Completed: r.Completed,
	}{{{ end }}}
}

type rows []*row

func (x rows) ToAudits() Audits {
	return lo.Map(x, func(r *row, _ int) *Audit {
		return r.ToAudit()
	})
}
