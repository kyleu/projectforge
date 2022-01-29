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
	recordColumns       = []string{"id", "audit_id", "t", "pk", "changes", "occurred"}
	recordColumnsQuoted = util.StringArrayQuoted(recordColumns)
	recordColumnsString = strings.Join(recordColumnsQuoted, ", ")
	recordDefaultWC     = "\"id\" = $1"
)

type recordDTO struct {
	ID       uuid.UUID       `db:"id"`
	AuditID  uuid.UUID       `db:"audit_id"`
	T        string          `db:"t"`
	Pk       string          `db:"pk"`
	Changes  json.RawMessage `db:"changes"`
	Occurred time.Time       `db:"occurred"`
}

func (d *recordDTO) ToRecord() *Record {
	if d == nil {
		return nil
	}
	changesArg := util.ValueMap{}
	_ = util.FromJSON(d.Changes, &changesArg)
	return &Record{ID: d.ID, AuditID: d.AuditID, T: d.T, Pk: d.Pk, Changes: changesArg, Occurred: d.Occurred}
}

type recordDTOs []*recordDTO

func (x recordDTOs) ToRecords() Records {
	ret := make(Records, 0, len(x))
	for _, d := range x {
		ret = append(ret, d.ToRecord())
	}
	return ret
}
