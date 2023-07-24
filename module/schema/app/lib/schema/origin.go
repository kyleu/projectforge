package schema

import (
	"github.com/samber/lo"

	"{{{ .Package }}}/app/util"
)

type Origin struct {
	Key         string `json:"key"`
	Title       string `json:"title,omitempty"`
	Icon        string `json:"icon,omitempty"`
	Description string `json:"description,omitempty"`
}

var (
	OriginMySQL      = Origin{Key: "mysql", Title: "MySQL", Icon: "database", Description: "MySQL database schema"}
	OriginPostgres   = Origin{Key: "postgres", Title: "PostgreSQL", Icon: "database", Description: "PostgreSQL database schema"}
	OriginSQLite     = Origin{Key: "sqlite", Title: "SQLite", Icon: "database", Description: "SQLite database schema"}
	OriginSQLServer  = Origin{Key: "sqlserver", Title: "SQL Server", Icon: "database", Description: "SQL Server database schema"}
	OriginGraphQL    = Origin{Key: "graphql", Title: "GraphQL", Icon: "star", Description: "GraphQL schema and queries"}
	OriginProtobuf   = Origin{Key: "protobuf", Title: "Protobuf", Icon: "star", Description: "File describing proto3 definitions"}
	OriginJSONSchema = Origin{Key: "jsonschema", Title: "JSON Schema", Icon: "star", Description: "JSON Schema definition files"}
	OriginMock       = Origin{Key: "mock", Title: "Mock", Icon: "star", Description: "Simple type that returns mock data"}
	OriginUnknown    = Origin{Key: "unknown", Title: "Unknown", Icon: "star", Description: "Not quite sure what this is"}
)

func OriginKeys() []string {
	return lo.Map(AllOrigins, func(x Origin, _ int) string {
		return x.Key
	})
}

func OriginTitles() []string {
	return lo.Map(AllOrigins, func(x Origin, _ int) string {
		return x.Title
	})
}

var AllOrigins = []Origin{OriginMySQL, OriginPostgres, OriginSQLite, OriginSQLServer, OriginGraphQL, OriginProtobuf, OriginJSONSchema, OriginMock}

func OriginFromString(key string) Origin {
	return lo.FindOrElse(AllOrigins, OriginUnknown, func(t Origin) bool {
		return t.Key == key
	})
}

func (t *Origin) String() string {
	return t.Key
}

func (t *Origin) MarshalText() ([]byte, error) {
	return []byte(t.Key), nil
}

func (t *Origin) UnmarshalText(data []byte) error {
	*t = OriginFromString(string(data))
	return nil
}
