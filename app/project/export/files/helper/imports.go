package helper

import (
	"fmt"

	"github.com/samber/lo"
	"golang.org/x/mod/semver"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/util"
)

var (
	ImpAudit         = AppImport("lib/audit")
	ImpApp           = &golang.Import{Type: golang.ImportTypeApp, Value: "{{{ .Package }}}/app"}
	ImpAppController = AppImport("controller")
	ImpAppDatabase   = AppImport("lib/database")
	ImpAppMenu       = AppImport("lib/menu")
	ImpAppUtil       = AppImport("util")
	ImpContext       = golang.NewImport(golang.ImportTypeInternal, "context")
	ImpComponents    = ViewImport("components")
	ImpCutil         = AppImport("controller/cutil")
	ImpDBDriver      = golang.NewImport(golang.ImportTypeInternal, "database/sql/driver")
	ImpErrors        = golang.NewImport(golang.ImportTypeExternal, "github.com/pkg/errors")
	ImpFastHTTP      = golang.NewImport(golang.ImportTypeExternal, "github.com/valyala/fasthttp")
	ImpFilter        = AppImport("lib/filter")
	ImpFmt           = golang.NewImport(golang.ImportTypeInternal, "fmt")
	ImpJSON          = golang.NewImport(golang.ImportTypeInternal, "encoding/json")
	ImpLayout        = ViewImport("layout")
	ImpLo            = golang.NewImport(golang.ImportTypeExternal, "github.com/samber/lo")
	ImpMSSQL         = golang.NewImport(golang.ImportTypeExternal, "github.com/denisenkom/go-mssqldb").WithAlias("mssql")
	ImpRouter        = golang.NewImport(golang.ImportTypeExternal, "github.com/fasthttp/router")
	ImpSearchResult  = AppImport("lib/search/result")
	ImpSlices        = golang.NewImport(golang.ImportTypeInternal, "slices")
	ImpSlices119     = golang.NewImport(golang.ImportTypeExternal, "golang.org/x/exp/slices")
	ImpSQL           = golang.NewImport(golang.ImportTypeInternal, "database/sql")
	ImpSQLx          = golang.NewImport(golang.ImportTypeExternal, "github.com/jmoiron/sqlx")
	ImpStrconv       = golang.NewImport(golang.ImportTypeInternal, "strconv")
	ImpStrings       = golang.NewImport(golang.ImportTypeInternal, "strings")
	ImpTime          = golang.NewImport(golang.ImportTypeInternal, "time")
	ImpURL           = golang.NewImport(golang.ImportTypeInternal, "net/url")
	ImpUUID          = golang.NewImport(golang.ImportTypeExternal, "github.com/google/uuid")
	ImpXML           = golang.NewImport(golang.ImportTypeInternal, "encoding/xml")
)

func AppImport(path string) *golang.Import {
	return &golang.Import{Type: golang.ImportTypeApp, Value: "{{{ .Package }}}/app/" + path}
}

func ViewImport(path string) *golang.Import {
	return &golang.Import{Type: golang.ImportTypeApp, Value: "{{{ .Package }}}/views/" + path}
}

func ImpSlicesForGo(v string) *golang.Import {
	c := semver.Compare("v"+v, "v"+project.DefaultGoVersion)
	if c < 0 {
		return ImpSlices119
	}
	return ImpSlices
}

func ImportsForTypes(ctx string, database string, ts ...types.Type) golang.Imports {
	return lo.FlatMap(ts, func(t types.Type, _ int) []*golang.Import {
		return importsForType(ctx, t, database)
	})
}

func importsForType(ctx string, t types.Type, database string) golang.Imports {
	switch ctx {
	case "go":
		return importsForTypeCtxGo(t)
	case "row":
		return importsForTypeCtxRow(t, database)
	case types.KeyString:
		return importsForTypeCtxString(t)
	case "parse":
		return importsForTypeCtxParse(t)
	case "webedit":
		return importsForTypeCtxWebEdit(t)
	default:
		return golang.Imports{{Type: golang.ImportTypeInternal, Value: fmt.Sprintf("ERROR:invalid import context [%s]", ctx)}}
	}
}

func importsForTypeCtxGo(t types.Type) golang.Imports {
	switch t.Key() {
	case types.KeyMap, types.KeyValueMap, types.KeyList:
		return golang.Imports{ImpAppUtil}
	case types.KeyDate, types.KeyTimestamp:
		return golang.Imports{ImpTime}
	case types.KeyUUID:
		return golang.Imports{ImpUUID}
	default:
		return nil
	}
}

func importsForTypeCtxRow(t types.Type, database string) golang.Imports {
	switch t.Key() {
	case types.KeyAny:
		if database == util.DatabaseSQLite {
			return nil
		}
		return golang.Imports{ImpJSON}
	case types.KeyList, types.KeyMap, types.KeyValueMap, types.KeyReference:
		if database == util.DatabaseSQLite {
			return golang.Imports{ImpAppUtil}
		}
		return golang.Imports{ImpJSON, ImpAppUtil}
	case types.KeyDate, types.KeyTimestamp:
		return golang.Imports{ImpTime}
	case types.KeyUUID:
		if database == util.DatabaseSQLServer {
			return nil
		}
		return golang.Imports{ImpUUID}
	default:
		return nil
	}
}

func importsForTypeCtxString(t types.Type) golang.Imports {
	switch t.Key() {
	case types.KeyInt, types.KeyFloat:
		return golang.Imports{ImpFmt}
	case types.KeyMap, types.KeyValueMap:
		return golang.Imports{ImpAppUtil}
	default:
		return nil
	}
}

func importsForTypeCtxParse(t types.Type) golang.Imports {
	switch t.Key() {
	case types.KeyInt, types.KeyFloat:
		return golang.Imports{ImpStrconv}
	case types.KeyUUID:
		return golang.Imports{ImpAppUtil}
	default:
		return nil
	}
}

func importsForTypeCtxWebEdit(t types.Type) golang.Imports {
	switch t.Key() {
	case types.KeyAny:
		return golang.Imports{ImpAppUtil, ImpFmt}
	case types.KeyList:
		lt := types.TypeAs[*types.List](t)
		if x := types.TypeAs[*types.Enum](lt.V); x != nil {
			return golang.Imports{}
		}
		return golang.Imports{ImpAppUtil}
	case types.KeyMap, types.KeyValueMap, types.KeyReference:
		return golang.Imports{ImpAppUtil}
	default:
		return nil
	}
}
