package template

func (t *Context) AuditData() string {
	if t.SQLServer() {
		return "[]any{a.ID.String(), a.App, a.Act, a.Client, a.Server, a.User, util.ToJSON(a.Metadata), a.Message, a.Started, a.Completed}"
	}
	if t.SQLite() {
		return "[]any{a.ID, a.App, a.Act, a.Client, a.Server, a.User, util.ToJSON(a.Metadata), a.Message, a.Started, a.Completed}"
	}
	return "[]any{a.ID, a.App, a.Act, a.Client, a.Server, a.User, a.Metadata, a.Message, a.Started, a.Completed}"
}

func (t *Context) AuditRecordData() string {
	if t.SQLServer() {
		return "[]any{a.ID.String(), a.AuditID.String(), a.T, a.PK, util.ToJSON(a.Changes), util.ToJSON(a.Metadata), a.Occurred}"
	}
	if t.SQLite() {
		return "[]any{a.ID, a.AuditID, a.T, a.PK, util.ToJSON(a.Changes), util.ToJSON(a.Metadata), a.Occurred}"
	}
	return "[]any{a.ID, a.AuditID, a.T, a.PK, a.Changes, a.Metadata, a.Occurred}"
}
