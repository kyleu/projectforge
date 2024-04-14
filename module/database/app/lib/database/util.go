package database

import (
	"strings"

	"github.com/samber/lo"
)

const (
	localhost = "localhost"

	cfgHost     = "db_host"
	cfgPort     = "db_port"
	cfgUser     = "db_user"
	cfgPassword = "db_password"
	cfgDatabase = "db_database"
	cfgSchema   = "db_schema"
	cfgFile     = "db_file"
	cfgMaxConns = "db_max_connections"
	cfgDebug    = "db_debug"
)

func ArrayToString(a []string) string {
	return "{" + strings.Join(a, ",") + "}"
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
