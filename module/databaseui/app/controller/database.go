package controller

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/valyala/fasthttp"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/database"{{{ if .DatabaseUISQLEditor }}}
	"{{{ .Package }}}/app/util"{{{ end }}}
	"{{{ .Package }}}/views/vdatabase"
)

func DatabaseList(rc *fasthttp.RequestCtx) {
	act("database.list", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		keys := database.RegistryKeys()
		if len(keys) == 1 {
			return "/admin/database/" + keys[0], nil
		}
		svcs := make(map[string]*database.Service, len(keys))
		for _, key := range keys {
			svc, err := database.RegistryGet(key)
			if err != nil {
				return "", errors.Wrapf(err, "no database found with key [%s]", key)
			}
			svcs[key] = svc
		}
		return render(rc, as, &vdatabase.List{Keys: keys, Services: svcs}, ps, "admin", "Database")
	})
}

func DatabaseDetail(rc *fasthttp.RequestCtx) {
	act("database.detail", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		svc, err := getDatabaseService(rc)
		if err != nil {
			return "", err
		}
		return render(rc, as, &vdatabase.Detail{Mode: "", Svc: svc}, ps, "admin", "Database||/admin/database", svc.Key)
	})
}

func DatabaseAction(rc *fasthttp.RequestCtx) {
	act("database.action", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		svc, err := getDatabaseService(rc)
		if err != nil {
			return "", err
		}
		act, err := RCRequiredString(rc, "act", true)
		if err != nil {
			return "", err
		}
		switch act {
		case "enable":
			_ = svc.EnableTracing(true, true, ps.Logger)
			return "/admin/database/" + svc.Key + "/recent", nil
		case "recent":
			if idxStr := string(rc.URI().QueryArgs().Peek("idx")); idxStr != "" {
				idx, _ := strconv.Atoi(idxStr)
				st := database.GetDebugStatement(svc.Key, idx)
				return render(rc, as, &vdatabase.Statement{Statement: st}, ps, "admin", "Database||/admin/database", svc.Key)
			}
			recent := database.GetDebugStatements(svc.Key)
			return render(rc, as, &vdatabase.Detail{Mode: "recent", Svc: svc, Recent: recent}, ps, "admin", "Database||/admin/database", svc.Key)
		case "tables":
			sizes, err := svc.Sizes(ps.Context, ps.Logger)
			if err != nil {
				return "", errors.Wrapf(err, "unable to calculate sizes for database [%s]", svc.Key)
			}
			return render(rc, as, &vdatabase.Detail{Mode: "tables", Svc: svc, Sizes: sizes}, ps, "admin", "Database||/admin/database", svc.Key)
		case "analyze":
			t := util.TimerStart()
			var tmp []interface{}
			err = svc.Select(ps.Context, &tmp, "analyze", nil, ps.Logger)
			if err != nil {
				return "", err
			}
			msg := fmt.Sprintf("Analyzed database in [%s]", util.MicrosToMillis(t.End()))
			return flashAndRedir(true, msg, "/admin/database/"+svc.Key+"/tables", rc, ps){{{ if .DatabaseUISQLEditor }}}
		case "sql":
			return render(rc, as, &vdatabase.Detail{Mode: "sql", Svc: svc, SQL: "select 1;"}, ps, "admin", "Database||/admin/database", svc.Key){{{ end }}}
		default:
			return "", errors.Errorf("invalid database action [%s]", act)
		}
	})
}{{{ if .DatabaseUISQLEditor }}}

func DatabaseSQLRun(rc *fasthttp.RequestCtx) {
	act("database.sql.run", rc, func(as *app.State, ps *cutil.PageState) (string, error) {
		svc, err := getDatabaseService(rc)
		if err != nil {
			return "", err
		}
		f := rc.PostArgs()
		sql := string(f.Peek("sql"))
		c := string(f.Peek("commit"))
		commit := c == util.BoolTrue
		action := string(f.Peek("action"))
		if action == "analyze" {
			sql = "explain analyze " + sql
		}

		tx, err := svc.StartTransaction(ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to start transaction")
		}

		var columns []string
		results := [][]any{}

		timer := util.TimerStart()
		result, err := svc.Query(ps.Context, sql, tx, ps.Logger)
		if err != nil {
			return "", err
		}
		defer func() { _ = result.Close() }()

		elapsed := timer.End()

		if result != nil {
			for result.Next() {
				if columns == nil {
					columns, _ = result.Columns()
				}
				row, e := result.SliceScan()
				if e != nil {
					return "", errors.Wrap(e, "unable to read row")
				}
				results = append(results, row)
			}
		}{{{ if .DatabaseUIReadOnly }}}
		_ = tx.Rollback(){{{ else }}}
		if commit {
			err = tx.Commit()
			if err != nil {
				return "", errors.Wrap(err, "unable to commit transaction")
			}
		} else {
			_ = tx.Rollback()
		}{{{ end }}}

		ps.Title = "SQL Results"
		ps.Data = results
		page := &vdatabase.Detail{Mode: "sql", Svc: svc, SQL: sql, Columns: columns, Results: results, Timing: elapsed, Commit: commit}
		return render(rc, as, page, ps, "admin", "Database||/admin/database", svc.Key+"||/admin/database/"+svc.Key, "Results")
	})
}{{{ end }}}

func getDatabaseService(rc *fasthttp.RequestCtx) (*database.Service, error) {
	key, err := RCRequiredString(rc, "key", true)
	if err != nil {
		return nil, err
	}
	svc, err := database.RegistryGet(key)
	if err != nil {
		return nil, errors.Wrapf(err, "no database found with key [%s]", key)
	}
	return svc, nil
}
