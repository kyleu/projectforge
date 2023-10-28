package clib

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/audit"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vaudit"
)

const auditDefaultTitle = "Audits"

var auditBreadcrumb = "Audit||/admin/audit"

func AuditList(rc *fasthttp.RequestCtx) {
	controller.Act("audit.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := as.Services.Audit.List(ps.Context, nil, ps.Params.Get("audit", nil, ps.Logger), ps.Logger)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(auditDefaultTitle, ret)
		return controller.Render(rc, as, &vaudit.List{Models: ret, Params: ps.Params}, ps, "admin", "Audit")
	})
}

func AuditDetail(rc *fasthttp.RequestCtx) {
	controller.Act("audit.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromPath(rc, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(ret.String(), ret)
		records, err := as.Services.Audit.RecordsForAudits(ps.Context, nil, ps.Params.Get("auditRecord", nil, ps.Logger), ps.Logger, ret.ID)
		if err != nil {
			return "", errors.Wrap(err, "unable to retrieve child auditrecords")
		}
		return controller.Render(rc, as, &vaudit.Detail{Model: ret, Params: ps.Params, Records: records}, ps, "admin", auditBreadcrumb, ret.String())
	})
}

func AuditCreateForm(rc *fasthttp.RequestCtx) {
	controller.Act("audit.create.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := &audit.Audit{}
		ps.SetTitleAndData("Create [Audit]", ret)
		return controller.Render(rc, as, &vaudit.Edit{Model: ret, IsNew: true}, ps, "admin", auditBreadcrumb, "Create")
	})
}

func AuditCreateFormRandom(rc *fasthttp.RequestCtx) {
	controller.Act("audit.create.form.random", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := audit.Random()
		ps.SetTitleAndData("Create Random [Audit]", ret)
		return controller.Render(rc, as, &vaudit.Edit{Model: ret, IsNew: true}, ps, "admin", auditBreadcrumb, "Create")
	})
}

func AuditCreate(rc *fasthttp.RequestCtx) {
	controller.Act("audit.create", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromForm(rc, true)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse Audit from form")
		}
		err = as.Services.Audit.Create(ps.Context, nil, ps.Logger, ret)
		if err != nil {
			return "", errors.Wrap(err, "unable to save newly-created Audit")
		}
		msg := fmt.Sprintf("Audit [%s] created", ret.String())
		return controller.FlashAndRedir(true, msg, ret.WebPath(), rc, ps)
	})
}

func AuditEditForm(rc *fasthttp.RequestCtx) {
	controller.Act("audit.edit.form", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromPath(rc, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData("Edit ["+ret.String()+"]", ret)
		return controller.Render(rc, as, &vaudit.Edit{Model: ret}, ps, "admin", auditBreadcrumb, ret.String())
	})
}

func AuditEdit(rc *fasthttp.RequestCtx) {
	controller.Act("audit.edit", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromPath(rc, as, ps)
		if err != nil {
			return "", err
		}
		frm, err := auditFromForm(rc, false)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse Audit from form")
		}
		frm.ID = ret.ID
		err = as.Services.Audit.Update(ps.Context, nil, frm, ps.Logger)
		if err != nil {
			return "", errors.Wrapf(err, "unable to update Audit [%s]", frm.String())
		}
		msg := fmt.Sprintf("Audit [%s] updated", frm.String())
		return controller.FlashAndRedir(true, msg, frm.WebPath(), rc, ps)
	})
}

func AuditDelete(rc *fasthttp.RequestCtx) {
	controller.Act("audit.delete", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromPath(rc, as, ps)
		if err != nil {
			return "", err
		}
		err = as.Services.Audit.Delete(ps.Context, nil, ret.ID, ps.Logger)
		if err != nil {
			return "", errors.Wrapf(err, "unable to delete audit [%s]", ret.String())
		}
		msg := fmt.Sprintf("Audit [%s] deleted", ret.String())
		return controller.FlashAndRedir(true, msg, "/admin/audit", rc, ps)
	})
}

func RecordDetail(rc *fasthttp.RequestCtx) {
	controller.Act("audit.record.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := recordFromPath(rc, as, ps)
		if err != nil {
			return "", err
		}
		aud, err := as.Services.Audit.Get(ps.Context, nil, ret.AuditID, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(ret.String(), ret)
		return controller.Render(rc, as, &vaudit.RecordDetail{Model: ret, Audit: aud}, ps, "admin", auditBreadcrumb, ret.String())
	})
}

func auditFromPath(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (*audit.Audit, error) {
	idArgStr, err := cutil.RCRequiredString(rc, "id", false)
	if err != nil {
		return nil, errors.Wrap(err, "must provide [id] as an argument")
	}
	idArgP := util.UUIDFromString(idArgStr)
	if idArgP == nil {
		return nil, errors.Errorf("argument [id] (%s) is not a valid UUID", idArgStr)
	}
	idArg := *idArgP
	return as.Services.Audit.Get(ps.Context, nil, idArg, ps.Logger)
}

func auditFromForm(rc *fasthttp.RequestCtx, setPK bool) (*audit.Audit, error) {
	frm, err := cutil.ParseForm(rc)
	if err != nil {
		return nil, err
	}
	return audit.FromMap(frm, setPK)
}

func recordFromPath(rc *fasthttp.RequestCtx, as *app.State, ps *cutil.PageState) (*audit.Record, error) {
	idArgStr, err := cutil.RCRequiredString(rc, "id", false)
	if err != nil {
		return nil, errors.Wrap(err, "must provide [id] as an argument")
	}
	idArgP := util.UUIDFromString(idArgStr)
	if idArgP == nil {
		return nil, errors.Errorf("argument [id] (%s) is not a valid UUID", idArgStr)
	}
	idArg := *idArgP
	return as.Services.Audit.GetRecord(ps.Context, nil, idArg, ps.Logger)
}
