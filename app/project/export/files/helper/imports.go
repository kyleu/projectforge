package helper

import (
	"fmt"


	"projectforge.dev/projectforge/app/lib/types"
	"projectforge.dev/projectforge/app/project/export/enum"
	"projectforge.dev/projectforge/app/project/export/golang"
	"projectforge.dev/projectforge/app/project/export/model"
)

var (
	ImpAudit         = AppImport("app/lib/audit")
	ImpApp           = AppImport("app")
	ImpAppController = AppImport("app/controller")
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
	ImpAppMenu       = AppImport("app/lib/menu")
	ImpRouter        = golang.NewImport(golang.ImportTypeExternal, "github.com/fasthttp/router")
	ImpSlices        = golang.NewImport(golang.ImportTypeExternal, "golang.org/x/exp/slices")
	ImpSQLx          = golang.NewImport(golang.ImportTypeExternal, "github.com/jmoiron/sqlx")
	ImpStrconv       = golang.NewImport(golang.ImportTypeInternal, "strconv")
	ImpStrings       = golang.NewImport(golang.ImportTypeInternal, "strings")
	ImpTime          = golang.NewImport(golang.ImportTypeInternal, "time")
	ImpUUID          = golang.NewImport(golang.ImportTypeExternal, "github.com/google/uuid")
)

func AppImport(path string) *golang.Import {
	return &golang.Import{Type: golang.ImportTypeApp, Value: "{{{ .Package }}}/" + path}
}

func ImportsForTypes(ctx string, enums enum.Enums, ts ...types.Type) golang.Imports {
	var ret golang.Imports
	for _, t := range ts {
		ret = ret.Add(importsForType(ctx, enums, t)...)
	}
	return ret
}

func importsForType(ctx string, enums enum.Enums, t types.Type) golang.Imports {
	switch ctx {
	case "go":
		return importsForTypeCtxGo(t, enums)
	case "dto":
		return importsForTypeCtxDTO(t, enums)
	case "string":
		return importsForTypeCtxString(t, enums)
	case "parse":
		return importsForTypeCtxParse(t, enums)
	case "webedit":
		return importsForTypeCtxWebEdit(t, enums)
	default:
		return golang.Imports{{Type: golang.ImportTypeInternal, Value: fmt.Sprintf("ERROR:invalid import context [%s]", ctx)}}
	}
}

func importsForTypeCtxGo(t types.Type, enums enum.Enums) golang.Imports {
	switch t.Key() {
	case types.KeyEnum:
		return importsForEnum(t, enums)
	case types.KeyMap, types.KeyValueMap:
		return golang.Imports{ImpAppUtil}
	case types.KeyTimestamp:
		return golang.Imports{ImpTime}
	case types.KeyUUID:
		return golang.Imports{ImpUUID}
	default:
		return nil
	}
}

func importsForTypeCtxDTO(t types.Type, enums enum.Enums) golang.Imports {
	switch t.Key() {
	case types.KeyAny:
		return golang.Imports{ImpJSON}
	case types.KeyEnum:
		return importsForEnum(t, enums)
	case types.KeyList, types.KeyMap, types.KeyValueMap, types.KeyReference:
		return golang.Imports{ImpJSON, ImpAppUtil}
	case types.KeyTimestamp:
		return golang.Imports{ImpTime}
	case types.KeyUUID:
		return golang.Imports{ImpUUID}
	default:
		return nil
	}
}

func importsForEnum(t types.Type, enums enum.Enums) golang.Imports {
	e, err := model.AsEnumInstance(t, enums)
	if err != nil {
		return golang.Imports{{Type: "error", Value: err.Error()}}
	}
	return golang.Imports{AppImport("app/" + e.PackageWithGroup(""))}
}

func importsForTypeCtxString(t types.Type, enums enum.Enums) golang.Imports {
	switch t.Key() {
	case types.KeyInt, types.KeyFloat:
		return golang.Imports{ImpFmt}
	case types.KeyMap, types.KeyValueMap:
		return golang.Imports{ImpAppUtil}
	default:
		return nil
	}
}

func importsForTypeCtxParse(t types.Type, enums enum.Enums) golang.Imports {
	switch t.Key() {
	case types.KeyInt, types.KeyFloat:
		return golang.Imports{ImpStrconv}
	case types.KeyUUID:
		return golang.Imports{ImpAppUtil}
	default:
		return nil
	}
}

func importsForTypeCtxWebEdit(t types.Type, enums enum.Enums) golang.Imports {
	switch t.Key() {
	case types.KeyAny:
		return golang.Imports{ImpAppUtil, ImpFmt}
	case types.KeyList, types.KeyMap, types.KeyValueMap, types.KeyReference:
		return golang.Imports{ImpAppUtil}
	default:
		return nil
	}
}
