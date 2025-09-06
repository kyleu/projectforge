package database

import (
	"strings"
	"unicode/utf8"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

const (
	localhost = "localhost"

	{{{ if not .SQLiteOnly }}}cfgHost     = "db_host"
	cfgPort     = "db_port"
	{{{ end }}}cfgUser     = "db_user"
	cfgPassword = "db_password"{{{ if not .SQLiteOnly }}}
	cfgDatabase = "db_database"{{{ end }}}
	cfgSchema   = "db_schema"{{{ if .HasModule "sqlite" }}}
	cfgFile     = "db_file"{{{ end }}}{{{ if not .SQLiteOnly }}}
	cfgMaxConns = "db_max_connections"{{{ end }}}
	cfgDebug    = "db_debug"
)

func ArrayToString(a []string) string {
	return "{" + util.StringJoin(a, ",") + "}"
}

func StringToArray(s string) []string {
	split := strings.Split(strings.TrimPrefix(strings.TrimSuffix(s, "}"), "{"), ",")
	ret := make([]string, 0)
	lo.ForEach(split, func(x string, _ int) {
		y := strings.TrimSpace(x)
		if y != "" {
			ret = append(ret, y)
		}
	})
	return ret
}

func MapScan(row *sqlx.Rows) (util.ValueMap, error) {
	x := util.ValueMap{}
	err := row.MapScan(x)
	if err != nil {
		return nil, err
	}
	for k, v := range x {
		switch t := v.(type) {
		case []byte:
			if utf8.Valid(t) {
				x[k] = string(t)
			}
		}
	}
	return x, nil
}
