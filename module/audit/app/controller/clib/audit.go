package clib

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/audit"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vaudit"
)

const auditDefaultTitle = "Audits"

var auditBreadcrumb = "Audit||/admin/audit"

func AuditList(w http.ResponseWriter, r *http.Request) {
	controller.Act("audit.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := as.Services.Audit.List(ps.Context, nil, ps.Params.Get("audit", nil, ps.Logger), ps.Logger)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(auditDefaultTitle, ret)
		return controller.Render(w, r, as, &vaudit.List{Models: ret, Params: ps.Params}, ps, "admin", "Audit")
	})
}

func AuditDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("audit.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromPath(r, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(ret.String(), ret)
		records, err := as.Services.Audit.RecordsForAudits(ps.Context, nil, ps.Params.Get("auditRecord", nil, ps.Logger), ps.Logger, ret.ID)
		if err != nil {
			return "", errors.Wrap(err, "unable to retrieve child auditrecords")
		}
		return controller.Render(w, r, as, &vaudit.Detail{Model: ret, Params: ps.Params, Records: records}, ps, "admin", auditBreadcrumb, ret.String())
	})
}

func AuditCreateForm(w http.ResponseWriter, r *http.Request) {
	controller.Act("audit.create.form", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := &audit.Audit{}
		ps.SetTitleAndData("Create [Audit]", ret)
		return controller.Render(w, r, as, &vaudit.Edit{Model: ret, IsNew: true}, ps, "admin", auditBreadcrumb, "Create")
	})
}

func AuditCreateFormRandom(w http.ResponseWriter, r *http.Request) {
	controller.Act("audit.create.form.random", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret := audit.Random()
		ps.SetTitleAndData("Create Random [Audit]", ret)
		return controller.Render(w, r, as, &vaudit.Edit{Model: ret, IsNew: true}, ps, "admin", auditBreadcrumb, "Create")
	})
}

func AuditCreate(w http.ResponseWriter, r *http.Request) {
	controller.Act("audit.create", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromForm(r, ps.RequestBody, true)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse Audit from form")
		}
		err = as.Services.Audit.Create(ps.Context, nil, ps.Logger, ret)
		if err != nil {
			return "", errors.Wrap(err, "unable to save newly-created Audit")
		}
		msg := fmt.Sprintf("Audit [%s] created", ret.String())
		return controller.FlashAndRedir(true, msg, ret.WebPath(), w, ps)
	})
}

func AuditEditForm(w http.ResponseWriter, r *http.Request) {
	controller.Act("audit.edit.form", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromPath(r, as, ps)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData("Edit ["+ret.String()+"]", ret)
		return controller.Render(w, r, as, &vaudit.Edit{Model: ret}, ps, "admin", auditBreadcrumb, ret.String())
	})
}

func AuditEdit(w http.ResponseWriter, r *http.Request) {
	controller.Act("audit.edit", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromPath(r, as, ps)
		if err != nil {
			return "", err
		}
		frm, err := auditFromForm(r, ps.RequestBody, false)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse Audit from form")
		}
		frm.ID = ret.ID
		err = as.Services.Audit.Update(ps.Context, nil, frm, ps.Logger)
		if err != nil {
			return "", errors.Wrapf(err, "unable to update Audit [%s]", frm.String())
		}
		msg := fmt.Sprintf("Audit [%s] updated", frm.String())
		return controller.FlashAndRedir(true, msg, frm.WebPath(), w, ps)
	})
}

func AuditDelete(w http.ResponseWriter, r *http.Request) {
	controller.Act("audit.delete", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := auditFromPath(r, as, ps)
		if err != nil {
			return "", err
		}
		err = as.Services.Audit.Delete(ps.Context, nil, ret.ID, ps.Logger)
		if err != nil {
			return "", errors.Wrapf(err, "unable to delete audit [%s]", ret.String())
		}
		msg := fmt.Sprintf("Audit [%s] deleted", ret.String())
		return controller.FlashAndRedir(true, msg, "/admin/audit", w, ps)
	})
}

func RecordDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("audit.record.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		ret, err := recordFromPath(r, as, ps)
		if err != nil {
			return "", err
		}
		aud, err := as.Services.Audit.Get(ps.Context, nil, ret.AuditID, ps.Logger)
		if err != nil {
			return "", err
		}
		ps.SetTitleAndData(ret.String(), ret)
		return controller.Render(w, r, as, &vaudit.RecordDetail{Model: ret, Audit: aud}, ps, "admin", auditBreadcrumb, ret.String())
	})
}

func auditFromPath(r *http.Request, as *app.State, ps *cutil.PageState) (*audit.Audit, error) {
	idArgStr, err := cutil.PathString(r, "id", false)
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

func auditFromForm(r *http.Request, requestBody []byte, setPK bool) (*audit.Audit, error) {
	frm, err := cutil.ParseForm(r, requestBody)
	if err != nil {
		return nil, err
	}
	return audit.FromMap(frm, setPK)
}

func recordFromPath(r *http.Request, as *app.State, ps *cutil.PageState) (*audit.Record, error) {
	idArgStr, err := cutil.PathString(r, "id", false)
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
