package files

import (
	"fmt"

	"github.com/kyleu/projectforge/app/export/golang"
	"github.com/kyleu/projectforge/app/export/model"
)

var (
	impAppUtil = golang.NewImport(golang.ImportTypeApp, "{{{ .Package }}}/app/util")
	impFmt     = golang.NewImport(golang.ImportTypeInternal, "fmt")
	impJSON    = golang.NewImport(golang.ImportTypeInternal, "encoding/json")
	impStrconv = golang.NewImport(golang.ImportTypeInternal, "strconv")
	impTime    = golang.NewImport(golang.ImportTypeInternal, "time")
	impUUID    = golang.NewImport(golang.ImportTypeExternal, "github.com/google/uuid")
)

func importsForTypes(ctx string, types ...*model.Type) golang.Imports {
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
		return golang.Imports{impAppUtil}
	case model.TypeTimestamp.Key:
		return golang.Imports{impTime}
	case model.TypeUUID.Key:
		return golang.Imports{impUUID}
	default:
		return nil
	}
}

func importsForTypeCtxDTO(t *model.Type) golang.Imports {
	switch t.Key {
	case model.TypeMap.Key:
		return golang.Imports{impJSON, impAppUtil}
	case model.TypeTimestamp.Key:
		return golang.Imports{impTime}
	case model.TypeUUID.Key:
		return golang.Imports{impUUID}
	default:
		return nil
	}
}

func importsForTypeCtxString(t *model.Type) golang.Imports {
	switch t.Key {
	case model.TypeInt.Key:
		return golang.Imports{impFmt}
	case model.TypeMap.Key:
		return golang.Imports{impAppUtil}
	default:
		return nil
	}
}

func importsForTypeCtxParse(t *model.Type) golang.Imports {
	switch t.Key {
	case model.TypeInt.Key:
		return golang.Imports{impStrconv}
	case model.TypeUUID.Key:
		return golang.Imports{impAppUtil}
	default:
		return nil
	}
}

func importsForTypeCtxWebEdit(t *model.Type) golang.Imports {
	switch t.Key {
	case model.TypeMap.Key:
		return golang.Imports{impAppUtil}
	default:
		return nil
	}
}
