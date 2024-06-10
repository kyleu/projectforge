package helper

import (
	"fmt"

	"github.com/samber/lo"
	"golang.org/x/mod/semver"

	"projectforge.dev/projectforge/app/lib/metamodel/model"
	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project"
	"projectforge.dev/projectforge/app/util"
)

var (
	ImpAudit          = AppImport("lib/audit")
	ImpApp            = &model.Import{Type: model.ImportTypeApp, Value: "{{{ .Package }}}/app"}
	ImpAppController  = AppImport("controller")
	ImpAppDatabase    = AppImport("lib/database")
	ImpAppMenu        = AppImport("lib/menu")
	ImpAppSvc         = AppImport("lib/svc")
	ImpAppUtil        = AppImport("util")
	ImpContext        = model.NewImport(model.ImportTypeInternal, "context")
	ImpComponents     = ViewImport("components")
	ImpComponentsView = ViewImport("components/view")
	ImpComponentsEdit = ViewImport("components/edit")
	ImpCutil          = AppImport("controller/cutil")
	ImpDBDriver       = model.NewImport(model.ImportTypeInternal, "database/sql/driver")
	ImpErrors         = model.NewImport(model.ImportTypeExternal, "github.com/pkg/errors")
	ImpHTTP           = model.NewImport(model.ImportTypeInternal, "net/http")
	ImpFilter         = AppImport("lib/filter")
	ImpFmt            = model.NewImport(model.ImportTypeInternal, "fmt")
	ImpJSON           = model.NewImport(model.ImportTypeInternal, "encoding/json")
	ImpLayout         = ViewImport("layout")
	ImpLo             = model.NewImport(model.ImportTypeExternal, "github.com/samber/lo")
	ImpMSSQL          = model.NewImport(model.ImportTypeExternal, "github.com/denisenkom/go-mssqldb").WithAlias("mssql")
	ImpRouter         = model.NewImport(model.ImportTypeExternal, "github.com/gorilla/mux")
	ImpSearchResult   = AppImport("lib/search/result")
	ImpSlices         = model.NewImport(model.ImportTypeInternal, "slices")
	ImpSlices119      = model.NewImport(model.ImportTypeExternal, "model.org/x/exp/slices")
	ImpSQL            = model.NewImport(model.ImportTypeInternal, "database/sql")
	ImpSQLx           = model.NewImport(model.ImportTypeExternal, "github.com/jmoiron/sqlx")
	ImpStrconv        = model.NewImport(model.ImportTypeInternal, "strconv")
	ImpStrings        = model.NewImport(model.ImportTypeInternal, "strings")
	ImpTime           = model.NewImport(model.ImportTypeInternal, "time")
	ImpURL            = model.NewImport(model.ImportTypeInternal, "net/url")
	ImpUUID           = model.NewImport(model.ImportTypeExternal, "github.com/google/uuid")
	ImpXML            = model.NewImport(model.ImportTypeInternal, "encoding/xml")
)

func AppImport(path string) *model.Import {
	return &model.Import{Type: model.ImportTypeApp, Value: "{{{ .Package }}}/app/" + path}
}

func ViewImport(path string) *model.Import {
	return &model.Import{Type: model.ImportTypeApp, Value: "{{{ .Package }}}/views/" + path}
}

func ImpSlicesForGo(v string) *model.Import {
	c := semver.Compare("v"+v, "v"+project.DefaultGoVersion)
	if c < 0 {
		return ImpSlices119
	}
	return ImpSlices
}

func ImportsForTypes(ctx string, database string, ts ...types.Type) model.Imports {
	return lo.FlatMap(ts, func(t types.Type, _ int) []*model.Import {
		return importsForType(ctx, t, database)
	})
}

func importsForType(ctx string, t types.Type, database string) model.Imports {
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
		return model.Imports{{Type: model.ImportTypeInternal, Value: fmt.Sprintf("ERROR:invalid import context [%s]", ctx)}}
	}
}

func importsForTypeCtxGo(t types.Type) model.Imports {
	switch t.Key() {
	case types.KeyMap, types.KeyValueMap, types.KeyList:
		return model.Imports{ImpAppUtil}
	case types.KeyDate, types.KeyTimestamp:
		return model.Imports{ImpTime}
	case types.KeyUUID:
		return model.Imports{ImpUUID}
	default:
		return nil
	}
}

func importsForTypeCtxRow(t types.Type, database string) model.Imports {
	switch t.Key() {
	case types.KeyAny:
		if SimpleJSON(database) {
			return nil
		}
		return model.Imports{ImpJSON}
	case types.KeyList, types.KeyMap, types.KeyValueMap, types.KeyReference:
		if SimpleJSON(database) {
			return model.Imports{ImpAppUtil}
		}
		return model.Imports{ImpJSON, ImpAppUtil}
	case types.KeyDate, types.KeyTimestamp:
		return model.Imports{ImpTime}
	case types.KeyUUID:
		if database == util.DatabaseSQLServer {
			return nil
		}
		return model.Imports{ImpUUID}
	default:
		return nil
	}
}

func importsForTypeCtxString(t types.Type) model.Imports {
	switch t.Key() {
	case types.KeyInt, types.KeyFloat:
		return model.Imports{ImpFmt}
	case types.KeyMap, types.KeyValueMap:
		return model.Imports{ImpAppUtil}
	default:
		return nil
	}
}

func importsForTypeCtxParse(t types.Type) model.Imports {
	switch t.Key() {
	case types.KeyInt, types.KeyFloat:
		return model.Imports{ImpStrconv}
	case types.KeyUUID:
		return model.Imports{ImpAppUtil}
	default:
		return nil
	}
}

func importsForTypeCtxWebEdit(t types.Type) model.Imports {
	switch t.Key() {
	case types.KeyAny:
		return model.Imports{ImpAppUtil, ImpFmt}
	case types.KeyList:
		lt := types.TypeAs[*types.List](t)
		if x := types.TypeAs[*types.Enum](lt.V); x != nil {
			return model.Imports{}
		}
		return model.Imports{ImpAppUtil}
	case types.KeyMap, types.KeyValueMap, types.KeyReference:
		return model.Imports{ImpAppUtil}
	default:
		return nil
	}
}
