package helper

import (
	"fmt"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
)

var (
	ImpApp        = AppImport("app")
	ImpAppUtil    = AppImport("app/util")
	ImpContext    = golang.NewImport(golang.ImportTypeInternal, "context")
	ImpComponents = AppImport("views/components")
	ImpCutil      = AppImport("app/controller/cutil")
	ImpDatabase   = AppImport("app/lib/database")
	ImpErrors     = golang.NewImport(golang.ImportTypeExternal, "github.com/pkg/errors")
	ImpFastHTTP   = golang.NewImport(golang.ImportTypeExternal, "github.com/valyala/fasthttp")
	ImpFilter     = AppImport("app/lib/filter")
	ImpFmt        = golang.NewImport(golang.ImportTypeInternal, "fmt")
	ImpJSON       = golang.NewImport(golang.ImportTypeInternal, "encoding/json")
	ImpLayout     = AppImport("views/layout")
	ImpLogging    = golang.NewImport(golang.ImportTypeExternal, "go.uber.org/zap")
	ImpSQLx       = golang.NewImport(golang.ImportTypeExternal, "github.com/jmoiron/sqlx")
	ImpStrconv    = golang.NewImport(golang.ImportTypeInternal, "strconv")
	ImpStrings    = golang.NewImport(golang.ImportTypeInternal, "strings")
	ImpTime       = golang.NewImport(golang.ImportTypeInternal, "time")
	ImpUUID       = golang.NewImport(golang.ImportTypeExternal, "github.com/google/uuid")
)

func AppImport(path string) *golang.Import {
	return &golang.Import{Type: golang.ImportTypeApp, Value: "{{{ .Package }}}/" + path}
}

func ImportsForTypes(ctx string, types ...*model.Type) golang.Imports {
	var ret golang.Imports
	for _, t := range types {
		ret = ret.Add(importsForType(ctx, t)...)
	}
	return ret
}

func importsForType(ctx string, t *model.Type) golang.Imports {
	switch ctx {
	case "go":
		return importsForTypeCtxGo(t)
	case "dto":
		return importsForTypeCtxDTO(t)
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

func importsForTypeCtxGo(t *model.Type) golang.Imports {
	switch t.Key {
	case model.TypeMap.Key:
		return golang.Imports{ImpAppUtil}
	case model.TypeTimestamp.Key:
		return golang.Imports{ImpTime}
	case model.TypeUUID.Key:
		return golang.Imports{ImpUUID}
	default:
		return nil
	}
}

func importsForTypeCtxDTO(t *model.Type) golang.Imports {
	switch t.Key {
	case model.TypeInterface.Key:
		return golang.Imports{ImpJSON}
	case model.TypeMap.Key:
		return golang.Imports{ImpJSON, ImpAppUtil}
	case model.TypeTimestamp.Key:
		return golang.Imports{ImpTime}
	case model.TypeUUID.Key:
		return golang.Imports{ImpUUID}
	default:
		return nil
	}
}

func importsForTypeCtxString(t *model.Type) golang.Imports {
	switch t.Key {
	case model.TypeInt.Key:
		return golang.Imports{ImpFmt}
	case model.TypeMap.Key:
		return golang.Imports{ImpAppUtil}
	default:
		return nil
	}
}

func importsForTypeCtxParse(t *model.Type) golang.Imports {
	switch t.Key {
	case model.TypeInt.Key:
		return golang.Imports{ImpStrconv}
	case model.TypeUUID.Key:
		return golang.Imports{ImpAppUtil}
	default:
		return nil
	}
}

func importsForTypeCtxWebEdit(t *model.Type) golang.Imports {
	switch t.Key {
	case model.TypeInterface.Key:
		return golang.Imports{ImpAppUtil, ImpFmt}
	case model.TypeMap.Key:
		return golang.Imports{ImpAppUtil}
	default:
		return nil
	}
}
