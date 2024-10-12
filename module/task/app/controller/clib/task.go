package clib

import (
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/task"
	"projectforge.dev/projectforge/views/vtask"
)

const taskIcon = "gift"

func TaskList(w http.ResponseWriter, r *http.Request) {
	controller.Act("task.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		tasks := as.Services.Task.RegisteredTasks
		ps.SetTitleAndData("Tasks", tasks)
		ps.DefaultNavIcon = taskIcon
		return controller.Render(r, as, &vtask.List{Tasks: tasks}, ps, taskBC(nil)...)
	})
}

func TaskDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("task.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", false)
		if err != nil {
			return "", err
		}
		t := as.Services.Task.RegisteredTasks.Get(key)
		if t == nil {
			return "", errors.Errorf("no task found with key [%s]", key)
		}
		ps.SetTitleAndData(t.TitleSafe(), t)
		ps.DefaultNavIcon = taskIcon
		return controller.Render(r, as, &vtask.Detail{Task: t}, ps, taskBC(t)...)
	})
}

func TaskRun(w http.ResponseWriter, r *http.Request) {
	controller.Act("task.run", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", false)
		if err != nil {
			return "", err
		}
		t := as.Services.Task.RegisteredTasks.Get(key)
		if t == nil {
			return "", errors.Errorf("unable to find task [%s]", key)
		}
		args := cutil.QueryArgsMap(r.URL)
		ret := t.Run(ps.Context, args, ps.Logger)
		page := &vtask.Detail{Task: t, Result: ret}
		ps.SetTitleAndData(t.TitleSafe(), ret)
		return controller.Render(r, as, page, ps, taskBC(t, "Run**play")...)
	})
}

func TaskRemove(w http.ResponseWriter, r *http.Request) {
	controller.Act("task.remove", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", false)
		if err != nil {
			return "", err
		}
		as.Services.Task.RemoveTask(key)
		if err != nil {
			return "", errors.Wrap(err, "unable to remove task")
		}
		return controller.ReturnToReferrer("removed task ["+key+"]", "/admin/task", ps)
	})
}

func taskBC(t *task.Task, extra ...string) []string {
	ret := []string{"admin", "Tasks||/admin/task**task"}
	if t != nil {
		ret = append(ret, t.TitleSafe()+"||"+t.WebPath()+"**"+t.IconSafe())
	}
	return append(ret, extra...)
}
