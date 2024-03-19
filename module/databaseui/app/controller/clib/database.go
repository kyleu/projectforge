package clib

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/lib/database"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vdatabase"
)

const KeyAnalyze = "analyze"

func DatabaseList(w http.ResponseWriter, r *http.Request) {
	controller.Act("database.list", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
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
		return controller.Render(w, r, as, &vdatabase.List{Keys: keys, Services: svcs}, ps, "admin", "Database")
	})
}

func DatabaseDetail(w http.ResponseWriter, r *http.Request) {
	controller.Act("database.detail", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		svc, err := getDatabaseService(r)
		if err != nil {
			return "", err
		}
		return controller.Render(w, r, as, &vdatabase.Detail{Mode: "", Svc: svc}, ps, "admin", "Database||/admin/database", svc.Key)
	})
}

func DatabaseAction(w http.ResponseWriter, r *http.Request) {
	controller.Act("database.action", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		svc, err := getDatabaseService(r)
		if err != nil {
			return "", err
		}
		act, err := cutil.RCRequiredString(r, "act", true)
		if err != nil {
			return "", err
		}
		bc := []string{"admin", "Database||/admin/database", fmt.Sprintf("%s||/admin/database/%s", svc.Key, svc.Key), act}
		switch act {
		case "enable":
			_ = svc.EnableTracing(r.URL.Query().Get("tracing"), ps.Logger)
			return "/admin/database/" + svc.Key + "/recent", nil
		case "recent":
			if idxStr := r.URL.Query().Get("idx"); idxStr != "" {
				idx, _ := strconv.ParseInt(idxStr, 10, 32)
				st := database.GetDebugStatement(svc.Key, int(idx))
				if st != nil {
					return controller.Render(w, r, as, &vdatabase.Statement{Statement: st}, ps, bc...)
				}
			}
			recent := database.GetDebugStatements(svc.Key)
			return controller.Render(w, r, as, &vdatabase.Detail{Mode: "recent", Svc: svc, Recent: recent}, ps, bc...)
		case "tables":
			sizes, dberr := svc.Sizes(ps.Context, ps.Logger)
			if dberr != nil {
				return "", errors.Wrapf(dberr, "unable to calculate sizes for database [%s]", svc.Key)
			}
			return controller.Render(w, r, as, &vdatabase.Detail{Mode: "tables", Svc: svc, Sizes: sizes}, ps, bc...)
		case KeyAnalyze:
			t := util.TimerStart()
			var tmp []any
			err = svc.Select(ps.Context, &tmp, KeyAnalyze, nil, ps.Logger)
			if err != nil {
				return "", err
			}
			msg := fmt.Sprintf("Analyzed database in [%s]", util.MicrosToMillis(t.End()))
			return controller.FlashAndRedir(true, msg, "/admin/database/"+svc.Key+"/tables", w, ps){{{ if .DatabaseUISQLEditor }}}
		case "sql":
			return controller.Render(w, r, as, &vdatabase.Detail{Mode: "sql", Svc: svc, SQL: "select 1;"}, ps, bc...){{{ end }}}
		default:
			return "", errors.Errorf("invalid database action [%s]", act)
		}
	})
}

func DatabaseTableView(w http.ResponseWriter, r *http.Request) {
	controller.Act("database.sql.run", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		prms := ps.Params.Get("table", []string{"*"}, ps.Logger).Sanitize("table")
		svc, err := getDatabaseService(r)
		if err != nil {
			return "", err
		}
		schema, _ := cutil.RCRequiredString(r, "schema", true)
		table, _ := cutil.RCRequiredString(r, "table", true)

		tbl := fmt.Sprintf("%q", table)
		if schema != "default" {
			tbl = fmt.Sprintf("%q.%q", schema, table)
		}

		q := database.SQLSelect("*", tbl, "", prms.OrderByString(), prms.Limit, prms.Offset, svc.Type)
		res, err := svc.QueryRows(ps.Context, q, nil, ps.Logger)
		ps.Data = res
		bc := []string{"admin", "Database||/admin/database", fmt.Sprintf("%s||/admin/database/%s", svc.Key, svc.Key), "Tables"}
		return controller.Render(w, r, as, &vdatabase.Results{Svc: svc, Schema: schema, Table: table, Results: res, Params: prms, Error: err}, ps, bc...)
	})
}{{{ if .DatabaseUISQLEditor }}}

func DatabaseSQLRun(w http.ResponseWriter, r *http.Request) {
	controller.Act("database.sql.run", w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		svc, err := getDatabaseService(r)
		if err != nil {
			return "", err
		}
		_ = r.ParseForm()
		f := r.PostForm
		sql := f.Get("sql")
		c := f.Get("commit")
		commit := c == util.BoolTrue
		action := f.Get("action")
		if action == KeyAnalyze {
			sql = "explain analyze " + sql
		}

		tx, err := svc.StartTransaction(ps.Logger)
		if err != nil {
			return "", errors.Wrap(err, "unable to start transaction")
		}
		defer func() { _ = tx.Rollback() }()

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

		ps.SetTitleAndData("SQL Results", results)
		page := &vdatabase.Detail{Mode: "sql", Svc: svc, SQL: sql, Columns: columns, Results: results, Timing: elapsed, Commit: commit}
		return controller.Render(w, r, as, page, ps, "admin", "Database||/admin/database", svc.Key+"||/admin/database/"+svc.Key, "Results")
	})
}{{{ end }}}

func getDatabaseService(r *http.Request) (*database.Service, error) {
	key, err := cutil.RCRequiredString(r, "key", true)
	if err != nil {
		return nil, err
	}
	svc, err := database.RegistryGet(key)
	if err != nil {
		return nil, errors.Wrapf(err, "no database found with key [%s]", key)
	}
	return svc, nil
}
