package audit

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"{{{ .Package }}}/app/util"
)

type Record struct {
	ID       uuid.UUID     `json:"id"`
	AuditID  uuid.UUID     `json:"auditID"`
	T        string        `json:"t"`
	Pk       string        `json:"pk"`
	Changes  util.ValueMap `json:"changes"`
	Occurred time.Time     `json:"occurred"`
}

func NewRecord(id uuid.UUID) *Record {
	return &Record{ID: id}
}

func RandomRecord() *Record {
	return &Record{
		ID:       util.UUID(),
		AuditID:  util.UUID(),
		T:        util.RandomString(12),
		Pk:       util.RandomString(12),
		Changes:  util.RandomValueMap(4),
		Occurred: time.Now(),
	}
}

func RecordFromMap(m util.ValueMap, setPK bool) (*Record, error) {
	ret := &Record{}
	var err error
	if setPK {
		retID, e := m.ParseUUID("id", true, true)
		if e != nil {
			return nil, e
		}
		if retID != nil {
			ret.ID = *retID
		}
		// $PF_SECTION_START(pkchecks)$
		// $PF_SECTION_END(pkchecks)$
	}
	retAuditID, e := m.ParseUUID("auditID", true, true)
	if e != nil {
		return nil, e
	}
	if retAuditID != nil {
		ret.AuditID = *retAuditID
	}
	ret.T, err = m.ParseString("t", true, true)
	if err != nil {
		return nil, err
	}
	ret.Pk, err = m.ParseString("pk", true, true)
	if err != nil {
		return nil, err
	}
	ret.Changes, err = m.ParseMap("changes", true, true)
	if err != nil {
		return nil, err
	}
	retOccurred, e := m.ParseTime("occurred", true, true)
	if e != nil {
		return nil, e
	}
	if retOccurred != nil {
		ret.Occurred = *retOccurred
	}
	// $PF_SECTION_START(extrachecks)$
	// $PF_SECTION_END(extrachecks)$
	return ret, nil
}

func (a *Record) Clone() *Record {
	return &Record{
		ID:       a.ID,
		AuditID:  a.AuditID,
		T:        a.T,
		Pk:       a.Pk,
		Changes:  a.Changes,
		Occurred: a.Occurred,
	}
}

func (a *Record) String() string {
	return a.ID.String()
}

func (a *Record) WebPath() string {
	return "/admin/audit/record" + "/" + a.ID.String()
}

func (a *Record) Diff(ax *Record) util.Diffs {
	var diffs util.Diffs
	if a.ID != ax.ID {
		diffs = append(diffs, util.NewDiff("id", a.ID.String(), ax.ID.String()))
	}
	if a.AuditID != ax.AuditID {
		diffs = append(diffs, util.NewDiff("auditID", a.AuditID.String(), ax.AuditID.String()))
	}
	if a.T != ax.T {
		diffs = append(diffs, util.NewDiff("t", a.T, ax.T))
	}
	if a.Pk != ax.Pk {
		diffs = append(diffs, util.NewDiff("pk", a.Pk, ax.Pk))
	}
	diffs = append(diffs, util.DiffObjects(a.Changes, ax.Changes, "changes")...)
	if a.Occurred != ax.Occurred {
		diffs = append(diffs, util.NewDiff("occurred", fmt.Sprint(a.Occurred), fmt.Sprint(ax.Occurred)))
	}
	return diffs
}

func (a *Record) ToData() []interface{} {
	return []interface{}{a.ID, a.AuditID, a.T, a.Pk, a.Changes, a.Occurred}
}

type Records []*Record
