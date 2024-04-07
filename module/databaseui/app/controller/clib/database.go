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
	act, err := cutil.PathString(r, "act", true)
	controller.Act("database.action."+act, w, r, func(as *app.State, ps *cutil.PageState) (string, error) {
		if err != nil {
			return "", err
		}
		svc, err := getDatabaseService(r)
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
			return controller.FlashAndRedir(true, msg, "/admin/database/"+svc.Key+"/tables", w, ps)
		case "sql":
			return controller.Render(w, r, as, &vdatabase.Detail{Mode: "sql", Svc: svc, SQL: "select 1;"}, ps, bc...)
		default:
			return "", errors.Errorf("invalid database action [%s]", act)
		}
	})
}

func getDatabaseService(r *http.Request) (*database.Service, error) {
	key, err := cutil.PathString(r, "key", true)
	if err != nil {
		return nil, err
	}
	svc, err := database.RegistryGet(key)
	if err != nil {
		return nil, errors.Wrapf(err, "no database found with key [%s]", key)
	}
	return svc, nil
}
