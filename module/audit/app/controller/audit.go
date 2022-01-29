package controller

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/audit"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vaudit"
)

const auditDefaultTitle = "Audits"

var auditBreadcrumb = "Audit"

func AuditList(rc *fasthttp.RequestCtx) {
	act("audit.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ps.Title = auditDefaultTitle
		params := cutil.ParamSetFromRequest(rc)
		ret, err := as.Services.Audit.List(ps.Context, nil, params.Get("audit", nil, ps.Logger))
		if err != nil {
			return "", err
		}
		ps.Data = ret
		return render(rc, as, &vaudit.List{Models: ret, Params: params}, ps, "admin", "Audit")
	})
}

func AuditDetail(rc *fasthttp.RequestCtx) {
	act("audit.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		params := cutil.ParamSetFromRequest(rc)
		ret, err := auditFromPath(rc, as, ps)
		if err != nil {
			return "", err
		}
		ps.Title = ret.String()
		ps.Data = ret
		records, err := as.Services.Audit.RecordsForAudit(ps.Context, nil, ret.ID, params.Get("auditRecord", nil, ps.Logger))
		if err != nil {
			return "", errors.Wrap(err, "unable to retrieve child auditrecords")
		}
		return render(rc, as, &vaudit.Detail{Model: ret, Params: params, Records: records}, ps, "admin", auditBreadcrumb, ret.String())
	})
}

func AuditCreateForm(rc *fasthttp.RequestCtx) {
	act("audit.create.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := &audit.Audit{}
		ps.Title = "Create [Audit]"
		ps.Data = ret
		return render(rc, as, &vaudit.Edit{Model: ret, IsNew: true}, ps, "admin", auditBreadcrumb, "Create")
	})
}

func AuditCreateFormRandom(rc *fasthttp.RequestCtx) {
	act("audit.create.form.random", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := audit.Random()
		ps.Title = "Create Random [Audit]"
		ps.Data = ret
		return render(rc, as, &vaudit.Edit{Model: ret, IsNew: true}, ps, "admin", auditBreadcrumb, "Create")
	})
}

func AuditCreate(rc *fasthttp.RequestCtx) {
	act("audit.create", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromForm(rc, true)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse Audit from form")
		}
		err = as.Services.Audit.Create(ps.Context, nil, ret)
		if err != nil {
			return "", errors.Wrap(err, "unable to save newly-created Audit")
		}
		msg := fmt.Sprintf("Audit [%s] created", ret.String())
		return flashAndRedir(true, msg, ret.WebPath(), rc, ps)
	})
}

func AuditEditForm(rc *fasthttp.RequestCtx) {
	act("audit.edit.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromPath(rc, as, ps)
		if err != nil {
			return "", err
		}
		ps.Title = "Edit [" + ret.String() + "]"
		ps.Data = ret
		return render(rc, as, &vaudit.Edit{Model: ret}, ps, "admin", auditBreadcrumb, ret.String())
	})
}

func AuditEdit(rc *fasthttp.RequestCtx) {
	act("audit.edit", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromPath(rc, as, ps)
		if err != nil {
			return "", err
		}
		frm, err := auditFromForm(rc, false)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse Audit from form")
		}
		frm.ID = ret.ID
		err = as.Services.Audit.Update(ps.Context, nil, frm)
		if err != nil {
			return "", errors.Wrapf(err, "unable to update Audit [%s]", frm.String())
		}
		msg := fmt.Sprintf("Audit [%s] updated", frm.String())
		return flashAndRedir(true, msg, frm.WebPath(), rc, ps)
	})
}

func AuditDelete(rc *fasthttp.RequestCtx) {
	act("audit.delete", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromPath(rc, as, ps)
		if err != nil {
			return "", err
		}
		err = as.Services.Audit.Delete(ps.Context, nil, ret.ID)
		if err != nil {
			return "", errors.Wrapf(err, "unable to delete audit [%s]", ret.String())
		}
		msg := fmt.Sprintf("Audit [%s] deleted", ret.String())
		return flashAndRedir(true, msg, "/audit", rc, ps)
	})
}

func RecordDetail(rc *fasthttp.RequestCtx) {
	act("audit.record.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := recordFromPath(rc, as, ps)
		if err != nil {
			return "", err
		}
		ps.Title = ret.String()
		ps.Data = ret
		return render(rc, as, &vaudit.RecordDetail{Model: ret}, ps, "admin", auditBreadcrumb, ret.String())
	})
}

func auditFromPath(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (*audit.Audit, error) {
	idArgStr, err := RCRequiredString(rc, "id", false)
	if err != nil {
		return nil, errors.Wrap(err, "must provide [id] as an argument")
	}
	idArgP := util.UUIDFromString(idArgStr)
	if idArgP == nil {
		return nil, errors.Errorf("argument [id] (%s) is not a valid UUID", idArgStr)
	}
	idArg := *idArgP
	return as.Services.Audit.Get(ps.Context, nil, idArg)
}

func auditFromForm(rc *fasthttp.RequestCtx, setPK bool) (*audit.Audit, error) {
	frm, err := cutil.ParseForm(rc)
	if err != nil {
		return nil, errors.Wrap(err, "unable to parse form")
	}
	return audit.FromMap(frm, setPK)
}

func recordFromPath(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (*audit.Record, error) {
	idArgStr, err := RCRequiredString(rc, "id", false)
	if err != nil {
		return nil, errors.Wrap(err, "must provide [id] as an argument")
	}
	idArgP := util.UUIDFromString(idArgStr)
	if idArgP == nil {
		return nil, errors.Errorf("argument [id] (%s) is not a valid UUID", idArgStr)
	}
	idArg := *idArgP
	return as.Services.Audit.GetRecord(ps.Context, nil, idArg)
}
