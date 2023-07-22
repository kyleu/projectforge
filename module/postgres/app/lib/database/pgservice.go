package database

import (
	"os"
	"strconv"

	"github.com/go-ini/ini"
	"github.com/pkg/errors"
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

// Valid parameters supported in reading `PGSERVICEFILE`.
var validParams = map[string]string{
	"host":        "PGHOST",
	"port":        "PGPORT",
	"dbname":      "PGDATABASE",
	"user":        "PGUSER",
	"password":    "PGPASSWORD",
	"sslmode":     "PGSSLMODE",
	"sslcert":     "PGSSLCERT",
	"sslkey":      "PGSSLKEY",
	"sslrootcert": "PGSSLROOTCERT",
}

type PostgresServiceParams struct {
	Host        string `json:"host"`
	Port        int    `json:"port,omitempty"`
	Schema      string `json:"schema,omitempty"`
	Username    string `json:"username"`
	Password    string `json:"password,omitempty"`
	Database    string `json:"database,omitempty"`
	SSLMode     string `json:"sslmode,omitempty"`
	SSLKey      string `json:"sslkey,omitempty"`
	SSLCert     string `json:"sslcert,omitempty"`
	SSLRootCert string `json:"sslrootcert,omitempty"`
}

// PostgresParamsFromService parses connection service config from DB_SERVICE env var from file located at DB_SERVICEFILE.
// If DB_SERVICEFILE is not provided, it searches for the pg service conf file in home directory.
func PostgresParamsFromService() (*PostgresServiceParams, error) {
	pgservice := util.GetEnv("db_service")
	if pgservice == "" {
		return nil, errors.New("missing 'DB_SERVICE' env var")
	}

	pgservicefile := util.GetEnv("db_servicefile")
	if pgservicefile == "" {
		pgservicefile = os.ExpandEnv("${HOME}/.pg_service.conf")
	}

	paramMap, err := parseConfigSection(pgservice, pgservicefile)
	if err != nil {
		return nil, err
	}

	port, err := strconv.ParseInt(paramMap["port"], 10, 32)
	if err != nil {
		port = 5432
	}
	params := PostgresServiceParams{
		Host:        paramMap["host"],
		Port:        int(port),
		Username:    paramMap["user"],
		Password:    paramMap["password"],
		Database:    paramMap["dbname"],
		SSLMode:     paramMap["sslmode"],
		SSLKey:      paramMap["sslkey"],
		SSLCert:     paramMap["sslcert"],
		SSLRootCert: paramMap["sslrootcert"],
	}

	return &params, nil
}

// parseConfigSection parses options specified in a config section of a pg service file and returns them as a map.
func parseConfigSection(service, file string) (map[string]string, error) {
	result := make(map[string]string)

	cfg, err := ini.Load(file)
	if err != nil {
		return result, errors.Errorf("error loading pg service file at '%s'", file)
	}

	cfg.BlockMode = false

	section, err := cfg.GetSection(service)
	if err != nil {
		return result, err
	}

	lo.ForEach(lo.Keys(validParams), func(key string, _ int) {
		if value, err := section.GetKey(key); err == nil {
			result[key] = value.String()
		}
	})

	return result, nil
}
