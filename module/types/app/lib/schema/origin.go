package schema

import (
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
	OriginGraphQL    = Origin{Key: "graphql", Title: "GraphQL", Icon: "star", Description: "GraphQL schema and queries"}
	OriginProtobuf   = Origin{Key: "protobuf", Title: "Protobuf", Icon: "star", Description: "File describing proto3 definitions"}
	OriginJSONSchema = Origin{Key: "jsonschema", Title: "JSON Schema", Icon: "star", Description: "JSON Schema definition files"}
	OriginMock       = Origin{Key: "mock", Title: "Mock", Icon: "star", Description: "Simple type that returns mock data"}
	OriginUnknown    = Origin{Key: "unknown", Title: "Unknown", Icon: "star", Description: "Not quite sure what this is"}
)

func OriginKeys() []string {
	ret := make([]string, 0, len(AllOrigins))
	for _, x := range AllOrigins {
		ret = append(ret, x.Key)
	}
	return ret
}

func OriginTitles() []string {
	ret := make([]string, 0, len(AllOrigins))
	for _, x := range AllOrigins {
		ret = append(ret, x.Title)
	}
	return ret
}

var AllOrigins = []Origin{OriginMySQL, OriginPostgres, OriginSQLite, OriginGraphQL, OriginProtobuf, OriginJSONSchema, OriginMock}

func OriginFromString(s string) Origin {
	for _, t := range AllOrigins {
		if t.Key == s {
			return t
		}
	}
	return OriginUnknown
}

func (t *Origin) String() string {
	return t.Key
}

func (t *Origin) MarshalJSON() ([]byte, error) {
	return util.ToJSONBytes(t.Key, false), nil
}

func (t *Origin) UnmarshalJSON(data []byte) error {
	var s string
	if err := util.FromJSON(data, &s); err != nil {
		return err
	}
	*t = OriginFromString(s)
	return nil
}
