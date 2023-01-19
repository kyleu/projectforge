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
	recordTable         = "audit_record"
	recordTableQuoted   = fmt.Sprintf("%q", recordTable)
	recordColumns       = []string{"id", "audit_id", "t", "pk", "changes", "metadata", "occurred"}
	recordColumnsQuoted = util.StringArrayQuoted(recordColumns)
	recordColumnsString = strings.Join(recordColumnsQuoted, ", ")
)

type recordRow struct {
	ID       uuid.UUID       `db:"id"`
	AuditID  uuid.UUID       `db:"audit_id"`
	T        string          `db:"t"`
	PK       string          `db:"pk"`
	Changes  json.RawMessage `db:"changes"`
	Metadata json.RawMessage `db:"metadata"`
	Occurred time.Time       `db:"occurred"`
}

func (r *recordRow) ToRecord() *Record {
	if r == nil {
		return nil
	}
	changesArg := util.Diffs{}
	_ = util.FromJSON(r.Changes, &changesArg)
	metadataArg := util.ValueMap{}
	_ = util.FromJSON(r.Metadata, &metadataArg)
	return &Record{ID: r.ID, AuditID: r.AuditID, T: r.T, PK: r.PK, Changes: changesArg, Metadata: metadataArg, Occurred: r.Occurred}
}

type recordRows []*recordRow

func (x recordRows) ToRecords() Records {
	ret := make(Records, 0, len(x))
	for _, r := range x {
		ret = append(ret, r.ToRecord())
	}
	return ret
}
