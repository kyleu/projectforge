package audit

import (
	{{{ if .SQLServer }}}{{{ else }}}{{{ if .SQLite }}}{{{ else }}}"encoding/json"
	{{{ end }}}{{{ end }}}"fmt"
	"strings"
	"time"

	{{{ if .SQLServer }}}mssql "github.com/denisenkom/go-mssqldb"{{{ else }}}"github.com/google/uuid"{{{ end }}}
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

var (
	recordTable         = "audit_record"
	recordTableQuoted   = fmt.Sprintf("%q", recordTable)
	recordColumns       = []string{"id", "audit_id", "t", "pk", "changes", "metadata", "occurred"}
	recordColumnsQuoted = util.StringArrayQuoted(recordColumns)
	recordColumnsString = strings.Join(recordColumnsQuoted, ", ")
)

{{{ if .SQLServer }}}type recordRow struct {
	ID       mssql.UniqueIdentifier `db:"id"`
	AuditID  mssql.UniqueIdentifier `db:"audit_id"`
	T        string                 `db:"t"`
	PK       string                 `db:"pk"`
	Changes  string                 `db:"changes"`
	Metadata string                 `db:"metadata"`
	Occurred time.Time              `db:"occurred"`
}{{{ else }}}{{{ if .SQLite }}}type recordRow struct {
	ID       uuid.UUID `db:"id"`
	AuditID  uuid.UUID `db:"audit_id"`
	T        string    `db:"t"`
	PK       string    `db:"pk"`
	Changes  string    `db:"changes"`
	Metadata string    `db:"metadata"`
	Occurred time.Time `db:"occurred"`
}{{{ else }}}type recordRow struct {
	ID       uuid.UUID       `db:"id"`
	AuditID  uuid.UUID       `db:"audit_id"`
	T        string          `db:"t"`
	PK       string          `db:"pk"`
	Changes  json.RawMessage `db:"changes"`
	Metadata json.RawMessage `db:"metadata"`
	Occurred time.Time       `db:"occurred"`
}{{{ end }}}{{{ end }}}

func (r *recordRow) ToRecord() *Record {
	if r == nil {
		return nil
	}
	changesArg := util.Diffs{}{{{ if .SQLServer }}}
	_ = util.FromJSON([]byte(r.Changes), &changesArg)
	metadataArg := util.ValueMap{}
	_ = util.FromJSON([]byte(r.Metadata), &metadataArg)
	return &Record{
		ID: util.UUIDFromStringOK(r.ID.String()), AuditID: util.UUIDFromStringOK(r.AuditID.String()), T: r.T, PK: r.PK,
		Changes: changesArg, Metadata: metadataArg, Occurred: r.Occurred,
	}{{{ else }}}{{{ if .SQLite }}}
	_ = util.FromJSON([]byte(r.Changes), &changesArg)
	metadataArg := util.ValueMap{}
	_ = util.FromJSON([]byte(r.Metadata), &metadataArg)
	return &Record{ID: r.ID, AuditID: r.AuditID, T: r.T, PK: r.PK, Changes: changesArg, Metadata: metadataArg, Occurred: r.Occurred}{{{ else }}}
	_ = util.FromJSON(r.Changes, &changesArg)
	metadataArg := util.ValueMap{}
	_ = util.FromJSON(r.Metadata, &metadataArg)
	return &Record{ID: r.ID, AuditID: r.AuditID, T: r.T, PK: r.PK, Changes: changesArg, Metadata: metadataArg, Occurred: r.Occurred}{{{ end }}}{{{ end }}}
}

type recordRows []*recordRow

func (x recordRows) ToRecords() Records {
	return lo.Map(x, func(r *recordRow, _ int) *Record {
		return r.ToRecord()
	})
}
