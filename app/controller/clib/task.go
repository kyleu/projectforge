package clib

import (
	"fmt"
	"net/http"

	"github.com/pkg/errors"

	"projectforge.dev/projectforge/app"
	"projectforge.dev/projectforge/app/controller"
	"projectforge.dev/projectforge/app/controller/cutil"
	"projectforge.dev/projectforge/app/lib/task"
	"projectforge.dev/projectforge/app/lib/websocket"
	"projectforge.dev/projectforge/app/util"
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
		args, cat := argsAndCategory(r)
		page := &vtask.Detail{Task: t, Args: args}
		if args.GetBoolOpt("async") {
			delete(args, "async")
			page.SocketURL = t.WebPath() + "/start" + "?" + args.ToQueryString()
			ps.SetTitleAndData(t.TitleSafe(), t)
		} else {
			ret := t.Run(ps.Context, cat, args, ps.Logger)
			page.Result = ret
			ps.SetTitleAndData(t.TitleSafe(), ret)
		}
		return controller.Render(r, as, page, ps, taskBC(t, "Run**play")...)
	})
}

func TaskStart(w http.ResponseWriter, r *http.Request) {
	controller.Act("task.start", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		key, err := cutil.PathString(r, "key", false)
		if err != nil {
			return "", err
		}
		t := as.Services.Task.RegisteredTasks.Get(key)
		if t == nil {
			return "", errors.Errorf("unable to find task [%s]", key)
		}
		args, cat := argsAndCategory(r)
		ch := fmt.Sprintf("%s-%d", t.Key, util.RandomInt(1000))
		id, err := as.Services.Socket.Upgrade(ps.Context, w, r, ch, ps.Profile, websocket.EchoHandler, ps.Logger)
		if err != nil {
			ps.Logger.Warnf("unable to upgrade connection to WebSocket: %s", err.Error())
			return "", err
		}

		fn := as.Services.Socket.Terminal(ch, ps.Logger)
		go func() {
			res := t.Run(ps.Context, cat, args, ps.Logger, fn)
			html := vtask.ResultSummary(as, res, ps)
			x := util.ValueMap{"result": res, "html": html}
			_ = as.Services.Socket.WriteMessage(id, websocket.NewMessage(nil, ch, "complete", x), ps.Logger)
			_ = as.Services.Socket.WriteCloseRequest(id, ps.Logger)
		}()

		return "", as.Services.Socket.ReadLoop(ps.Context, id, ps.Logger)
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

func argsAndCategory(r *http.Request) (util.ValueMap, string) {
	args := cutil.QueryArgsMap(r.URL)
	category := args.GetStringOpt("category")
	if category == "" {
		category = "ad-hoc"
	}
	return args.WithoutKeys("category"), category
}

func taskBC(t *task.Task, extra ...string) []string {
	ret := []string{keyAdmin, "Tasks||/admin/task**task"}
	if t != nil {
		ret = append(ret, t.Breadcrumb())
	}
	return append(ret, extra...)
}
