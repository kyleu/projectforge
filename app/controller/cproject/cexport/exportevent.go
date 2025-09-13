package cexport

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cproject"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/project/export/files"
	"projectforge.dev/projectforge/app/util"
	"projectforge.dev/projectforge/views/vexport"
)

func ProjectExportEventDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.event.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, evt, err := exportLoadEvent(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		fls, err := files.EventAll(prj, prj.ExportArgs.Events, util.StringDefaultLinebreak)
		if err != nil {
			ps.Logger.Warnf("unable to generate files for event [%s]", evt.Name)
		}
		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), evt.Title()}
		ps.SetTitleAndData(fmt.Sprintf("[%s] %s", prj.Key, evt.Name), evt)
		return controller.Render(r, as, &vexport.EventDetail{BaseURL: prj.WebPathEvents(), Event: evt, Files: fls}, ps, bc...)
	})
}

func ProjectExportEventNew(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.event.new", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := cproject.GetProjectWithArgs(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		mdl := &model.Event{}
		bc := []string{"projects", prj.Key, fmt.Sprintf("Export||/p/%s/export", prj.Key), "New"}
		ps.SetTitleAndData(fmt.Sprintf("[%s] New Event", prj.Key), mdl)
		return controller.Render(r, as, &vexport.EventForm{BaseURL: prj.WebPathEvents(), Event: mdl, Examples: model.ExampleEvent}, ps, bc...)
	})
}

func ProjectExportEventCreate(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.event.create", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, err := cproject.GetProjectWithArgs(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}
		mdl := &model.Event{}
		err = exportEventFromForm(frm, mdl)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse event from form")
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		err = as.Services.Projects.SaveExportEvent(pfs, mdl)
		if err != nil {
			return "", err
		}
		msg := "event created successfully"
		u := fmt.Sprintf("/p/%s/export/events/%s", prj.Key, mdl.Name)
		return controller.FlashAndRedir(true, msg, u, ps)
	})
}

func ProjectExportEventForm(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.event.form", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, err := exportLoadEvent(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		bc := []string{
			"projects",
			prj.Key,
			fmt.Sprintf("Export||/p/%s/export", prj.Key),
			mdl.Title() + "||" + prj.WebPath() + "/export/events/" + mdl.Name,
			"Edit",
		}
		ps.SetTitleAndData(fmt.Sprintf("[%s] %s", prj.Key, mdl.Name), mdl)
		return controller.Render(r, as, &vexport.EventForm{BaseURL: prj.WebPathEvents(), Event: mdl, Examples: model.ExampleEvent}, ps, bc...)
	})
}

func ProjectExportEventSave(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.event.save", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, err := exportLoadEvent(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		frm, err := cutil.ParseForm(r, ps.RequestBody)
		if err != nil {
			return "", err
		}
		err = exportEventFromForm(frm, mdl)
		if err != nil {
			return "", errors.Wrap(err, "unable to parse event from form")
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		err = as.Services.Projects.SaveExportEvent(pfs, mdl)
		if err != nil {
			return "", err
		}
		msg := "event saved successfully"
		u := fmt.Sprintf("/p/%s/export/events/%s", prj.Key, mdl.Name)
		return controller.FlashAndRedir(true, msg, u, ps)
	})
}

func ProjectExportEventDelete(w http.ResponseWriter, r *http.Request) {
	controller.Act("project.export.event.delete", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prj, mdl, err := exportLoadEvent(r, as, ps.Logger)
		if err != nil {
			return "", err
		}
		pfs, err := as.Services.Projects.GetFilesystem(prj)
		if err != nil {
			return "", err
		}
		err = as.Services.Projects.DeleteExportEvent(pfs, mdl.Name, ps.Logger)
		if err != nil {
			return "", err
		}
		msg := "event deleted successfully"
		return controller.FlashAndRedir(true, msg, fmt.Sprintf("/p/%s/export", prj.Key), ps)
	})
}
