package helper

import (
	"fmt"

	"github.com/samber/lo"

	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

var (
	ImpAudit         = AppImport("app/lib/audit")
	ImpApp           = AppImport("app")
	ImpAppController = AppImport("app/controller")
	ImpAppMenu       = AppImport("app/lib/menu")
	ImpAppUtil       = AppImport("app/util")
	ImpContext       = golang.NewImport(golang.ImportTypeInternal, "context")
	ImpComponents    = AppImport("views/components")
	ImpCutil         = AppImport("app/controller/cutil")
	ImpDatabase      = AppImport("app/lib/database")
	ImpErrors        = golang.NewImport(golang.ImportTypeExternal, "github.com/pkg/errors")
	ImpFastHTTP      = golang.NewImport(golang.ImportTypeExternal, "github.com/valyala/fasthttp")
	ImpFilter        = AppImport("app/lib/filter")
	ImpFmt           = golang.NewImport(golang.ImportTypeInternal, "fmt")
	ImpJSON          = golang.NewImport(golang.ImportTypeInternal, "encoding/json")
	ImpLayout        = AppImport("views/layout")
	ImpMSSQL         = golang.NewImport(golang.ImportTypeExternal, "github.com/denisenkom/go-mssqldb").WithAlias("mssql")
	ImpRouter        = golang.NewImport(golang.ImportTypeExternal, "github.com/fasthttp/router")
	ImpSlices        = golang.NewImport(golang.ImportTypeExternal, "golang.org/x/exp/slices")
	ImpSQL           = golang.NewImport(golang.ImportTypeInternal, "database/sql")
	ImpSQLx          = golang.NewImport(golang.ImportTypeExternal, "github.com/jmoiron/sqlx")
	ImpStrconv       = golang.NewImport(golang.ImportTypeInternal, "strconv")
	ImpStrings       = golang.NewImport(golang.ImportTypeInternal, "strings")
	ImpTime          = golang.NewImport(golang.ImportTypeInternal, "time")
	ImpUUID          = golang.NewImport(golang.ImportTypeExternal, "github.com/google/uuid")
)

func AppImport(path string) *golang.Import {
	return &golang.Import{Type: golang.ImportTypeApp, Value: "{{{ .Package }}}/" + path}
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
	case "string":
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
		return golang.Imports{ImpJSON}
	case types.KeyList, types.KeyMap, types.KeyValueMap, types.KeyReference:
		return golang.Imports{ImpJSON, ImpAppUtil}
	case types.KeyDate, types.KeyTimestamp:
		return golang.Imports{ImpTime}
	case types.KeyUUID:
		if database == model.SQLServer {
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
	case types.KeyList, types.KeyMap, types.KeyValueMap, types.KeyReference:
		return golang.Imports{ImpAppUtil}
	default:
		return nil
	}
}
