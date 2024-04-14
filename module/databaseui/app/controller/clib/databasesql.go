package clib

import (
	"net/http"

	"github.com/pkg/errors"

	"{{{ .Package }}}/app"
	"{{{ .Package }}}/app/controller"
	"{{{ .Package }}}/app/controller/cutil"
	"{{{ .Package }}}/app/util"
	"{{{ .Package }}}/views/vdatabase"
)

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
		}
		if commit {
			err = tx.Commit()
			if err != nil {
				return "", errors.Wrap(err, "unable to commit transaction")
			}
		} else {
			_ = tx.Rollback()
		}

		ps.SetTitleAndData("SQL Results", results)
		page := &vdatabase.Detail{Mode: "sql", Svc: svc, SQL: sql, Columns: columns, Results: results, Timing: elapsed, Commit: commit}
		return controller.Render(w, r, as, page, ps, databaseBC(svc.Key, "Results")...)
	})
}
